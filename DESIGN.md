---
name: BreakLab
description: A personal, offline-first DevOps break/fix training ground — a warm paper-cream dashboard, cards outlined in solid ink, for tracking hands-on incident-response practice.
colors:
  bg: "#f2ecdc"
  bg-elevated: "#ffffff"
  bg-sunken: "#e9e0c8"
  bg-overlay: "#ded2b0"
  fg: "#1a1711"
  fg-muted: "#615948"
  fg-subtle: "#9c9480"
  border: "#e2d8bc"
  border-strong: "#1a1711"
  accent: "#c14e1f"
  accent-strong: "#963812"
  accent-subtle: "#f7e2d2"
  ok: "#2f6b2f"
  warning: "#8c5a07"
  danger: "#ab2e2e"
  info: "#1d5f7a"
  pastel-mint: "#dbe4ca"
  pastel-lavender: "#e7dcc0"
  pastel-sky: "#ded7c2"
  pastel-butter: "#efdfb0"
  pastel-blush: "#f3d9c4"
  terminal-bg: "#18181b"
  terminal-fg: "#f4f4f5"
  terminal-muted: "#b3b3b3"
  diff-add: "#4ade80"
  diff-remove: "#f87171"
  diff-ctx: "#9a9a9a"
typography:
  h1:
    fontFamily: "Geist, ui-sans-serif, system-ui, -apple-system, Segoe UI, sans-serif"
    fontSize: "1.8rem"
    fontWeight: 800
    lineHeight: 1.2
    letterSpacing: "-0.03em"
  h2:
    fontFamily: "Geist, ui-sans-serif, system-ui, -apple-system, Segoe UI, sans-serif"
    fontSize: "1.125rem"
    fontWeight: 700
    lineHeight: 1.2
    letterSpacing: "-0.02em"
  body:
    fontFamily: "Geist, ui-sans-serif, system-ui, -apple-system, Segoe UI, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.5
  label:
    fontFamily: "JetBrains Mono, ui-monospace, SFMono-Regular, Consolas, monospace"
    fontSize: "1rem"
    fontWeight: 500
    lineHeight: 1.5
rounded:
  sm: "10px"
  md: "16px"
  lg: "22px"
  full: "999px"
spacing:
  xs: "0.25rem"
  sm: "0.5rem"
  md: "1rem"
  lg: "1.5rem"
  xl: "2.5rem"
components:
  button-primary:
    backgroundColor: "{colors.fg}"
    textColor: "#ffffff"
    rounded: "{rounded.full}"
    padding: "0 1.75rem"
  button-primary-hover:
    backgroundColor: "{colors.accent-strong}"
  button-secondary:
    backgroundColor: "{colors.pastel-lavender}"
    textColor: "{colors.fg}"
    rounded: "{rounded.full}"
    padding: "0 1.75rem"
  button-secondary-hover:
    backgroundColor: "{colors.accent-subtle}"
  badge:
    textColor: "{colors.fg-muted}"
    rounded: "{rounded.full}"
    padding: "0.2rem 1rem"
  card:
    backgroundColor: "{colors.bg-elevated}"
    rounded: "{rounded.md}"
---

# Design System: BreakLab

## Overview

**Creative North Star: "Warm Ledger"**

BreakLab trains one person at a time to work real incidents: a container breaks, a ticket shows only the symptom, and the user debugs with no instructions. The interface's job is to make that loop feel like a game worth returning to daily, not like staring at a pager at 3am. "Warm Ledger" is the resolution of that tension — a warm paper-cream canvas where every card is a crisp white (or lightly tanned) rectangle outlined in solid ink, cycling through four muted warm fills so a list of labs or roadmap stages reads as varied rather than a wall of identical gray boxes. But the "ledger" half of the name is a real constraint: status colors (ok/warning/danger/info) stay distinct, legible, and never purely decorative — a lab that's locked, in progress, or complete has to be instantly scannable even inside all that warmth, because the practice only works if the state is always trustworthy.

This replaces the earlier **"Bright Signal"** system (a light lavender-gray canvas, one saturated violet accent, soft shadow-only card definition, candy-pastel fills) — itself a successor to a cool-neutral zinc "Quiet Indigo" direction and, before that, a dark monochrome "Pedagogical Clarity" direction, both now historical reference material under `breaklab-design/`. Warm Ledger is a deliberate identity replacement pinned to a reference image, not an incremental refinement — including an explicit reversal of Bright Signal's own stated "not the generic warm-cream default" call: the canvas went from lavender-gray to paper-cream, the one accent went from violet to terracotta, card definition moved from a soft ambient shadow to a real ink-colored border, and the four pastel fills were re-dyed from candy hues to muted warm neutrals (sage/sand/clay/gold). Shape (pill buttons, rounded cards) and the near-black "console" sidebar carry over unchanged — this is a re-skin of color and surface treatment, not a structural rebuild.

**Key Characteristics:**
- A warm paper-cream canvas (`#f2ecdc`), replacing the lavender-gray tint
- A near-black console surface (`#1a1711`, now warm-black rather than cool-black) for the sidebar and primary buttons — the one serious, structural element against an otherwise warm page
- Every card (`.panel`, `.active-lab`, `.stat-tile`, `.stage-node`) is outlined in a real 1.5px solid ink border (`--border-strong`) instead of floating on a soft shadow — the reference's signature "white card, black line" look
- Four pastel card fills that rotate across repeating lists (stat tiles, lab cards, roadmap stage ribbons, room-mark chips), re-dyed warm/muted so color still carries variety without reading as candy
- Pill radius (`999px`) on every badge, button, tag, and search field; `16–22px` on cards
- Status color is never decorative-only: every badge/difficulty pairs its hue with a text label and a leading mark (dot or small emoji), and the status hues were re-picked to stay legible as *text* sitting on their own pastel tint, not just as a 3:1 accent against the warm page background
- A dashed connector line (with a circular arrow "bead") between roadmap stages, echoing the dashed underline on every `<h1>` — the one recurring decorative motif, used exactly twice

## Colors

The palette is a warm, paper-cream base carrying one saturated terracotta accent, four legible-but-muted status hues, and four warm pastel fills reserved strictly for card backgrounds.

### Primary
- **Signal Terracotta** (`#c14e1f`): the one actionable/link hue — the active sidebar item's icon tint, focus rings, `.btn-secondary`'s hover border, and solid-fill UI components (`.sidebar-active`, `.brand-mark`) where it's a background behind white text, not text on a card. Nothing else in the system uses this hue, so a terracotta element is always "click me." **Text specifically** — every `<a>`, both its default and `:hover` state — uses the darker **Accent Strong** (`#963812`) instead: `#c14e1f` only clears 4.5:1 against white, failing on every pastel card fill and even `--bg` itself, so it's reserved for 3:1 UI-component contexts (borders, rings, icon/button fills) and never used as body text color. `a:hover` adds only `text-decoration: underline`, no color change, since Accent Strong is already the hover-safe color at rest.

### Secondary
- **Console Ink** (`#1a1711`): the near-black chrome color — the sidebar's background and `.btn-primary`'s fill, and now also every card's outline (`--border-strong`). It's deliberately *not* pure accent-colored; a black sidebar against a warm page reads as a stable dock, not a second competing brand color.

### Tertiary — Pastel Fills (card backgrounds only, never text)
- **Sage** (`#dbe4ca`), **Sand** (`#e7dcc0`), **Clay** (`#ded7c2`), **Gold** (`#efdfb0`), **Blush** (`#f3d9c4`, reserved for the one "hero" card per page — the active-incident banner, the lab-detail `panel-hero`). The first four rotate in a fixed 4-color cycle across any repeating list of cards; Blush never rotates in — it's the one loud card a page is allowed to have.

### Neutral
- **Ink** (`#1a1711`): primary text.
- **Dim Ink** (`#615948`): secondary/muted text — darkened specifically so it still clears 4.5:1 sitting on a pastel card, not only against the page background.
- **Faint Ink** (`#9c9480`): tertiary text (meta labels, placeholders).
- **Hairline** (`#e2d8bc`): internal dividers and quiet borders. **Hairline Strong** (`#1a1711`, same value as Ink) is now the card-outline color, not a darker tint of Hairline — see the Elevation & Depth section below.
- **Paper** (`#ffffff`) / **Sunken Paper** (`#e9e0c8`) / **Pressed Paper** (`#ded2b0`): elevated surfaces, recessed surfaces (table headers, code inline), and pressed/active states.

### Status
- **Signal Green** (`#2f6b2f`), **Signal Amber** (`#8c5a07`), **Signal Rose** (`#ab2e2e`), **Signal Teal** (`#1d5f7a`): ok/warning/danger/info. Each one drives both a badge's border/text color *and*, via `color-mix(currentColor, white)`, its own pastel fill — so the badge is self-tinting, never a separately-authored color pair.

### Terminal Surface (dark, one-off)
- **Terminal Ink** (`#18181b`), **Terminal Paper** (`#f4f4f5`), **Terminal Dim** (`#b3b3b3`): the background/text/muted-text set for the app's one deliberate dark "technical readout" surface (`.log-tail`, the topology graph canvas wrapper) — not a second theme, just this one component family's own backdrop. The topology graph itself (Cytoscape.js, canvas-rendered) takes its node/edge colors straight from the Primary/Status hues above (`#c14e1f`/`#1d5f7a`/`#2f6b2f` as literal values in `topology.js`'s style array, since a canvas library can't consume CSS custom properties) rather than a separate dark-surface palette.

### Named Rules
**The Self-Tinting Badge Rule.** A status badge's background is always `color-mix(in srgb, currentColor 22%, var(--bg-elevated))` — never a second, independently chosen pastel. One hue produces both the text color and the fill, so they can never drift out of sync, and the hue has to be dark enough to stay legible on its *own* 22%-tint, not just against the page.

**The One Hero Card Rule.** `--pastel-blush` never enters the four-color rotation; it is reserved for exactly one "this matters most right now" card per page (the active-incident banner, a hero panel). Everywhere else, boldness comes from rotation and shape, not from a louder color.

## Typography

**Display/Body Font:** Geist (with `ui-sans-serif, system-ui` fallback)
**Label/Mono Font:** JetBrains Mono (with `ui-monospace` fallback)

**Character:** Geist carries every heading, label, and paragraph — a single rounded-geometric sans doing all the talking, so the "fun" register comes from weight (800 on `<h1>`) and color, never from mixing in a second display face. JetBrains Mono is reserved for things that are literally code, data, or IDs (ticket numbers, shell commands, log output, diffs), so its appearance always reads as "real system output."

### Hierarchy
- **H1** (800, 1.8rem, 1.2 line-height, −0.03em tracking): page titles only, one per page, underlined with a 3px dashed rule instead of a solid one.
- **H2** (700, 1.125rem, 1.2 line-height, −0.02em tracking): section headings inside a page (Task, Hints, Workspace, ...).
- **Body** (400, 1rem/16px, 1.5 line-height): all prose. 1rem is a hard floor — nothing in the app ever renders smaller, carried over unchanged from every prior system as a standing accessibility commitment.
- **Label/Mono** (500, 1rem, JetBrains Mono): ticket IDs, shell commands, hint costs, log tails.

### Named Rules
**The No-Shrinking Rule.** `--text-sm` equals `--text-base` (both 1rem). There is no smaller type size anywhere in the system; a "small" label is marked by weight or color, never by dropping under 16px.

## Layout

A persistent left console (15.5rem, collapsing to a 3.75rem icon rail under 48rem) plus a single scrollable content column, max-width 75rem, 2rem/1.5rem outer padding (narrower 54rem on the postmortem long-form reader). Repeating content uses two patterns: a `.stat-grid` (auto-fit columns, 11rem minimum) for metric tiles, and a `.split` two-column grid (2fr/1fr, collapsing to one column under 56rem) for a primary panel beside a secondary sidebar panel (diagnostics beside hints/launch, activity beside a sandbox feed). Spacing runs on a loose rem scale (0.25 / 0.5 / 1 / 1.5 / 2.5rem) rather than a tight 8px grid — cards get generous internal padding (1–1.5rem) to support the bigger radius.

## Elevation & Depth

Line-leaning, not shadow-leaning: every card (`.panel`, `.active-lab`, `.stat-tile`, `.stage-node`) gets a real 1.5px solid `--border-strong` (ink) outline instead of a soft ambient shadow — depth comes from the pastel fill plus that crisp line, echoing the reference's black-bordered white/tan cards. The one exception is the sidebar itself, which keeps its real directional shadow (`4px 0 16px rgba(0,0,0,0.12)`) to read as a floating dock in front of the page rather than a panel embedded in it — the sidebar is chrome, not a card, so it's exempt from the ink-outline rule.

### Border/Shadow Vocabulary
- **Card outline** (`1.5px solid var(--border-strong)`): the default on every card, panel, stat tile, and stage node — this replaced the prior system's ambient shadow as the primary depth cue.
- **Hero lift** (`0 4px 6px -1px rgba(0,0,0,0.06), 0 2px 4px -2px rgba(0,0,0,0.06)`): the one shadow upgrade left in the system, reserved for `.panel-hero` (paired with the Blush fill) — a card can have a line *or* a lift, Blush gets both.
- **Console dock** (`4px 0 16px rgba(0,0,0,0.12)`): the sidebar only — the single strongest shadow in the system, and the only shadow most pages show at all.

### Named Rules
**The Line-By-Default Rule.** Cards don't lean on shadow for definition; the ink outline does that job uniformly, so no card reads as "more important" purely from elevation. A card only gets a shadow on top of its outline when it's also wearing Blush (`.panel-hero`).

## Shapes

Pill-round is the signature: `--radius-full` (999px) on every button, badge, tag, hairline pill (search input, loop-stepper links, hint slots), and the sidebar's own nav-item highlight. Cards step down from there — `16px` default, `22px` for a handful of larger hero-ish containers (`.active-lab`, `.stage-node`) — never square, never the old system's crisp 6–12px. The one deliberate non-pill, non-rounded motif is the dashed line: a 3px dashed rule under every `<h1>` and a 3px dashed vertical connector between roadmap stages, each terminating in a small pill "bead" (the heading has none; the connector's bead is a circular terracotta arrow button).

### Named Rules
**The No Side-Tab Rule.** Don't add a thick, flat-colored border on one edge of a rounded card as a stand-alone "accent" — it's the single most recognizable AI-generated-UI tell, and on a card with real border-radius it visibly clashes at the corner where the thick line meets the curve. If a section needs to announce what it is, say so in its heading or its badge, not with a colored bar. (This replaced an earlier left-border-per-section treatment on the lab page's Task/Hints/Launch blocks, removed after both a design-detector pass and a direct code review flagged it.)

## Components

Rounded and tactile: buttons are full pills you can't mistake for anything but a button, cards are crisply ink-bordered rather than soft-shadowed, and badges are small colored capsules rather than square outlined tags.

### Buttons
- **Shape:** full pill (`border-radius: 999px`), `2.75rem` tall, `1.75rem` horizontal padding.
- **Primary:** Console Ink fill (`#1a1711`), white text — the one black pill on an otherwise warm page, so "the main action" is unambiguous. Hover shifts to Signal Terracotta, not a darker black.
- **Secondary:** Pastel Sand fill, Ink text, transparent border. Hover shifts to Accent Subtle with a visible Signal Terracotta border.
- **Hover / Focus:** background-color and shadow transition (0.15s), plus a 1px press-down translate on `:active`. Focus-visible gets a 2px solid Signal Terracotta outline, 2px offset, on every focusable element — one rule, no per-component overrides.

### Chips / Badges
- **Style:** full pill, `0.2rem / 1rem` padding, 700-weight text, 1px border at 25% of the status hue, background from the Self-Tinting Badge Rule.
- **State:** `badge-available` (Dim Ink, neutral dot), `badge-in-progress` (Signal Amber, ⏳), `badge-complete` (Signal Green, 🎉), `badge-locked` (Signal Rose, 🔒), `badge-review-due` (Signal Teal, 🔁). The emoji replaces the plain dot only where it adds real scannability; difficulty badges keep the plain dot — a difficulty level isn't an event worth an emoji.

### Cards / Containers
- **Corner Style:** 16px default, 22px for `.active-lab` / `.stage-node`.
- **Background:** rotates through Sage → Sand → Clay → Gold across any repeating list (`:nth-child`-driven); Blush is reserved for the one hero card per page.
- **Shadow Strategy:** none by default — a 1.5px solid ink (`--border-strong`) outline carries definition instead; see Elevation & Depth.
- **Border:** 1.5px solid Ink on every card, regardless of fill color — not a lighter Hairline tint.
- **Internal Padding:** 1–1.5rem.

### Inputs / Fields
- **Style:** 1px Hairline border, 10px radius (full pill for the search bar specifically), white background, JetBrains Mono for anything code-flavored.
- **Focus:** border shifts to Signal Terracotta plus a 2px Accent Subtle ring — no glow/blur.
- **Labeling:** every input/textarea carries a programmatic name — a visible `<label>` where there's room, `aria-label` where the layout is a compact inline search/textarea paired with its own section heading.

### Navigation (Sidebar / Console)
- **Style:** near-black console (`#1a1711`), floating with rounded right corners and a real directional shadow — the one surface in the system that isn't warm-pastel.
- **Default state:** white text/icon at full opacity (never faded to a muted gray on a dark surface — that reads as disabled, not "just not current").
- **Active state:** a solid white pill (`.sidebar-nav-current`), Ink text — the "chip lifted off a dark dock" pattern.
- **Signature element:** `.sidebar-active` — the one Signal Terracotta–filled card in the whole sidebar, reserved for "a lab is running right now."
- **Mobile:** collapses to a 3.75rem icon-only rail under 48rem; labels and meta text hide, icons and the live-pulse dot stay.

### Roadmap Connector (signature component)
A 3px dashed vertical line running behind each stage transition, with the connecting blurb (if any) and a small circular Accent-Subtle "bead" with a terracotta arrow floating on top of the dash (the line is interrupted visually by the bead's own opaque background, not literally cut). It's the one place in the system where a dashed line *is* the UI, not just a heading underline.

## Do's and Don'ts

### Do:
- **Do** rotate card fills from the fixed four-pastel set (Sage/Sand/Clay/Gold) across any repeating list — never invent a fifth ad hoc pastel, and never use Blush outside the one-hero-per-page rule.
- **Do** keep every status hue dark enough to pass 4.5:1 as *text on its own tint*, not just 3:1 against the page background — verify any new status color against `color-mix(currentColor, white)`, not only `--bg`.
- **Do** pair every status/difficulty signal with a text label and a leading mark (dot or emoji); color never carries meaning alone.
- **Do** use the dashed-line motif sparingly — it's the `<h1>` underline and the roadmap connector, not a general-purpose divider.

### Don't:
- **Don't** add a thick, flat-colored border on one side of a rounded card as a stand-alone accent (see The No Side-Tab Rule) — say what a section is in its heading, not in a colored bar.
- **Don't** reach for opacity to convey a "locked/disabled" state on text — it silently erodes contrast depending on what's underneath; set an explicit muted color instead.
- **Don't** give Console Ink (the sidebar/primary-button black) a second job as a decorative brand color elsewhere on the page — it stays reserved for structural chrome so it keeps reading as "the dock," not "another accent."
- **Don't** rely on placeholder text as a form control's only accessible name — add a visible `<label>` or `aria-label`.
