# AGENTS.md

Platform: BreakLab — breaks a throwaway Docker environment in a random way, gives
the learner only a symptom, and checks their fix automatically (see
`docs/breaklab-spec.md`). Separately, there is a private, team-only deployment of
BreakLab on Proxmox (Authelia + TOTP, Caddy, a provisioning service, per-teammate
LXCs) — see the Architecture, Playground, Proxmox Runbook, and Security & Code Audit
artifacts. That deployment has no public-facing instance and is not a lab topic
itself; it's infrastructure this repo's owner runs and must be able to operate
without me.

## Environment
- Labs run as Docker containers, driven by shelling out to the `docker` CLI
  (`internal/docker`) — `os/exec` with argument slices only, never a shell string.
  Names (container, image, volume) are validated against `^[a-z0-9-]+$` before they
  touch Docker or a filesystem path.
- The web UI and hosted-account layer are Go (`internal/web/*`), served by
  `breaklab serve`. The hosted `--db` account/login system has no MFA — it must stay
  disabled (`breaklab serve` with no `--db`) on the Proxmox team deployment, where
  Authelia is the single front door. See the Security & Code Audit artifact.
- No cloud VMs or Kubernetes labs yet (Docker only; `kind`-based Kubernetes labs are
  future curriculum, not infra this repo runs itself).

## Curriculum (stages, from `roadmap.yaml`)
Linux fundamentals · Networking & Web · Git · Docker · CI/CD ·
Infrastructure as code · Cloud fundamentals (Azure) · Kubernetes fundamentals ·
Kubernetes for developers · Kubernetes administration · Kubernetes security ·
Monitoring and logging · Capstone

Each stage holds many `lab_id`s (see `labs/*/lab.yaml`, e.g. `web-down`, `docker-001`,
`lnx-014`). A new lab is added to a stage, never a new stage, unless the learner
request is explicitly "add a stage."

## Agents (authoring a new lab)

| Agent | Owns | Reads | Writes |
|---|---|---|---|
| Orchestrator | Plan, split, merge | Learner request, `roadmap.yaml`, existing `labs/` | `spec` |
| Mentor | `learn`/`guided`/`practice` narrative | `spec` | `lab.yaml`'s `loop:` block |
| Builder | `Dockerfile`, `variants/*.sh`, `check.sh` | `spec`, existing `labs/<similar-id>/` | `labs/<id>/` (minus `loop:`) |
| Verifier | Correctness | `spec`, final `lab.yaml`, final `labs/<id>/` | `report` |

- The Orchestrator writes the spec once. Mentor and Builder then work from it in
  parallel; neither waits on the other.
- Mentor refers to files by the paths in the spec and never writes `Dockerfile` or
  script contents. Builder writes no narrative prose (the `learn`/`why` text).
- Verifier starts with a fresh context: it sees only the final `lab.yaml` and
  `labs/<id>/`, not drafts or either agent's reasoning. It reports problems; it
  doesn't fix them.
- The Orchestrator routes each finding to the agent that owns that file. After 2 fix
  rounds on the same finding, it escalates to the learner (or, outside a learner
  session, to whoever asked for the lab).
- No agent reads anything outside its Reads column.

### Spec (Orchestrator → Mentor, Builder)
```yaml
lab_id: 04-timers
stage: linux              # must be one of the Curriculum stages above
difficulty: 1-5
estimated_minutes: 20
concepts: [max 3, one line each]
environment:
  dockerfile: Dockerfile
  ports: ["host:container"]
check:
  script: check.sh         # exit 0 = lab complete
variants: [id, weight, one-line hint]  # from docs/breaklab-spec.md's variant model
```

### Verifier report
```
PASS|FAIL
- <file|field>: <issue> → <owner>
```
One line per finding. Don't list what passed.

### Verifier checklist
- `docker build` succeeds from the lab's `Dockerfile` with no warnings treated as fatal.
- `check.sh` exits 0 on a lab where the fix has been applied, and nonzero on a fresh,
  unfixed container — both are tested, not assumed.
- Every `variants/*.sh` actually leaves the environment broken and exits 0 itself
  (per `docs/breaklab-spec.md`'s "Variant scripts" rule).
- Every command in the `guided` loop and every `hints`/`practice` answer matches what
  the target container actually prints — run it, don't guess.
- Lab/image/container/volume names used anywhere in the spec match
  `^[a-z0-9-]+$` (the same pattern `internal/docker` enforces at runtime).
- `concepts` has at most 3 entries; `guided` introduces one new idea per step and
  links to an earlier lab instead of re-explaining it.

## Lab narrative format (Mentor)
Tone: a mentor next to the learner — second person, short sentences, explain why
once per concept, then let the commands teach. This is the `loop:` block in
`lab.yaml` (see `labs/web-down/lab.yaml` for a worked example):
- `learn`: 2-4 sentences, the one idea this lab turns on.
- `guided`: a `do`/`why` pair per command — `why` is one line, never restates `do`.
- `practice`: a predict-then-reveal question per guided command, not a restatement
  of `guided` itself.
- Hints (top-level `hints:` and per-variant `hints:`) point at a diagnostic command,
  never at the literal fix.

## Deployment & ops docs (teaching discipline)
This applies to the Architecture, Playground, Proxmox Runbook, and Security & Code
Audit artifacts, and anything like them — documents describing how to run and
maintain BreakLab's own Proxmox/Authelia/provisioning deployment, as opposed to a
learner-facing lab.

- The reader is the repo's owner, learning to operate infrastructure they'll
  maintain themselves — not someone receiving a finished system to click through
  once. Write accordingly: no step skipped because it's "obvious," no command
  without the one-line reason it's there.
- Every non-trivial architectural concept gets a diagram (mermaid: flowchart for
  topology/context, sequence for request flows) — not prose describing a diagram
  that doesn't exist. A diagram that just restates the paragraph above it is
  redundant; cut one of the two.
- Steps are numbered, runnable as written, and state the expected result, the same
  discipline the lab format uses for learners.
- When a screenshot or photo would answer "what should I be looking at," include
  one; don't substitute a wall of text for an image that already exists or is cheap
  to produce (e.g. a Proxmox console view, a Tailscale admin screen).
- Keep the four deployment artifacts cross-linked and mutually consistent — a
  change to one (a new hardening step, a renamed service) gets reflected in the
  others in the same pass, not left to drift.

## Token discipline (all agents)
- Load only what your Reads column allows. Request files by path, never whole trees.
- Write spec, report, and handoffs in the formats above, with no surrounding prose.
- Show only the key lines of expected command output; trim the rest with `...`.
- Never restate the spec, another lab, or another agent's output — refer to it by
  `lab_id` or path.
- When revising, return only the changed step, field, or file section.

## Code discipline
- Write the minimum code that correctly solves the task. No speculative features,
  unused params, wrapper functions that add nothing, single-implementation
  interfaces, or config nobody asked for.
- No logging, comments, validation, or error handling beyond what the task needs.
  Handle real errors (returned `error` values, I/O, user input); don't guard against
  impossible states.
- Reuse existing helpers and patterns in this codebase before writing new ones.
  Never duplicate logic that already exists.
- Delete dead code you touch. No commented-out code, no TODOs.
- Prefer Go stdlib. A new module dependency needs a one-line justification in the
  PR/commit message.
- Default to efficient algorithms and data structures: no repeated work, no
  needless allocations, no N+1 calls (file reads, subprocess spawns, HTTP calls in
  a loop).

## Context discipline
- Treat the files as currently shown (via Read, not memory) as source of truth. If
  a file hasn't been read in its current form this session, read it before editing it.
- Don't reintroduce code, names, or patterns that were removed or replaced.
- Keep changes scoped to the task. Don't refactor, rename, or reformat unrelated code.
- If a request conflicts with existing code or an earlier decision, flag it in one
  line before proceeding.

## Destructive actions
- Never delete or overwrite state files (e.g. `progress.json`, lab state) without
  explicit confirmation — this repo has no safety net for that data (see memory:
  no-rm-on-real-state-files).

## Output format
- Builder: code only, a diff or the changed functions. Full files only when new or
  explicitly requested.
- Mentor: the `loop:` block only, in the format above.
- Orchestrator / Verifier: spec / report formats only.
- Deployment & ops docs: follow the teaching discipline above; otherwise, no intro,
  recap, or summary. State assumptions and non-obvious decisions in one line each.
