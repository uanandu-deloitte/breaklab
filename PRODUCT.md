# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

Individual engineers building hands-on DevOps/SRE troubleshooting skill —
starting with the repo owner, who uses it as their own personal training
ground, but built so any other engineer can clone the repo and run their
own private instance the same way. Each instance is single-user; there is
no multi-tenant or shared-account use case.

## Product Purpose

BreakLab builds a fresh, throwaway Docker environment, secretly applies
one random breakage, and hands the user a symptom-only incident ticket
(never the cause). The user debugs inside the container with no
instructions, `breaklab check` verifies the symptom is gone, and
`breaklab postmortem` requires a real write-up (Summary, Detection, Root
cause, Fix, Prevention; 20+ words each) before the lab counts as
complete. Success means durable incident-response skill, built the way
real incidents are worked, with a written record the user keeps.

## Positioning

Modeled on hands-on break/fix platforms like devopslabs.dev (symptom-only
incident briefs, a roadmap from first shell to Kubernetes, no
instructions) but reimplemented standalone for one person, running fully
locally: Docker is the only dependency, there is no cloud account or
service to subscribe to, progress is plain files (`progress/progress.json`,
`postmortems/*.md`) committed to the user's own git repo rather than rows
in someone else's database, and the whole tool ships as one
statically-linked binary fast enough to run well on a Raspberry Pi 5.

## Operating Context

- CLI workflow: `start <lab>` (build + break) → debug inside the
  container → `check` → `hint` (tiered: general hints first, then the
  active variant's, one at a time, counted) → `postmortem` → repeat;
  `reset`, `stop`, `status`, `list` support the loop.
- `breaklab serve` runs a local web UI (dashboard, roadmap, challenge
  browser filterable by stage/difficulty/status/topic tag, lab
  detail/postmortem viewer) reading live progress files, no publish
  step. Binds to `127.0.0.1` by default; exposing it on a LAN is opt-in
  (`--host`) and prints a one-time warning, since the UI has no
  authentication by design (single-user tool).
- `breaklab publish` rebuilds a public static site (`site/`) from
  progress, the lab catalog, and published postmortems, deployed to
  GitHub Pages by the user's own CI on every push to `main`. Nothing
  with `publish: false` in a postmortem's front matter, and no
  variant/solution script, ever reaches `site/`.
- Labs are data, not code: each `labs/<id>/` folder holds `lab.yaml`,
  `TICKET.md` (symptom only), a `Dockerfile`, `check.sh`, `variants/`
  (each applies one breakage), and `solutions/` (reference fixes used
  only by integration tests, never run for the user). Adding a lab never
  touches the CLI's Go code.
- `roadmap.yaml` is the single source of stage order (Linux fundamentals
  → Networking & Web → Git → Docker → CI/CD → Infrastructure as code →
  Cloud fundamentals (Azure) → Kubernetes fundamentals/developer/admin/
  security → Monitoring & logging → Capstone) and drives the CLI
  listing, the web UI, and the public site identically.
- Built and tracked as a sequence of OpenSpec changes
  (`openspec/specs/` = merged requirements, `openspec/changes/` =
  in-flight proposals, `docs/breaklab-spec.md` = full build spec).

## Capabilities and Constraints

- Only required dependency beyond the binary itself is Docker (plus Git
  for owning your own progress/postmortem history). Go 1.26+ to build
  from source.
- Ships as one statically-linked binary; a `linux/arm64` cross-compile
  is a CI build-smoke check for the Raspberry Pi 5 target.
- Outcome-based checks: a check verifies the symptom is gone (e.g. the
  site serves the right page), not a specific fix, so any valid repair
  passes.
- Random, non-repeating variant selection: each lab has two or more ways
  to break, weighted selection, never the same variant twice in a row.
- Security posture is a product constraint, not an afterthought:
  `os/exec` with argument slices only (no shell interpolation), strict
  id/path validation (`^[a-z0-9-]+$`), `govulncheck` as a CI dependency
  audit gate, SHA-pinned GitHub Actions.

## Brand Commitments

- Name "BreakLab" is fixed. Positioned explicitly as a standalone,
  personal-use reimplementation of the devopslabs.dev format (incident
  briefs, roadmap, certification-style structure) — that lineage is a
  stated reference, not a claim of affiliation.
- The shipped visual system is **"Monochrome Minimalist"**, a
  high-contrast editorial light theme recorded in `DESIGN.md` (with its
  `.impeccable/design.json` sidecar): a pure white canvas and cards,
  pure black ink/borders/accent, heavy 3-4px structural rules, and a
  small deliberate set of status hues (ok/warning/danger/info) kept as
  real color on purpose — protected from the monochrome collapse by a
  design-critique pass that caught an earlier draft over-collapsing
  `--warning` to black too. It replaced the earlier warm **"Bright
  Progress"** system (cream canvas, indigo accent), which itself had
  replaced the dark-zinc **"Pedagogical Clarity"** system; a light theme
  remains the product's only theme. Lives in `templates/web/style.css`
  (kept byte-identical with `templates/site/style.css` for shared
  tokens), Geist + JetBrains Mono carried over unchanged (JetBrains Mono
  now also covers every form input, not just code), with contrast
  enforced by a test (`TestDesignTokensMeetWCAGAAContrast`).
  Status/difficulty signaling never relies on color alone. The 1rem
  accessibility floor carried across every system unchanged — it's an
  accessibility commitment, not a side effect of any one theme.
- `breaklab-design/` (the four dark-theme page mockups plus their
  reference DESIGN.md) is now historical reference material for the
  superseded "Pedagogical Clarity" direction, not the current system.
- No confirmed logo or other brand assets exist yet. The reference
  screenshot that pinned this direction is another company's branded
  product (Codecademy); BreakLab's system borrows its palette/layout
  family, not its name, logo, or wordmark.

## Evidence on Hand

- Lab catalog: 150+ labs under `labs/<id>/`, spanning every roadmap
  stage from Linux fundamentals through Cloud fundamentals, the four
  Kubernetes tracks, Observability, and a Capstone.
- One real postmortem on file:
  `postmortems/2026-10-18-web-down-bad-perms.md`.
- `openspec/specs/` and `docs/breaklab-spec.md` hold the full,
  currently-merged requirements and build spec for every area of the
  app.
- No testimonials, pricing, benchmarks, or third-party case studies
  exist or should be invented; this is a personal/shareable tool, not a
  commercial product.

## Product Principles

1. Local-first and owned: no server or account beyond Docker; progress
   and postmortems are plain files in the user's own git repo, not rows
   in someone else's database.
2. Symptom-only, outcome-based learning: never reveal the cause; verify
   the fix by effect, so any valid repair counts, not just the one
   solution script got right.
3. The write-up is part of the lab: a lab isn't "done" until its
   postmortem has real content — reflection is treated as how the skill
   sticks, not paperwork bolted on after.
4. Minimal footprint, maximum portability: one statically-linked binary,
   fast enough for a Raspberry Pi 5, Docker the only dependency.
5. Security and reliability are load-bearing, not optional: argument-slice
   exec, strict input validation, CI vulnerability scanning, and
   SHA-pinned Actions apply even though this is a personal tool.

## Accessibility & Inclusion

Contrast is a tested, enforced requirement, not a guideline: `--fg` and
`--fg-muted` must clear 4.5:1, and the status/accent color set
(`--ok`/`--warning`/`--danger`/`--info`/`--accent`) must clear 3.0:1,
both against `--bg`, verified by
`TestDesignTokensMeetWCAGAAContrast`. Status and difficulty signaling
must never rely on color alone — it always pairs a hue with a text
label, icon, or shape/shade distinction.
