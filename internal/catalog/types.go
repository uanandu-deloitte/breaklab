package catalog

// Roadmap is the parsed shape of roadmap.yaml.
type Roadmap struct {
	Stages []Stage `yaml:"stages"`
}

// Stage is one entry in roadmap.yaml.
type Stage struct {
	ID    string `yaml:"id"`
	Title string `yaml:"title"`

	// Blurb is one short line on why this stage comes next; shown next
	// to its connector on the roadmap tree view. Optional — a stage
	// with no blurb just shows no annotation.
	Blurb string `yaml:"blurb"`

	// Capstone marks this stage as the gated finale, shown on its own
	// page rather than folded into the regular stage-by-stage roadmap.
	// At most one stage should set this — Catalog.CapstoneStage finds
	// the first and RegularStages excludes every stage that sets it.
	Capstone bool `yaml:"capstone"`
}

// Lab is the parsed shape of a labs/<id>/lab.yaml file.
type Lab struct {
	ID               string      `yaml:"id"`
	Title            string      `yaml:"title"`
	Stage            string      `yaml:"stage"`
	Difficulty       int         `yaml:"difficulty"`
	EstimatedMinutes int         `yaml:"estimated_minutes"`
	Environment      Environment `yaml:"environment"`
	Check            Check       `yaml:"check"`
	Hints            []string    `yaml:"hints"`
	Variants         []Variant   `yaml:"variants"`
	Tags             []string    `yaml:"tags"`

	// Summary is one line describing the lab, shown next to its title
	// in the roadmap's room-list rows. Optional — an empty summary just
	// shows no second line.
	Summary string `yaml:"summary"`

	// Loop is this lab's learn/guided/practice on-ramp, shown before
	// the break/fix ticket. Optional — a lab with a zero-value Loop
	// just skips straight to the ticket, same as before this field
	// existed.
	Loop Loop `yaml:"loop"`

	// Type classifies this lab as one standalone step of the learning
	// loop (a separate lab per step, chained by Prerequisites) rather
	// than a single lab containing every step as a section. Empty for
	// every lab written before this field existed — those keep
	// rendering the way they always have. One of "guided", "practice",
	// "break-fix", "incident", or "project" when set.
	Type string `yaml:"type"`

	// PathID and Module are informational grouping metadata (shown as a
	// breadcrumb on the lab page) from the curriculum doc's
	// stage→path→module hierarchy. They don't affect which roadmap
	// stage a lab belongs to or how it's sorted there.
	PathID string `yaml:"path"`
	Module string `yaml:"module"`

	// Objectives and Skills are informational display lists: observable
	// outcomes and named skills this lab exercises.
	Objectives []string `yaml:"objectives"`
	Skills     []string `yaml:"skills"`

	// Prerequisites lists lab ids that must be "complete" before this
	// lab is available — see progress.ApplyStatus's locking pass. Empty
	// for every lab unless it genuinely depends on another.
	Prerequisites []string `yaml:"prerequisites"`

	// Lesson is a relative path to a markdown concept primer, rendered
	// through internal/markdown. Optional — distinct from Loop.Learn
	// (plain text, no markdown): a lab sets one or the other, never
	// both.
	Lesson string `yaml:"lesson"`

	// Checks is an ordered list of named verification checks, each run
	// and reported individually. Optional and additive: a lab with a
	// non-empty Checks uses this instead of the single legacy Check
	// field; a lab with neither keeps working as it always has.
	Checks []NamedCheck `yaml:"checks"`

	// Status is computed by the catalog loader, never parsed from YAML.
	Status string `yaml:"-"`
}

// NamedCheck is one of a lab's Checks: like Check, but with an id so its
// pass/fail can be reported individually rather than collapsed into one
// overall message.
type NamedCheck struct {
	ID     string `yaml:"id"`
	Script string `yaml:"script"`
	RunIn  string `yaml:"run_in"`
}

// Loop is a lab's on-ramp content: a concept primer, a fully-explained
// guided walkthrough, and a light rehearsal — all read-through content,
// never persisted to progress.json (only the break/fix check and the
// capstone carry real tracked completion).
type Loop struct {
	// Learn is a short (2-4 sentence) primer on the concept this lab
	// tests, shown before the ticket.
	Learn string `yaml:"learn"`

	// Guided is an ordered list of fully-explained steps that walk
	// through diagnosing and fixing this lab's category of problem.
	Guided []GuidedStep `yaml:"guided"`

	// Practice is a short rehearsal task plus a few revealable
	// question/answer pairs, no sandbox environment required.
	Practice PracticeBlock `yaml:"practice"`
}

// GuidedStep is one step of a Loop's guided walkthrough: a command to
// run and why it matters.
type GuidedStep struct {
	Do  string `yaml:"do"`
	Why string `yaml:"why"`
}

// PracticeBlock is a Loop's light rehearsal exercise.
type PracticeBlock struct {
	Task      string             `yaml:"task"`
	Questions []PracticeQuestion `yaml:"questions"`
}

// PracticeQuestion is one revealable question/answer pair in a
// PracticeBlock.
type PracticeQuestion struct {
	Prompt string `yaml:"prompt"`
	Answer string `yaml:"answer"`
}

// Environment describes how a lab's container is built and exposed.
type Environment struct {
	Dockerfile string   `yaml:"dockerfile"`
	Ports      []string `yaml:"ports"`

	// Privileged opts a lab into `docker run --privileged`, which grants
	// the container meaningfully more access to the host (devices,
	// kernel modules, weaker isolation) than every other lab gets by
	// default. It exists only because some environments (e.g. a real
	// single-node Kubernetes control plane) cannot start without cgroup
	// access that an unprivileged container is denied. False for every
	// lab unless explicitly declared true.
	Privileged bool `yaml:"privileged"`

	// Volumes bind-mounts a fresh, per-attempt host directory into the
	// container at ContainerPath — the runner creates and empties the
	// host side on every Start/Reset, so a lab never inherits another
	// attempt's leftover files. Empty for every lab unless a lab
	// genuinely needs real bind-mount semantics (e.g. host-side
	// ownership shadowing a Dockerfile's own chown).
	Volumes []Volume `yaml:"volumes"`

	// PidsLimit, when > 0, adds `docker run --pids-limit`. Zero means no
	// limit — the Docker default.
	PidsLimit int `yaml:"pids_limit"`
}

// Volume is one bind mount a lab's environment declares. Name must be a
// valid id (^[a-z0-9-]+$): it becomes part of the host-side directory
// the runner creates, never raw user input.
type Volume struct {
	Name          string `yaml:"name"`
	ContainerPath string `yaml:"container_path"`

	// Named, when true, backs this volume with a daemon-managed Docker
	// volume instead of a bind-mounted host directory. Use this when a
	// lab needs real multi-UID permission enforcement on the mount
	// (some Docker backends, notably Docker Desktop's default Mac file
	// sharing, don't honor chown to a non-root uid on a host bind
	// mount). False (the default) keeps the regular host-dir behavior.
	Named bool `yaml:"named"`
}

// Check describes how a lab's fix is verified.
type Check struct {
	Script string `yaml:"script"`
	RunIn  string `yaml:"run_in"`
}

// Variant is one possible breakage a lab can apply.
type Variant struct {
	ID     string   `yaml:"id"`
	Script string   `yaml:"script"`
	Weight int      `yaml:"weight"`
	Hints  []string `yaml:"hints"`
}

// StageGroup is every lab under one roadmap stage, in roadmap order.
type StageGroup struct {
	Stage Stage
	Labs  []Lab
}

// Catalog is the result of loading and validating the roadmap and every lab.
type Catalog struct {
	Stages []StageGroup

	// Warnings holds one message per lab that failed to load or validate;
	// the rest of the catalog still loads around them.
	Warnings []string
}

// RegularStages returns every stage group except the capstone, in
// roadmap order — what the normal stage-by-stage roadmap and challenge
// browser views should iterate over.
func (c *Catalog) RegularStages() []StageGroup {
	out := make([]StageGroup, 0, len(c.Stages))
	for _, g := range c.Stages {
		if !g.Stage.Capstone {
			out = append(out, g)
		}
	}
	return out
}

// CapstoneStage returns the capstone stage group, or nil if roadmap.yaml
// declares none.
func (c *Catalog) CapstoneStage() *StageGroup {
	for i := range c.Stages {
		if c.Stages[i].Stage.Capstone {
			return &c.Stages[i]
		}
	}
	return nil
}

// FindLab returns the lab matching labID, across every stage, and its
// stage's title.
func (c *Catalog) FindLab(labID string) (Lab, string, bool) {
	for _, group := range c.Stages {
		for _, lab := range group.Labs {
			if lab.ID == labID {
				return lab, group.Stage.Title, true
			}
		}
	}
	return Lab{}, "", false
}

// FindVariant returns l's variant matching variantID, if any.
func (l Lab) FindVariant(variantID string) (Variant, bool) {
	for _, v := range l.Variants {
		if v.ID == variantID {
			return v, true
		}
	}
	return Variant{}, false
}
