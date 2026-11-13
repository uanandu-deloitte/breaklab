---
version: 1
slug: "frontend-src"
primary_target: "frontend/src"
related_targets: []
---

## Direction contract

THESIS: BreakLab's dashboard reads as a live incident-investigation case file, not a course-platform scoreboard — it refuses the card-grid-plus-progress-bar SaaS-education default (KodeKloud / devopslabs / Codecademy) in favor of the same discipline a real investigator's report demands: exhibits, dispositions, and findings closed only once reviewed.

OWN-WORLD: Official report-paper ground (warm off-white), black ink body copy, one authoritative red for review marks and open-finding stamps, deep navy for exhibit numbering and reference marks. Serif report typeface (Source Serif 4) for headings and body text — the register of an investigation document, not UI chrome; JetBrains Mono (carried over unchanged) for every exhibit number, case ID, measurement, and code block. A rotated case-review stamp anchors each page's header. Every lab is an exhibit, numbered; every status is a disposition (OPEN / UNDER REVIEW / CLOSED), not a cheerful badge.

STORY: the engineer opens the tool and reads it like reviewing their own open casework — which exhibits are still open, which are under review, which are closed and why — then acts: resume the open exhibit or pull the next one from the docket.

FIRST VIEWPORT: Dashboard. Report-paper background. Left rail (brand mark + nav) in ink-on-paper, active nav item marked by weight plus a small case-tab, never a thick colored border. Case header across the top of the main area — case number and page title left, a rotated review stamp (red/navy) pinned to the outer frame's top-right corner. Three measured finding plots in a row directly below the header (exhibits closed, day streak, open findings). Below that, an exhibit-numbered table of labs: stage, difficulty, disposition tag, each row carrying its own exhibit number (EX-114 style). The primary action (resume / start next lab) lives inside each exhibit row, at the row's right edge — not a separate floating CTA; the docket itself is the action surface.

FORM: Incident Record (NTSB-style investigation report). IMPECCABLE'S PICK — candidate 1 of 7 on the grounded list, the most literal reading of the product's own postmortem structure (Summary/Detection/Root cause/Fix/Prevention already mirror a report's own findings). Chosen by the user over the dice-assigned Site Survey (cyanotype blueprint) direction. seed:98cc1d8593b1a3f8

FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance.
