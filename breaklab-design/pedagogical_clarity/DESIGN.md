---
name: Pedagogical Clarity
colors:
  surface: '#131315'
  surface-dim: '#131315'
  surface-bright: '#39393b'
  surface-container-lowest: '#0e0e10'
  surface-container-low: '#1c1b1d'
  surface-container: '#201f22'
  surface-container-high: '#2a2a2c'
  surface-container-highest: '#353437'
  on-surface: '#e5e1e4'
  on-surface-variant: '#c4c7c9'
  inverse-surface: '#e5e1e4'
  inverse-on-surface: '#313032'
  outline: '#8e9193'
  outline-variant: '#444749'
  surface-tint: '#c6c6c7'
  primary: '#ffffff'
  on-primary: '#2f3132'
  primary-container: '#e2e2e3'
  on-primary-container: '#636466'
  inverse-primary: '#5d5e60'
  secondary: '#4edea3'
  on-secondary: '#003824'
  secondary-container: '#00a572'
  on-secondary-container: '#00311f'
  tertiary: '#ffffff'
  on-tertiary: '#00354a'
  tertiary-container: '#c4e7ff'
  on-tertiary-container: '#006c93'
  error: '#ffb4ab'
  on-error: '#690005'
  error-container: '#93000a'
  on-error-container: '#ffdad6'
  primary-fixed: '#e2e2e3'
  primary-fixed-dim: '#c6c6c7'
  on-primary-fixed: '#1a1c1d'
  on-primary-fixed-variant: '#454748'
  secondary-fixed: '#6ffbbe'
  secondary-fixed-dim: '#4edea3'
  on-secondary-fixed: '#002113'
  on-secondary-fixed-variant: '#005236'
  tertiary-fixed: '#c4e7ff'
  tertiary-fixed-dim: '#7bd0ff'
  on-tertiary-fixed: '#001e2c'
  on-tertiary-fixed-variant: '#004c69'
  background: '#131315'
  on-background: '#e5e1e4'
  surface-variant: '#353437'
typography:
  headline-xl:
    fontFamily: Geist
    fontSize: 40px
    fontWeight: '600'
    lineHeight: 48px
    letterSpacing: -0.03em
  headline-xl-mobile:
    fontFamily: Geist
    fontSize: 30px
    fontWeight: '600'
    lineHeight: 38px
    letterSpacing: -0.025em
  headline-lg:
    fontFamily: Geist
    fontSize: 32px
    fontWeight: '600'
    lineHeight: 40px
    letterSpacing: -0.02em
  headline-lg-mobile:
    fontFamily: Geist
    fontSize: 24px
    fontWeight: '600'
    lineHeight: 32px
    letterSpacing: -0.015em
  headline-md:
    fontFamily: Geist
    fontSize: 24px
    fontWeight: '500'
    lineHeight: 32px
    letterSpacing: -0.015em
  headline-sm:
    fontFamily: Geist
    fontSize: 20px
    fontWeight: '500'
    lineHeight: 28px
    letterSpacing: -0.01em
  body-lg:
    fontFamily: Geist
    fontSize: 17px
    fontWeight: '400'
    lineHeight: 28px
    letterSpacing: -0.005em
  body-md:
    fontFamily: Geist
    fontSize: 15px
    fontWeight: '400'
    lineHeight: 24px
    letterSpacing: 0em
  body-sm:
    fontFamily: Geist
    fontSize: 13px
    fontWeight: '400'
    lineHeight: 20px
    letterSpacing: 0.005em
  label-md:
    fontFamily: JetBrains Mono
    fontSize: 13px
    fontWeight: '500'
    lineHeight: 18px
    letterSpacing: 0.02em
  label-sm:
    fontFamily: JetBrains Mono
    fontSize: 11px
    fontWeight: '500'
    lineHeight: 16px
    letterSpacing: 0.04em
  code-block:
    fontFamily: JetBrains Mono
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 22px
    letterSpacing: 0em
rounded:
  sm: 0.125rem
  DEFAULT: 0.25rem
  md: 0.375rem
  lg: 0.5rem
  xl: 0.75rem
  full: 9999px
spacing:
  gutter: 1.5rem
  gutter-sm: 1rem
  margin: 2rem
  margin-sm: 1rem
  margin-lg: 3rem
  space-xs: 0.25rem
  space-sm: 0.5rem
  space-md: 1rem
  space-lg: 1.5rem
  space-xl: 2.5rem
---

## Brand & Style

This design system is engineered for deep focus, long-duration technical learning, and zero cognitive fatigue. Tailored for engineers, researchers, and technical students, the visual language merges the discipline of high-end developer tooling with the legibility of academic printing.

The aesthetic is rooted in **Precision Functional Minimalism**:
- **Monochrome Foundation:** The interface relies on an absolute black-and-white tonal hierarchy using calibrated zinc and slate neutrals. Surfaces prioritize ink-on-paper readability over decorative ornament.
- **Surgical Accentuation:** Color is treated as functional telemetry. It is strictly excluded from general branding or decorative fills, reserved solely for execution states, syntax tokens, interactive indicators, and status badges.
- **Restraint & Dignity:** Clear structural hierarchy, ample breathing room, and strict alignment prevent visual crowding during dense technical lectures, terminal simulations, or code breakdowns.

## Colors

The palette establishes an eye-friendly, deeply focused dark foundation (`#09090B`) that dramatically reduces visual strain in low-light environments, paired with an inverted high-contrast scale meeting WCAG AAA requirements (7:1+ contrast ratio for primary text).

### Canvas & Surface Architecture
- **Canvas Base (`#09090B`):** Deep zinc floor for the outermost background viewport.
- **Surface Layer 1 (`#121215`):** Structural panels, sidebars, and course navigation trees.
- **Surface Layer 2 (`#18181B`):** Cards, terminal shells, quiz modules, and code editors.
- **Surface Elevated (`#27272A`):** Popovers, dropdown menus, and command palettes.
- **Borders (`#27272A` / `#3F3F46` active):** 1px structural definition lines separating content planes without aggressive visual cuts.

### Typography Scales
- **Foreground Primary (`#F4F4F5`):** Headers, active reading blocks, and critical metadata (contrast ratio > 15:1 against canvas).
- **Foreground Secondary (`#A1A1AA`):** Body context, chapter outlines, and structural labels.
- **Foreground Muted (`#71717A`):** Captions, timestamps, inactive pagination, and breadcrumbs.

### Surgical Semantic Accents
Accents are used at small scale (pills, badges, dots, focus rings, status outlines, and syntax nodes):
- **Emerald (`#10B981`):** Complete, verified, passed assessments, operational status.
- **Amber (`#F59E0B`):** In-progress units, queued checks, cautionary notices.
- **Rose (`#F43F5E`):** Broken builds, syntax errors, failed test suites, attention alerts.
- **Sky (`#38BDF8`):** Informational highlights, active focus states, active code parameters, and contextual tooltips.

## Typography

The typographic hierarchy pairs **Geist** for prose, navigation, and module headers with **JetBrains Mono** for status flags, metadata tags, metric readouts, and terminal/source code displays.

- **Prose Readability:** Long-form pedagogical content uses `body-lg` (17px on 28px line height) constrained to a maximum reading width of 68 characters (~720px) to prevent saccadic reading errors and eye tracking fatigue.
- **Technical Discipline:** All interactive keys, shortcuts (`⌘K`), tabular metrics, step counters, and completion percentages are set strictly in `label-sm` or `label-md` using tabular figures (`tnum`).
- **Heading Restraint:** Headings favor modest scale ratios and tight tracking over massive display banners, preserving vertical real estate for learning materials and schematics.

## Layout & Spacing

The layout is structured around an asymmetrical split-pane and fluid grid architecture designed for simultaneous reading, interactive execution, and reference navigation.

### Structural Framework
- **Master Shell:** 3-column setup on desktop (collapsible course sidebar: 280px; main reading document: max 780px centered; interactive sandbox/terminal pane: fluid 420px–640px).
- **Responsive Adaptations:**
  - **Desktop (>=1280px):** 12-column grid, 24px gutters, 48px outer canvas padding. Sidebar and sandbox can be pinned side-by-side.
  - **Tablet (768px - 1279px):** 8-column grid, 16px gutters, 24px outer margins. The auxiliary sandbox docks into a persistent bottom sheet or tabbed mode.
  - **Mobile (<768px):** 4-column grid, 16px margins, single-stack linear progression with fixed top navigation and bottom-anchored control trays.
- **Rhythm:** Layout spaces strictly honor an 8px base rhythm (`0.25rem` to `2.5rem`), ensuring consistent line-height locks across code editors and explanatory markdown.

## Elevation & Depth

This system avoids blurred colored drop shadows and skeuomorphic lighting in favor of **structural planar hierarchy** and **calibrated 1px boundary lines**.

1. **Floor (Level 0):** Background canvas (`#09090B`). Flat, absorbing, non-reflective.
2. **Surface Framing (Level 1):** Main documentation cards, execution panels, and lesson containers sit on `#121215` with a continuous `1px solid #27272A` perimeter.
3. **Interactive & Inset Focus (Level 2):** Input fields, terminal interiors, and code blocks adopt a sunken tonal contrast (`#0C0C0E`) bordered by `#27272A`.
4. **Overlays & Command Interfaces (Level 3):** Modals, search palettes, and tooltips use `#18181B` enclosed in `1px solid #3F3F46`, paired with an ultra-subtle non-tinted atmospheric falloff (`box-shadow: 0 12px 32px -8px rgba(0, 0, 0, 0.7)`).
5. **State Activation:** Hover and active layers do not elevate via Z-index or shadow expansion; they elevate strictly through border luminescence (transitioning from `#27272A` to `#52525B`) or high-contrast focus rings (`2px solid #38BDF8`).

## Shapes

The design system enforces a **Soft (Level 1)** geometric standard. Rounding is restrained, functional, and consistent with precision scientific instruments and developer IDEs.

- **Base Radius (`0.25rem` / 4px):** Applied to badges, code tokens, list items, inline chips, status pills, and input controls.
- **Container Radius (`0.5rem` / 8px):** Applied to code blocks, interactive sandboxes, module cards, and modal dialogs.
- **Inner-to-Outer Conformity:** When components are nested (e.g., an action button within a card header), inner elements strictly use 4px while the enclosing container uses 8px, maintaining geometric concentricity.
- **Circular Primitives:** Full rounding (`rounded-full`) is reserved solely for progress circular rings and atomic status presence indicators (e.g., 6px execution indicator dots).

## Components

### Buttons
- **Primary:** High-contrast solid zinc fill (`#F4F4F5`), text in `#09090B`, 4px radius, weight 500. On hover, background shifts to `#E4E4E7`.
- **Secondary / Outline:** Background transparent, 1px solid border `#3F3F46`, text `#F4F4F5`. On hover, background fills to `#18181B` with border `#71717A`.
- **Tertiary / Ghost:** No border, text `#A1A1AA`. On hover, text `#F4F4F5` and background `#18181B`.
- **Destructive:** Border `1px solid #F43F5E`, text `#F43F5E`, background transparent. Hover transitions to `#F43F5E` fill with `#FFFFFF` text.

### Badges & Semantic Status Indicators
- Constructed using `label-sm` font, uppercase or tabular formatting, 4px radius, and `0.25rem 0.5rem` padding.
- **Completed:** Border `1px solid #065F46`, background `#064E3B`/20%, text `#10B981`. Includes a 6px solid `#10B981` dot prefix.
- **In-Progress:** Border `1px solid #92400E`, background `#78350F`/20%, text `#F59E0B`.
- **Failed / Attention:** Border `1px solid #9F1239`, background `#881337`/20%, text `#F43F5E`.
- **Info / Metadata:** Border `1px solid #075985`, background `#0C4A6E`/20%, text `#38BDF8`.

### Input Fields & Terminal Prompts
- Background `#0C0C0E`, 1px solid `#27272A`, 4px radius, text `#F4F4F5`, placeholder `#52525B`.
- **Focus:** No soft halo; transitions to a sharp `1px solid #38BDF8` with an accompanying keyboard navigation outline ring (`2px offset`).
- **Code Inputs / Command Prompts:** Include a fixed leading terminal prompt glyph (`$` or `>`) in `#71717A` with monospace formatting.

### Checkboxes & Radio Controls
- Square (4px radius for checkbox) and circle (radio), sized at 16px by 16px.
- Unchecked: `1px solid #3F3F46`, surface `#121215`.
- Checked: Surface `#F4F4F5` with checkmark icon `#09090B` for general tasks; semantic completion assessments use `#10B981` fill with `#09090B` glyph.

### Cards & Code Playgrounds
- Background `#121215`, 1px solid `#27272A`, 8px radius.
- Headers are delineated by a 1px bottom border `#1F1F23`, featuring tabbed file titles in monospace (`label-md`) and integrated execution badges.
- Content margins inside containers follow strict 16px/24px padding rhythms.

### Interactive Module Lists
- Borderless rows separated by `1px solid #1F1F23`.
- Hover introduces a unified row highlight of `#18181B`.
- Active lesson item displays a left border indicator (`2px solid #38BDF8` or `#F4F4F5`) and text highlighted to `#FFFFFF`.