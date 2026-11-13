package topology

import (
	"context"
	"errors"
	"testing"

	"breaklab/internal/catalog"
	"breaklab/internal/runner"
	"breaklab/internal/web/webcore"
)

func fixtureCatalogForTopology() *catalog.Catalog {
	return &catalog.Catalog{Stages: []catalog.StageGroup{
		{
			Stage: catalog.Stage{ID: "linux", Title: "Linux fundamentals"},
			Labs: []catalog.Lab{
				{
					ID:    "web-down",
					Title: "The website returns nothing",
					Environment: catalog.Environment{
						Ports: []string{"8080:80"},
					},
				},
			},
		},
	}}
}

func fakeStatus(state, image string, err error) webcore.StatusFunc {
	return func(ctx context.Context, name string) (string, string, error) {
		return state, image, err
	}
}

func TestBuildTopologyNoActiveLab(t *testing.T) {
	view := buildTopology(context.Background(), fixtureCatalogForTopology(), nil, fakeStatus("running", "img", nil), nil, nil)
	if view.Active {
		t.Fatalf("Active = true with no session, want false")
	}
	if len(view.Containers) != 0 {
		t.Fatalf("Containers = %+v, want none", view.Containers)
	}
}

func TestBuildTopologyUnknownLabInSession(t *testing.T) {
	session := &runner.Session{Lab: "no-such-lab", Container: "breaklab-no-such-lab", Variant: "x"}
	view := buildTopology(context.Background(), fixtureCatalogForTopology(), session, fakeStatus("running", "img", nil), nil, nil)
	if view.Active {
		t.Fatalf("Active = true for a lab not in the catalog, want false")
	}
}

func TestBuildTopologyActiveLabRunning(t *testing.T) {
	session := &runner.Session{Lab: "web-down", Container: "breaklab-web-down", Variant: "config-typo"}
	view := buildTopology(context.Background(), fixtureCatalogForTopology(), session, fakeStatus("running", "breaklab-web-down", nil), nil, nil)

	if !view.Active {
		t.Fatalf("Active = false, want true")
	}
	if view.LabID != "web-down" {
		t.Errorf("LabID = %q, want the catalog id (not the display title)", view.LabID)
	}
	if len(view.Containers) != 1 {
		t.Fatalf("Containers = %+v, want exactly one", view.Containers)
	}
	c := view.Containers[0]
	if c.Name != "breaklab-web-down" || c.Status != "running" || c.Image != "breaklab-web-down" {
		t.Errorf("container = %+v", c)
	}
	if c.ExecCmd != "docker exec -it breaklab-web-down bash" {
		t.Errorf("ExecCmd = %q", c.ExecCmd)
	}
	if len(c.Ports) != 1 || c.Ports[0].HostPort != "8080" || c.Ports[0].ContainerPort != "80" {
		t.Fatalf("Ports = %+v", c.Ports)
	}
	if c.Ports[0].OpenHref != "http://localhost:8080" {
		t.Errorf("OpenHref = %q", c.Ports[0].OpenHref)
	}
}

func TestBuildTopologyLogsAndStats(t *testing.T) {
	session := &runner.Session{Lab: "web-down", Container: "breaklab-web-down", Variant: "config-typo"}
	logs := func(ctx context.Context, name string, tailLines int) ([]byte, error) {
		if name != "breaklab-web-down" || tailLines != 20 {
			t.Errorf("logs called with (%q, %d)", name, tailLines)
		}
		return []byte("  line one\nline two  "), nil
	}
	stats := func(ctx context.Context, name string) (webcore.Stats, error) {
		return webcore.Stats{CPUPercent: 12.5, MemUsedBytes: 1024, MemLimitBytes: 4096}, nil
	}

	view := buildTopology(context.Background(), fixtureCatalogForTopology(), session, fakeStatus("running", "img", nil), logs, stats)

	c := view.Containers[0]
	if c.LogTail != "line one\nline two" {
		t.Errorf("LogTail = %q, want trimmed log output", c.LogTail)
	}
	if c.Stats != (webcore.Stats{CPUPercent: 12.5, MemUsedBytes: 1024, MemLimitBytes: 4096}) {
		t.Errorf("Stats = %+v", c.Stats)
	}
}

func TestBuildTopologyNilLogsAndStatsDegradeGracefully(t *testing.T) {
	session := &runner.Session{Lab: "web-down", Container: "breaklab-web-down", Variant: "config-typo"}
	view := buildTopology(context.Background(), fixtureCatalogForTopology(), session, fakeStatus("running", "img", nil), nil, nil)

	c := view.Containers[0]
	if c.LogTail != "" {
		t.Errorf("LogTail = %q, want empty with a nil logs func", c.LogTail)
	}
	if c.Stats != (webcore.Stats{}) {
		t.Errorf("Stats = %+v, want zero value with a nil stats func", c.Stats)
	}
}

// web-ui spec, "Docker unreachable" scenario: a lookup error degrades
// status to "unknown" instead of failing the page.
func TestBuildTopologyDockerErrorDegradesToUnknown(t *testing.T) {
	session := &runner.Session{Lab: "web-down", Container: "breaklab-web-down", Variant: "config-typo"}
	view := buildTopology(context.Background(), fixtureCatalogForTopology(), session, fakeStatus("", "", errors.New("daemon unreachable")), nil, nil)

	if !view.Active {
		t.Fatalf("Active = false, want true (session data should still render)")
	}
	if view.Containers[0].Status != "unknown" {
		t.Errorf("Status = %q, want unknown", view.Containers[0].Status)
	}
}

func TestBuildTopologyNilStatusFuncDegradesToUnknown(t *testing.T) {
	session := &runner.Session{Lab: "web-down", Container: "breaklab-web-down", Variant: "config-typo"}
	view := buildTopology(context.Background(), fixtureCatalogForTopology(), session, nil, nil, nil)

	if view.Containers[0].Status != "unknown" {
		t.Errorf("Status = %q, want unknown", view.Containers[0].Status)
	}
}
