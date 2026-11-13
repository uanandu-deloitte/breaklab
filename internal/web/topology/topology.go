package topology

import (
	"context"
	"net/http"
	"strings"

	"breaklab/internal/catalog"
	"breaklab/internal/runner"
	"breaklab/internal/web/webcore"
)

// ContainerNodeView is the active lab's container. A slice of one today
// in TopologyView.Containers, so a later multi-container (compose-labs)
// change is additive rather than a restructure (design.md Decision 1).
type ContainerNodeView struct {
	Name    string
	Image   string
	Status  string // Docker's running state, or "unknown" when it couldn't be read
	ExecCmd string
	Ports   []webcore.PortNodeView

	// LogTail is the container's last few log lines, same best-effort
	// contract as Status: empty when logs is nil or the lookup fails,
	// never a page failure (design.md Decision 4).
	LogTail string

	// Stats is the container's live CPU/memory snapshot, zero-valued
	// under the same best-effort contract as LogTail.
	Stats webcore.Stats
}

// TopologyView is the Astro topology page's data (GET /api/topology).
type TopologyView struct {
	Active bool
	Lab    string

	// LabID is the catalog id (not Lab's display title) the topology
	// page needs to call the existing POST /api/labs/{id}/check route
	// for the "run check" action — runner.Check ignores the path's id
	// and always checks the active session's own lab, but the route
	// still expects a path segment, so this is that segment.
	LabID      string
	Containers []ContainerNodeView

	// Elements is Containers (plus the implicit host node) as a
	// Cytoscape.js elements array — see cytoscapeElements — ready to
	// assign straight into cytoscape({elements: ...}) with no further
	// transformation client-side.
	Elements []cytoElement
}

// cytoElement is one Cytoscape.js graph element (a node or an edge);
// Source/Target are set only on edges, Kind/Status only on nodes, which
// is why every field here is optional — a single flat shape matches
// Cytoscape's own elements-array format directly, needing no per-kind
// Go types.
type cytoElement struct {
	Data struct {
		ID     string `json:"id"`
		Label  string `json:"label"`
		Kind   string `json:"kind,omitempty"`
		Status string `json:"status,omitempty"`
		Source string `json:"source,omitempty"`
		Target string `json:"target,omitempty"`
	} `json:"data"`
}

func cytoNode(id, label, kind, status string) cytoElement {
	var e cytoElement
	e.Data.ID, e.Data.Label, e.Data.Kind, e.Data.Status = id, label, kind, status
	return e
}

func cytoEdge(id, source, target, label string) cytoElement {
	var e cytoElement
	e.Data.ID, e.Data.Source, e.Data.Target, e.Data.Label = id, source, target, label
	return e
}

// cytoscapeElements builds the host -> container -> port tree as a
// Cytoscape.js elements array, marshaled as this response's JSON
// directly (no client-side parsing of an embedded string needed).
func cytoscapeElements(containers []ContainerNodeView) []cytoElement {
	elements := []cytoElement{cytoNode("host", "you\nthis machine", "host", "")}
	for _, c := range containers {
		containerID := "container-" + c.Name
		elements = append(elements, cytoNode(containerID, c.Name+"\n"+c.Status, "container", c.Status))
		elements = append(elements, cytoEdge("edge-"+containerID, "host", containerID, "runs"))
		for _, p := range c.Ports {
			portID := "port-" + p.HostPort
			elements = append(elements, cytoNode(portID, p.HostPort+" → "+p.ContainerPort, "port", ""))
			elements = append(elements, cytoEdge("edge-"+portID, containerID, portID, "exposes"))
		}
	}
	return elements
}

// buildTopology assembles TopologyView for the active session, if any.
// status is called only with the container name already recorded in
// session — never anything from a request (web-ui spec's "Web UI attack
// surface covers Docker-backed routes" requirement) — and is expected to
// already carry whatever timeout the caller wants; a lookup error or a
// nil status func degrades the container's Status to "unknown" rather
// than failing the page (design.md Decision 4). logs and stats follow
// the same nil-degrades-gracefully contract, each leaving its field
// zero-valued rather than failing the page.
func buildTopology(ctx context.Context, cat *catalog.Catalog, session *runner.Session, status webcore.StatusFunc, logs webcore.LogsFunc, stats webcore.StatsFunc) TopologyView {
	if session == nil {
		return TopologyView{}
	}

	lab, _, found := cat.FindLab(session.Lab)
	if !found {
		return TopologyView{}
	}

	state := "unknown"
	var image string
	if status != nil {
		if s, img, err := status(ctx, session.Container); err == nil && s != "" {
			state, image = s, img
		}
	}

	var logTail string
	if logs != nil {
		if out, err := logs(ctx, session.Container, 20); err == nil {
			logTail = strings.TrimSpace(string(out))
		}
	}

	var liveStats webcore.Stats
	if stats != nil {
		if s, err := stats(ctx, session.Container); err == nil {
			liveStats = s
		}
	}

	var ports []webcore.PortNodeView
	for _, p := range lab.Environment.Ports {
		host, container, ok := strings.Cut(p, ":")
		pv := webcore.PortNodeView{HostPort: host, ContainerPort: container}
		if ok && host != "" {
			pv.OpenHref = "http://localhost:" + host
		}
		ports = append(ports, pv)
	}

	containers := []ContainerNodeView{{
		Name:    session.Container,
		Image:   image,
		Status:  state,
		ExecCmd: "docker exec -it " + session.Container + " bash",
		Ports:   ports,
		LogTail: logTail,
		Stats:   liveStats,
	}}

	return TopologyView{
		Active:     true,
		Lab:        lab.Title,
		LabID:      lab.ID,
		Containers: containers,
		Elements:   cytoscapeElements(containers),
	}
}

// topologyAPIHandler serves the Astro topology page's data: the same
// buildTopology computation the old html/template handler used, as
// JSON.
func topologyAPIHandler(cfg webcore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cat, err := webcore.LoadCatalog(cfg)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		session, err := runner.LoadSession(cfg.SessionPath)
		if err != nil {
			webcore.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), webcore.DockerStatusTimeout)
		defer cancel()

		view := buildTopology(ctx, cat, session, cfg.DockerStatus, cfg.ContainerLogs, cfg.ContainerStats)
		webcore.WriteJSON(w, view)
	}
}

// RegisterRoutes mounts the topology JSON API route on mux.
func RegisterRoutes(mux *http.ServeMux, cfg webcore.Config) {
	mux.HandleFunc("GET /api/topology", topologyAPIHandler(cfg))
}
