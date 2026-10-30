# BreakLab: Build Spec for OpenSpec + Claude

Oct 17, 2026 · @Ananda Narayanan Udayakumar

## Overview

BreakLab is a personal command-line app that breaks a throwaway Docker environment in a random way, gives you only a symptom, and checks your fix automatically. The first version ships four features: random variants, postmortems, a local web UI, and a public progress page.

BreakLab is modeled on sites like devopslabs.dev — throwaway machines, a symptom-only incident brief, no instructions, one roadmap from first shell to Kubernetes — but fixes the two things that matter most for personal use: everything runs on your own machine (no cloud-only labs you can't touch offline), and your progress and postmortems are plain files you own and can commit to git, not rows in someone else's database.

The core loop for every lab is the same:

1. `breaklab start <lab>` builds a fresh container and secretly applies one random breakage.
2. You read the ticket, which describes only the symptom, and debug inside the container.
3. `breaklab check` tells you whether the symptom is gone.
4. `breaklab postmortem` opens a pre-filled postmortem you must complete before the lab counts as done.
5. `breaklab serve` runs a local web UI — roadmap, dashboard, challenge browser — reading your live progress.
6. `breaklab publish` rebuilds your public progress site from completed labs and postmortems.

**Goals for v1**

- Each lab has at least two variants, and the same variant never repeats twice in a row.
- Checks test the outcome (the site serves the right page), not a specific fix, so any valid repair passes.
- A lab is "complete" only when its check passes and its postmortem is filled in.
- A local web UI lets you browse the roadmap, challenges, and your own dashboard without waiting on a publish step.
- The CLI and local web UI ship as a single lightweight, statically-linked binary (Go) with no language runtime or dependency install step beyond Docker, so it runs well on constrained hardware like a Raspberry Pi 5.
- The local web UI is highly interactive (live filtering, status updates) without a heavy client-side framework or JS build step — server-rendered HTML plus HTMX.
- The progress site deploys to GitHub Pages through your own CI pipeline, which is itself DevOps practice.
- A security audit covers process execution, input/path validation, secret handling, and dependencies before v1 is considered done — this app shells out to Docker, runs scripts as root inside containers, serves a web UI that can be LAN-reachable, and publishes a public site, so that surface gets checked deliberately rather than assumed safe.

**Non-goals for v1**

- No accounts, auth, or multi-user support, even when the local web UI is deployed LAN-reachable (for example on a headless Raspberry Pi) — it is still meant for the one person running it, just no longer confined to one machine.
- No cloud VMs or Kubernetes labs yet. Docker containers only; Kubernetes via kind comes later.
- No anti-cheating. Variant details sit on disk, so it runs on the honor system.

## How to use this doc with OpenSpec and Claude

Build BreakLab as nine small OpenSpec changes, one at a time, with Claude Code doing the planning and coding while you review. Each section below feeds one or more of those changes.

**One-time setup**

1. Install Go 1.22 or newer, Docker, and Git. Node.js 20.19 or newer is only needed to install the OpenSpec CLI tool itself, not required by BreakLab at runtime. On Windows, do all of this inside WSL2.
2. Install the CLI: `npm install -g @fission-ai/openspec@latest`.
3. Create the repo: `mkdir breaklab && cd breaklab && git init`.
4. Run `openspec init` and pick Claude Code when asked which tools to set up.
5. Paste the YAML from the Project context section into `openspec/config.yaml`. OpenSpec injects this context into every planning request.

**The loop for every change**

1. `/opsx:explore` (optional): talk the change through with Claude before committing to a plan.
2. `/opsx:propose <change-name>`: Claude creates `openspec/changes/<change-name>/` with `proposal.md`, `design.md`, `tasks.md`, and spec deltas.
3. Review the files yourself. Run `openspec validate <change-name>` and `openspec show <change-name>`. Edit anything that is wrong before any code exists.
4. `/opsx:apply`: Claude implements the tasks and ticks them off in `tasks.md`.
5. Run the tests and try the feature by hand. Commit.
6. `/opsx:archive`: merges the spec deltas into `openspec/specs/` and moves the change into the archive.

Start a fresh Claude session for each change. OpenSpec works best with a clean context window, and the specs on disk carry what Claude needs to know.

Sources: [OpenSpec OPSX docs](https://openspec.dev/docs/opsx), [OpenSpec migration guide](https://lzw.me/docs/openspec/en/migration-guide.html), [OpenSpec repository](https://github.com/Fission-AI/OpenSpec).

## Project context

Copy this into `openspec/config.yaml` so Claude knows the stack and rules in every planning step. Go was picked over the easier-to-read Python default because v1 targets a Raspberry Pi 5: a single static binary with no runtime or virtualenv to install beats a Python runtime and dependency tree on constrained hardware, and HTMX gives the local web UI real interactivity without a JS build pipeline.

```yaml
schema: spec-driven
context: |
  Project: BreakLab, a personal break-fix DevOps training CLI.
  One user. Labs are throwaway Docker containers that get broken on
  purpose; the user fixes them and the CLI verifies the fix. Built to be
  lightweight and fast enough to run well on a Raspberry Pi 5, and to be
  reachable from another device on the user's home LAN when deployed
  there headless.

  Stack:
  - Go (latest stable), one module, no other language runtime required
  - CLI built with Go's standard flag/subcommand pattern (or a minimal
    CLI library such as cobra if subcommand ergonomics need it)
  - gopkg.in/yaml.v3 for lab and roadmap files
  - html/template (Go stdlib) + HTMX for both the local web UI and the
    public site; no SPA framework, no JS build step
  - net/http (Go stdlib) for the local web UI; no separate app server
  - Docker is driven by shelling out to the docker CLI via os/exec with
    argument slices (never through a shell)
  - go test for tests; gofmt + go vet (and golangci-lint if configured)
    for lint and format
  - Static site published to GitHub Pages by GitHub Actions
  - Ships as a single statically-linked binary per platform (linux/amd64,
    linux/arm64 for Raspberry Pi 5, darwin/arm64, windows/amd64)

  Conventions:
  - Source in cmd/breaklab/ (CLI entry) and internal/<concern>/, one
    package per concern (catalog, docker, variants, runner, progress,
    postmortem, site, web)
  - Labs live in labs/<lab-id>/ and are data, not code: lab.yaml + scripts
  - All lab scripts are Bash with `set -euo pipefail`
  - Check scripts exit 0 when fixed, non-zero otherwise, and print one line
  - Checks verify the outcome (symptom gone), never a specific fix
  - User data: progress/progress.json and postmortems/*.md (committed)
  - Session state: .breaklab/session.json (gitignored)
  - The CLI must never print the active variant id until the lab is solved
  - The local web UI must never serve variants/, solutions/, or an unused
    hint's text through any route, the same rule as the CLI
  - All docker and script invocations use os/exec with argument slices,
    never a shell (sh -c) or a string-interpolated command
  - Lab ids, variant ids, and any path built from them are validated
    against ^[a-z0-9-]+$ before touching the filesystem or Docker
  - Every command has a --help text and clear error messages
  - Labs must work on Linux, macOS, and Windows via WSL2, and on a
    Raspberry Pi 5 running Raspberry Pi OS or Ubuntu Server (64-bit)
  - The local web UI's bind host defaults to 127.0.0.1 and is overridable
    with `breaklab serve --host <addr>`; binding to any non-loopback
    address prints a one-time warning that the server is unauthenticated
    and reachable by any device on that network

  Out of scope for v1: accounts, cloud VMs, Kubernetes labs, any
  authentication for the local web UI (even when LAN-exposed).
```

## Architecture

BreakLab is one Python CLI that reads lab definitions, drives Docker, and writes your progress into the repo; a local-only web server reads those same files live for browsing, and GitHub Actions turns that progress into a public site when you choose to publish.

&#91;embedded content: BreakLab architecture · 8 components\]

You only ever talk to the CLI, the local web UI on your own machine, and the lab container; everything public flows from files you commit.

Repo layout:

```text
breaklab/
  cmd/breaklab/
    main.go            # CLI entry, subcommand wiring
  internal/
    catalog/catalog.go     # load and validate lab.yaml, roadmap.yaml
    docker/docker.go       # build, run, exec, remove via the docker CLI
    variants/variants.go   # weighted random selection, no repeats
    runner/runner.go       # start, check, hint, reset, stop
    progress/progress.go   # progress.json and streaks
    postmortem/postmortem.go # create and validate postmortems
    site/site.go           # build site/ with html/template (public, publish-triggered)
    web/web.go             # breaklab serve: local HTTP server, reads live files
  labs/              # one folder per lab
  templates/         # postmortem.md, public site HTML and CSS
    web/             # local UI: dashboard, roadmap, challenge browser, lab detail (HTMX fragments)
  roadmap.yaml
  progress/progress.json
  postmortems/
  openspec/          # config.yaml, specs/, changes/
  docs/breaklab-spec.md
  docs/deploy-raspberry-pi.md
  .github/workflows/ # ci.yml, publish.yml
  go.mod
  go.sum
  .gitignore         # .breaklab/, site/
```

## Lab format

A lab is a folder of data and Bash scripts, so adding a lab never requires touching the CLI's Python code. Every lab follows this layout:

```text
labs/web-down/
  lab.yaml             # metadata, variants, hints
  TICKET.md            # the symptom only, never the cause
  Dockerfile           # the healthy environment
  check.sh             # exit 0 = fixed
  variants/
    config-typo.sh     # applies one breakage inside the container
    port-taken.sh
    bad-perms.sh
  solutions/
    config-typo.sh     # reference fix, used only by automated tests
    port-taken.sh
    bad-perms.sh
```

The `lab.yaml` schema:

```yaml
id: web-down                  # must match the folder name
title: The website returns nothing
stage: linux                  # must exist in roadmap.yaml
difficulty: 1                 # 1 to 3
estimated_minutes: 20
tags: [linux, web-servers]    # topic tags, each ^[a-z0-9-]+$, for the
                               # challenge browser's topic filter
environment:
  dockerfile: Dockerfile
  ports: ["8080:80"]          # host:container
check:
  script: check.sh
  run_in: host                # host or container
hints:                        # shown in order, any variant
  - Start with the symptom. What does curl localhost:8080 return?
  - Is the web server process running? Is anything listening on port 80?
variants:
  - id: config-typo
    script: variants/config-typo.sh
    weight: 1
    hints:                    # shown after the general hints
      - Ask nginx to test its own configuration.
  - id: port-taken
    script: variants/port-taken.sh
    weight: 1
    hints:
      - Which process owns port 80? Try ss -ltnp.
```

The roadmap lives in one file at the repo root, and its order sets the order on the progress site:

```yaml
# roadmap.yaml
stages:
  - id: linux
    title: Linux fundamentals
  - id: git
    title: Git
  - id: docker
    title: Docker
  - id: cicd
    title: CI/CD
  - id: iac
    title: Infrastructure as code
  - id: kubernetes
    title: Kubernetes
  - id: observability
    title: Monitoring and logging
```

Script contracts:

- **Variant scripts** run inside the container as root after it starts. They must leave the environment broken and exit 0.
- **Check scripts** print one line explaining the result. They receive `BREAKLAB_CONTAINER` and `BREAKLAB_PORT_<container port> (the mapped host port, for example BREAKLAB_PORT_80=8080)` as environment variables.
- **Solution scripts** apply the reference fix for one variant. The CLI never runs them for the user; only tests do.

## Spec: lab runner and random variants

These requirements become `openspec/specs/lab-runner/spec.md`. When Claude proposes the related changes, the spec deltas should match this text; correct them if they drift.

```markdown
# lab-runner

## Requirements

### Requirement: Lab catalog
The system SHALL load every labs/<id>/lab.yaml and roadmap.yaml at startup
and reject labs whose id, stage, or script paths are invalid.

#### Scenario: Listing labs
- WHEN the user runs `breaklab list`
- THEN labs are shown grouped by roadmap stage, in roadmap order
- AND each row shows id, title, difficulty, and status

#### Scenario: Invalid lab
- WHEN a lab.yaml references a stage missing from roadmap.yaml
- THEN `breaklab list` prints a warning naming the lab and the problem
- AND the other labs still load

### Requirement: Starting a lab
The system SHALL build the lab image, start a fresh container, and apply
exactly one variant chosen by the variant selector.

#### Scenario: Clean start
- WHEN the user runs `breaklab start web-down`
- THEN any previous container for that lab is removed first
- AND a container named breaklab-web-down is started with the lab's ports
- AND the selected variant script runs inside it as root
- AND the ticket text and the `docker exec` command to get a shell are printed
- AND the variant id is NOT printed

#### Scenario: Another lab already active
- WHEN a different lab is active and the user starts a new one
- THEN the CLI asks for confirmation before stopping the active lab

#### Scenario: Docker not running
- WHEN the docker CLI is missing or the daemon is unreachable
- THEN the CLI exits with code 2 and a message explaining how to start Docker

### Requirement: Random variant selection
The system SHALL pick a variant at random, weighted by `weight`, and SHALL
not pick the variant used on that lab's previous attempt when the lab has
two or more variants.

#### Scenario: No immediate repeat
- GIVEN web-down's last attempt used config-typo
- WHEN the user starts web-down again
- THEN the selected variant is not config-typo

#### Scenario: Reproducible selection
- WHEN the user runs `breaklab start web-down --seed 42` twice
- THEN the same variant is selected both times

#### Scenario: Forced variant for tests
- WHEN the BREAKLAB_FORCE_VARIANT environment variable is set
- THEN that variant is used and the selection rules are skipped

### Requirement: Checking a fix
The system SHALL run the lab's check script and report pass or fail with
the script's one-line message.

#### Scenario: Fix works
- WHEN `breaklab check` runs and the check script exits 0
- THEN the CLI prints the message, reveals the variant id, and records a solve
- AND it reminds the user to write the postmortem

#### Scenario: Not fixed yet
- WHEN the check script exits non-zero
- THEN the CLI prints the message and records a failed check attempt

### Requirement: Tiered hints
The system SHALL reveal one hint per `breaklab hint` call: general hints
first, then the active variant's hints, and count how many were used.

#### Scenario: Out of hints
- WHEN every hint has been shown
- THEN the CLI says so and suggests `breaklab reset`

### Requirement: Reset and stop
The system SHALL let the user retry with `breaklab reset` (same variant,
fresh container) and quit with `breaklab stop` (container removed).
```

## Spec: postmortems

A solved lab only becomes "complete" once its postmortem is written, because explaining the fix is what makes it stick and what you will reuse in interviews. These requirements become `openspec/specs/postmortems/spec.md`.

```markdown
# postmortems

## Requirements

### Requirement: Postmortem creation
The system SHALL create a postmortem file from templates/postmortem.md
for the most recently solved attempt and open it in $EDITOR.

#### Scenario: After a solve
- WHEN the user runs `breaklab postmortem` after a passing check
- THEN postmortems/<YYYY-MM-DD>-<lab-id>-<variant-id>.md is created
- AND its front matter holds lab, variant, stage, date, minutes to fix,
  hints used, and publish: false
- AND the file opens in $EDITOR (falling back to nano)

#### Scenario: Nothing solved
- WHEN no attempt has been solved since the last postmortem
- THEN the CLI exits non-zero with "Solve a lab first"

### Requirement: Required sections
Each postmortem SHALL contain these sections, each with at least
20 words of the user's own text: Summary, Detection, Root cause,
Fix, Prevention.

#### Scenario: Incomplete postmortem
- WHEN the user runs `breaklab postmortem --validate` on a file
  with an empty or placeholder Root cause section
- THEN the CLI lists each failing section and exits non-zero

#### Scenario: Valid postmortem
- WHEN every required section passes
- THEN the attempt is marked complete in progress.json

### Requirement: Publishing control
The system SHALL only publish postmortems whose front matter says
publish: true.
```

The template the CLI copies:

```markdown
---
lab: web-down
variant: config-typo
stage: linux
date: 2026-10-17
minutes_to_fix: 18
hints_used: 1
publish: false
---

# Website returns nothing

## Summary
<!-- What broke and what was the impact, in two or three sentences. -->

## Detection
<!-- Which commands and outputs led you to the cause, in order. -->

## Root cause
<!-- The actual cause, not the symptom. -->

## Fix
<!-- What you changed, with the exact commands. -->

## Prevention
<!-- What would stop this in production: a test, an alert, a CI check. -->
```

## Spec: progress tracking and public progress page

Progress is a JSON file committed to the repo, and the public page is a static site built from it plus your published postmortems. These requirements become `openspec/specs/progress/spec.md` and `openspec/specs/progress-site/spec.md`.

```markdown
# progress

## Requirements

### Requirement: Attempt log
The system SHALL append every attempt to progress/progress.json with:
lab id, variant id, started_at, solved_at, failed_checks, hints_used,
postmortem path, and status (active, solved, complete, abandoned).

#### Scenario: Abandoned attempt
- WHEN the user stops or switches labs before solving
- THEN the attempt is saved with status abandoned

### Requirement: Status summary
The system SHALL show progress with `breaklab status`.

#### Scenario: Viewing status
- WHEN the user runs `breaklab status`
- THEN each roadmap stage shows completed labs out of total
- AND the active lab, if any, shows elapsed minutes and hints used
- AND the current daily streak is shown
```

```markdown
# progress-site

## Requirements

### Requirement: Static site build
The system SHALL generate a static site in site/ with `breaklab publish`,
using only progress.json, roadmap.yaml, lab metadata, and postmortems
marked publish: true.

#### Scenario: Building the site
- WHEN the user runs `breaklab publish`
- THEN site/index.html shows every roadmap stage with a progress bar
- AND lists completed labs with date, difficulty, and minutes to fix
- AND links each published postmortem to its own page
- AND shows totals: labs completed, variants seen, current streak

#### Scenario: Nothing leaks
- WHEN a postmortem has publish: false
- THEN no page, link, or text from it appears anywhere in site/
- AND variant scripts, solutions, and hints are never copied into site/

#### Scenario: Works offline
- WHEN site/index.html is opened from disk
- THEN it renders fully with no external CSS, fonts, or scripts

### Requirement: Automated deploy
The repository SHALL include .github/workflows/publish.yml that rebuilds
the site and deploys it to GitHub Pages on every push to main.

#### Scenario: Push to main
- WHEN a commit with a new completed lab is pushed to main
- THEN the workflow runs `breaklab publish` and deploys site/
- AND the public page shows the new lab within one workflow run
```

The `progress.json` shape:

```json
{
  "version": 1,
  "attempts": [
    {
      "lab": "web-down",
      "variant": "config-typo",
      "started_at": "2026-10-17T10:02:00-04:00",
      "solved_at": "2026-10-17T10:20:00-04:00",
      "failed_checks": 2,
      "hints_used": 1,
      "postmortem": "postmortems/2026-10-17-web-down-config-typo.md",
      "status": "complete"
    }
  ]
}
```

## Spec: local web UI

Unlike the public site, this is the view you open every day: it reads
progress.json, roadmap.yaml, lab metadata, and your postmortems live, with no
publish step, and it's the direct answer to "I want roadmap/challenges/a
dashboard like devopslabs.dev, but local and mine." These requirements become
`openspec/specs/web-ui/spec.md`.

```markdown
# web-ui

## Requirements

### Requirement: Local server
The system SHALL serve a local web UI with `breaklab serve [--port] [--host]`,
bound to 127.0.0.1 by default, and SHALL open the user's default browser
to it when bound to a loopback address.

#### Scenario: Starting the server
- WHEN the user runs `breaklab serve`
- THEN the CLI starts a local HTTP server on 127.0.0.1 (default port 8787)
- AND prints the URL
- AND opens the URL in the default browser

#### Scenario: Starting on a headless host
- WHEN the user runs `breaklab serve --host 0.0.0.0` (or any other
  non-loopback address, as on a headless Raspberry Pi)
- THEN the CLI starts the server bound to that address and prints the URL
- AND does NOT attempt to open a local browser

#### Scenario: Port already in use
- WHEN the chosen port is already bound
- THEN the CLI prints a clear error naming the port and exits non-zero,
  rather than silently picking a different one

### Requirement: Dashboard view
The system SHALL render a dashboard showing the current streak, totals
(labs completed, variants seen), the active lab with elapsed time and
hints used if one is running, and the most recent postmortems.

#### Scenario: Dashboard with an active lab
- WHEN a lab is active and the user opens the dashboard
- THEN it shows the lab's elapsed time, hints used, and a reminder of the
  command to check it

#### Scenario: Dashboard with no active lab
- WHEN no lab is active
- THEN the dashboard shows totals and streak only, with a prompt to start
  a lab from the roadmap or challenge browser

### Requirement: Roadmap view
The system SHALL render every stage from roadmap.yaml in order, and every
lab under its stage, each labeled with status (available, in-progress,
complete), difficulty, and estimated minutes.

#### Scenario: Viewing the roadmap
- WHEN the user opens the roadmap page
- THEN stages appear in roadmap.yaml order
- AND each lab under a stage shows status, difficulty, and estimated minutes
- AND a completed lab links to its postmortem if one exists

### Requirement: Challenge browser view
The system SHALL render a browsable list of all labs across every stage,
filterable by stage, difficulty, status, and topic tag.

#### Scenario: Filtering challenges
- WHEN the user filters the challenge browser by stage and difficulty
- THEN only matching labs are shown
- AND clearing the filters restores the full list

#### Scenario: Filtering by topic
- WHEN the user filters the challenge browser by a topic tag
- THEN only labs whose lab.yaml declares that tag are shown
- AND the topic filter's options list every tag used by any lab

### Requirement: Lab detail view
The system SHALL render a lab's ticket text, status, hints used so far,
and the exact CLI command to start or resume it, for any lab the user
selects from the roadmap or challenge browser.

#### Scenario: Opening a lab's detail page
- WHEN the user opens a lab from the roadmap or challenge browser
- THEN the ticket text, status, and hints-used count are shown
- AND the CLI command to start or resume it is shown
- AND no variant id, unused hint text, or any variants/ or solutions/
  content ever appears, whether or not the lab is solved

### Requirement: Postmortem viewer
The system SHALL render any of the user's own completed postmortems in
the local web UI regardless of their publish flag, since that flag only
controls what becomes public.

#### Scenario: Viewing a local postmortem
- WHEN the user opens a completed lab's postmortem from the local UI
- THEN its full content renders, even when publish: false

### Requirement: Configurable bind host, still unauthenticated
The system SHALL bind to 127.0.0.1 unless the user explicitly passes
`--host`, and SHALL NOT require or implement authentication regardless
of which host it binds to.

#### Scenario: Default binding stays local
- WHEN the user runs `breaklab serve` with no `--host` flag
- THEN the server binds to 127.0.0.1 only, matching the non-goal of no
  authentication for the one person running it

#### Scenario: LAN binding is an explicit opt-in
- WHEN the user runs `breaklab serve --host 0.0.0.0` (or a specific
  LAN address), as when deploying to a headless Raspberry Pi
- THEN the server binds to that address
- AND the CLI prints a one-time warning that the UI is unauthenticated
  and reachable by any device that can reach that address
```

## Spec: security audit

BreakLab shells out to Docker, runs lab scripts as root inside throwaway
containers, serves a local web UI, and publishes a public site built from
your own files — enough moving parts that "it's just for me" isn't a reason
to skip a deliberate pass. This audit is a hardening change, not a new
user-facing feature: it reviews and tightens the other seven changes rather
than adding a view or command. These requirements become
`openspec/specs/security/spec.md`.

```markdown
# security

## Requirements

### Requirement: Safe process execution
The system SHALL invoke docker and lab scripts via `os/exec` using
argument slices, and SHALL NOT shell out through `sh -c` or build a
command by string interpolation anywhere in the codebase.

#### Scenario: Auditing exec calls
- WHEN the security audit runs
- THEN every `exec.Command` call under internal/ and cmd/ is checked
- AND any call that goes through a shell or a string-interpolated
  command is flagged and fixed before the audit passes

### Requirement: Input and path validation
The system SHALL validate lab ids, variant ids, and stage ids against
`^[a-z0-9-]+$` before using them to build a filesystem path, a Docker
container/image name, or a URL route, and SHALL reject anything else.

#### Scenario: Path traversal attempt
- GIVEN a lab.yaml, CLI argument, or URL path segment containing `..`,
  `/`, or other characters outside `^[a-z0-9-]+$`
- WHEN the system would use it to build a path or container name
- THEN the system rejects it with a clear error instead of touching the
  filesystem or Docker

### Requirement: No secret leakage
The system SHALL keep session state out of version control and SHALL
only expose postmortems through the public site when their front matter
says publish: true.

#### Scenario: Gitignore coverage
- WHEN the audit checks .gitignore
- THEN .breaklab/ and site/ are present
- AND no progress, postmortem, or session file is tracked that shouldn't be

#### Scenario: Unpublished content stays local
- WHEN `breaklab publish` runs
- THEN no file, filename, or text from a postmortem with publish: false,
  or from any variants/ or solutions/ script, appears anywhere in site/

### Requirement: Dependency audit
CI SHALL run `govulncheck` on every pull request and SHALL fail the
build on any finding that affects code reachable from main.

#### Scenario: Vulnerable dependency introduced
- WHEN a pull request adds or bumps a dependency with a known
  vulnerability reachable from the module's code
- THEN the CI audit step fails and blocks the merge

### Requirement: CI supply-chain hardening
GitHub Actions workflows SHALL pin third-party actions to a full commit
SHA rather than a mutable tag, and SHALL NOT use `pull_request_target`
with untrusted checkout.

#### Scenario: Reviewing workflow files
- WHEN the audit checks ci.yml and publish.yml
- THEN every `uses:` line for a non-GitHub-owned action references a
  commit SHA, not a tag like `@v4`

### Requirement: Web UI attack surface
The local web UI SHALL expose only read routes over the data it already
documents (roadmap, labs, progress, postmortems) and SHALL NOT include
any route that executes a shell command, reads an arbitrary path from a
query parameter, or writes to the filesystem. Binding to a non-loopback
host SHALL require the explicit `--host` flag from the web-ui spec and
SHALL NOT be the default.

#### Scenario: Reviewing web.go routes
- WHEN the audit checks every route in internal/web/web.go
- THEN none of them accept a filesystem path or shell command from the
  request and act on it directly

#### Scenario: LAN exposure requires opt-in
- WHEN the audit checks the default `breaklab serve` invocation
- THEN it binds to 127.0.0.1 unless `--host` was explicitly passed
```

## Change plan

Build in nine changes, in this order; each one is usable on its own, and each is small enough for one Claude session. Tick each off after `/opsx:archive`.

- [ ] **1. `add-lab-catalog`** (spec: lab-runner, Lab catalog)
  - Go module skeleton (`go.mod`), CLI entry point `breaklab` in cmd/breaklab/, gofmt/go vet, go test
  - Load and validate `lab.yaml` and `roadmap.yaml`; `breaklab list`
  - Add the `web-down` lab files (no Docker yet, just valid data)
  - Done when: `breaklab list` shows web-down under Linux fundamentals
- [ ] **2. `add-lab-runner`** (spec: lab-runner, Starting, Checking, Reset and stop)
  - Docker wrapper: build, run, exec, remove, with a clear error when Docker is down
  - `start`, `check`, `reset`, `stop`, plus `.breaklab/session.json`
  - Uses the first variant only for now
  - Done when: you can break and fix web-down end to end
- [ ] **3. `add-random-variants`** (spec: lab-runner, Random variant selection)
  - Weighted selection, no immediate repeat, `--seed`, `BREAKLAB_FORCE_VARIANT`
  - Add the remaining web-down variants and their solution scripts
  - Done when: five starts in a row never repeat back to back
- [ ] **4. `add-tiered-hints`** (spec: lab-runner, Tiered hints)
  - `breaklab hint`, general hints first, then variant hints
  - Done when: hints count increases and runs out cleanly
- [ ] **5. `add-progress-tracking`** (spec: progress)
  - `progress/progress.json`, attempt statuses, `breaklab status`, streaks
  - Done when: status shows a solved lab with minutes and hints
- [ ] **6. `add-postmortems`** (spec: postmortems)
  - Template, `breaklab postmortem`, `--validate`, marks attempts complete
  - Done when: a lab only counts as complete after a valid postmortem
- [ ] **7. `add-progress-site`** (spec: progress-site)
  - Go `html/template` templates, `breaklab publish`, plain CSS, no external assets
  - `publish.yml` GitHub Actions workflow deploying `site/` to GitHub Pages
  - Done when: your public URL shows your progress after a push
- [ ] **8. `add-local-web-ui`** (spec: web-ui)
  - Go `net/http` server in `internal/web/web.go`, `html/template` + HTMX
    templates in `templates/web/`
  - `breaklab serve [--host] [--port]`, bound to 127.0.0.1 by default, opens
    the browser only when bound to a loopback address
  - Dashboard, roadmap, challenge browser, lab detail, and postmortem viewer
  - `docs/deploy-raspberry-pi.md` documents running the compiled binary on a
    Pi 5 with `--host` set for LAN access
  - Depends on changes 1, 5, and 6 (catalog, progress tracking, postmortems)
  - Done when: `breaklab serve` shows live roadmap, dashboard, and challenge
    browser pages, and no variant/solution content is reachable from any route
- [ ] **9. `add-security-audit`** (spec: security)
  - Review every `exec.Command` call for shelling out through `sh -c` or
    string-interpolated commands; fix any found
  - Add id/path validation (`^[a-z0-9-]+$`) at every point a lab, variant,
    or stage id reaches the filesystem, Docker, or a URL route
  - Add a dependency vulnerability scan (`govulncheck`) to ci.yml
  - Pin third-party GitHub Actions to commit SHAs in ci.yml and publish.yml
  - Confirm `breaklab serve` defaults to 127.0.0.1 and only binds to a LAN
    address when `--host` is explicitly passed, with the exposure warning printed
  - Re-check `.gitignore` and the "Nothing leaks" and "Works offline"
    scenarios from progress-site still hold
  - Depends on changes 1 through 8 — it audits what they built, so run it last
  - Done when: every requirement in the security spec passes and CI's
    dependency audit step is green

After v1, good next changes are `add-git-lab`, `add-disk-full-lab`, `add-compose-labs` (multi-container Docker labs), and `add-incident-mode` (two breakages at once with a timer).

## Starter labs

Ship v1 with `web-down` as a tutorial lab, then add the other two as follow-up changes. Because this doc reveals web-down's variants, treat it as practice for the tooling; for every later lab, have Claude write the variants and avoid reading the `variants/` and `solutions/` folders.

| Lab | Stage | Variants | Check passes when |
| --- | --- | --- | --- |
| web-down | Linux | nginx config typo; port 80 taken by another process; homepage file permissions | `localhost:8080` returns 200 with the text "BreakLab OK" |
| disk-full | Linux | huge runaway log; deleted file still held open by a process; inode exhaustion | the app can write a 10 MB file and disk use is under 80% |
| git-lost-work | Git | release commit lost after a hard reset; work done on a detached HEAD; a secret committed then "deleted" | the release commit is on main and the secret is absent from all history |

**web-down files**

`TICKET.md`:

```markdown
# INC-101: Website down after maintenance

Customers report that http://localhost:8080 shows an error or the wrong
page since this morning's maintenance window. Restore the normal homepage.

Get a shell: docker exec -it breaklab-web-down bash
```

`Dockerfile`:

```dockerfile
FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx curl procps iproute2 lsof vim-tiny less python3 \
    && rm -rf /var/lib/apt/lists/*
RUN echo "BreakLab OK" > /var/www/html/index.html
CMD ["sleep", "infinity"]
```

`variants/config-typo.sh` (the other two follow the same pattern):

```bash
#!/usr/bin/env bash
set -euo pipefail
# Drop the semicolon so nginx refuses to start
sed -i 's/listen 80 default_server;/listen 80 default_server/' \
  /etc/nginx/sites-available/default
nginx 2>/dev/null || true
```

`check.sh`:

```bash
#!/usr/bin/env bash
set -uo pipefail
port="${BREAKLAB_PORT_80:-8080}"
body=$(curl -fsS --max-time 3 "http://localhost:${port}/") || {
  echo "FAIL: site unreachable or returned an HTTP error"; exit 1; }
if grep -q "BreakLab OK" <<<"$body"; then
  echo "PASS: site serves the expected page"
else
  echo "FAIL: site is up but serving the wrong content"; exit 1
fi
```

## Prompts for Claude

Use these in Claude Code from the repo root. Export this doc as Markdown, save it as `docs/breaklab-spec.md`, and the prompts can point Claude at it.

**Propose a change** (swap in the change name and number from the Change plan):

```text
/opsx:propose add-lab-catalog

Read docs/breaklab-spec.md. Implement change 1 from the "Change plan"
section only. Use the matching requirements and scenarios from the
"Spec" sections word for word as the spec delta, and the "Lab format"
section for file schemas. Do not add features from later changes.
In design.md, explain each decision in plain language for a DevOps
learner, and list anything in the spec you think is wrong or unclear.
```

**Review the plan before coding:**

```text
Before I run apply, review openspec/changes/add-lab-catalog/ against
docs/breaklab-spec.md. List any requirement or scenario that has no
task, any task not backed by the spec, and anything that would break
on macOS or WSL2. Do not change files yet.
```

**Apply and teach:**

```text
/opsx:apply

After each task, give me a two-line explanation of what you did and
why, so I learn the code as it is written. Run the tests before
ticking a task off.
```

**Write a new lab without spoiling it:**

```text
Create a new lab in labs/<lab-id>/ for the <stage> stage, difficulty
<1-3>, following the "Lab format" section of docs/breaklab-spec.md.
Include three variants, a solution script for each, and an
outcome-based check.sh. Do not describe the variants or causes in
your reply or in TICKET.md; only tell me the lab id, the ticket text,
and that the tests pass.
```

**Debug a lab that misbehaves:**

```text
breaklab check passes even though the site is still broken for
variant <id>. Without telling me the fix for the lab itself, find why
check.sh gives a false pass and propose a stricter check.
```

## Testing and definition of done

The most important test proves every variant really breaks the lab and its reference fix really repairs it; without it, a lab can silently become unsolvable or always pass.

**Lab integration test** (runs for every lab and every variant):

1. Set `BREAKLAB_FORCE_VARIANT=<variant>` and run `breaklab start <lab>`.
2. Assert `breaklab check` fails.
3. Run `solutions/<variant>.sh` inside the container.
4. Assert `breaklab check` passes.
5. Run `breaklab stop` and assert the container is gone.

**Unit tests**

- Catalog: valid labs load; bad stage ids, missing scripts, and duplicate ids are reported.
- Variant selection: seeded runs are reproducible; 1,000 seeded runs never repeat back to back; weights roughly hold.
- Postmortems: placeholder or short sections fail validation; a complete file passes.
- Site: unpublished postmortems and any variant or solution text never appear in `site/`.

**CI** in `.github/workflows/ci.yml` runs `gofmt -l` / `go vet` (and
`golangci-lint` if configured), `go test ./...`, and the lab integration
tests on `ubuntu-latest`, which has Docker available, on every pull
request. It also cross-compiles a `linux/arm64` build as a smoke check
for the Raspberry Pi 5 target (GitHub-hosted runners are amd64-only, so
this is a build check, not a test run, on that architecture).

**Definition of done for each change**

- [ ] `openspec validate <change>` passes
- [ ] Every task in `tasks.md` is ticked and every scenario has a test
- [ ] `gofmt`/`go vet` (and `golangci-lint` if used) and all `go test` pass locally and in CI
- [ ] You ran the feature by hand at least once
- [ ] You can explain the new code in your own words
- [ ] The change is archived with `/opsx:archive` and committed

## Deploying to a Raspberry Pi 5

BreakLab ships as a single statically-linked Go binary, so deploying to
a Pi 5 is: install Docker, copy one binary plus your `labs/`,
`templates/`, and `roadmap.yaml`, and run it with `--host` set so it's
reachable from another device on your LAN. Full step-by-step
instructions live in `docs/deploy-raspberry-pi.md`.
