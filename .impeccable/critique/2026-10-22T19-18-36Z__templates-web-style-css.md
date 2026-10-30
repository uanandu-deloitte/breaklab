---
target: templates/web/style.css
total_score: 25
max_score: 36
na_heuristics: 10
p0_count: 2
p1_count: 4
target_identity: "file:/Users/anandu/DevOps-labs/breaklab/templates/web/style.css"
target_fingerprint: "sha256:11f565f022ec1418dac48c9a75b3041792de4f40772800244220f8d3f92c7ca8"
target_path: /Users/anandu/DevOps-labs/breaklab/templates/web/style.css
timestamp: 2026-10-22T19-18-36Z
slug: templates-web-style-css
---
# Design Critique: templates/web/style.css (BreakLab web UI, "Monochrome Minimalist")

Method: dual-agent (A: design-director review · B: Impeccable detector + manual CSS skim)

## Design Health Score
| # | Heuristic | Score | Key Issue |
|---|-----------|-------|-----------|
| 1 | Visibility of System Status | 3/4 | Primary-button hover is a no-op |
| 2 | Match System / Real World | 4/4 | Ticket/postmortem vocabulary maps precisely onto incident response |
| 3 | User Control and Freedom | 3/4 | Stop-lab has no confirm; launch-over-running does |
| 4 | Consistency and Standards | 2/4 | Dead .panel-hero, stale comments, specificity collision |
| 5 | Error Prevention | 3/4 | Flash/check-result states structurally sound |
| 6 | Recognition Rather Than Recall | 4/4 | Ticket ID/stage/hint cost always visible |
| 7 | Flexibility and Efficiency | 1/4 | No shortcuts, no saved filters, no bulk actions |
| 8 | Aesthetic and Minimalist Design | 2/4 | Leftover/dead CSS, flattened type scale |
| 9 | Error Recovery | 3/4 | Clear per-check pass/fail |
| 10 | Help and Documentation | n/a | Not needed, single-user local tool |
| Total | | 25/36 | Acceptable (69%), borderline Good |

## Design Specificity Verdict
High conceptual specificity (incident-response-specific component vocabulary: .incident-banner, .ticket-id, .hint-slot, room-metaphor roadmap, topology node-graph), uneven execution (stale comments, dead class, broken states read as unfinished migration).

## Priority Issues
- [P0] Broken stylesheet path on templates/site/postmortem.html and postmortems.html (href="../style.css" should be "style.css") — 2 public pages render unstyled.
- [P0] --warning collapsed to #000000 makes "in progress"/difficulty-2 indistinguishable from plain text/borders, contradicting the CSS's own "never color alone" claim.
- [P1] Dead .panel-hero class (applied to dashboard's hero panel, zero CSS rule), leftover tan --border-strong (#cfc6ba) flashing on sidebar click, stale comments describing colors/theme that no longer exist.
- [P1] Neobrutalist hard-shadow button treatment reads as "indie/hackathon" regardless of execution quality — worth a deliberate decision given the stated "not self-made" goal.
- [P1] Primary-button :hover is a no-op since --accent == --accent-strong == #000000.
- [P1] .filters button (0,1,1) outranks .btn-primary/.btn-secondary (0,1,0) — padding diverges on the "same" button inside filter forms.
- [P2] --text-sm == --text-base (both 1rem) flattens hierarchy on dense metadata rows.
- [P2] input[type=checkbox]:checked::after font-size: 0.9rem violates the project's own explicit 1rem accessibility floor (introduced this session).
- [P3] Dead CSS (.kv-rows family; .totals/.bar-fill only consumed by templates/site's own separate stylesheet) and DESIGN.md/PRODUCT.md still describe the superseded "Bright Progress"/"Ruled Desk" theme.

## Persona Red Flags
- Alex (Power User): challenges.html exposes up to 9 simultaneous filter inputs with no saved presets for a daily-use tool; inconsistent confirm pattern between stop-lab and launch-over-running.
- Sam (Accessibility): focus ring is pure black against a page already saturated with 3-4px black borders — WCAG-compliant contrast, low practical distinguishability.

## Minor Observations
- Stray 2px border weight on table.rows th / .brand-mark, outside the stated 3px/4px system.
- Unstyled <pre>{{.Ticket}}</pre> inherits dark-terminal treatment meant for real logs/diffs.
- 8 duplicate inline style="margin:0" across 5 templates for the same .panel-header h1/h2 issue.
- --text-sm is now a dead alias.

## Questions to Consider
1. Is "monochrome + 3 status exceptions" the right model, or would full traffic-light semantics serve the incident-training premise better?
2. lab.html renders all 10 sections on load despite the product's symptom-only, progressive-disclosure premise — worth revisiting?
