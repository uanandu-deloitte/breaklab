// A high-contrast black-and-white theme, extending @astryxdesign/theme-neutral
// with overridden background/text/border tokens. Modeled on the reference
// screenshots: true black page background in dark mode, true white in light
// mode, hairline borders instead of shadows for card separation. Accent is
// Deloitte's brand green (the two tones swap by mode so the on-accent label
// always has AA contrast: the darker green on light backgrounds needs white
// text, the brighter green on dark backgrounds needs black text).
import {defineTheme} from '@astryxdesign/core/theme';
import {neutralTheme} from '@astryxdesign/theme-neutral';

const systemSans =
  '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif';

export const monoTheme = defineTheme({
  name: 'mono',
  extends: neutralTheme,
  // Override neutral's Figtree with the system stack: no font file to load
  // or warn about, and a plainer, more "dev tool" look fits the stark
  // black-and-white direction better than a rounded display font.
  typography: {
    scale: {base: 14, ratio: 1.2},
    body: {family: 'system-ui', fallbacks: systemSans},
    heading: {family: 'system-ui', fallbacks: systemSans, weights: {3: 'bold', 4: 'bold'}},
    code: {
      family: 'ui-monospace',
      fallbacks: '"SF Mono", Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
    },
  },
  tokens: {
    // Deloitte Green (Pantone 376C, ~#86BC25) for dark surfaces; the
    // deeper Deloitte Green Dark (~#046A38) for light surfaces, where the
    // brighter tone alone would be too low-contrast for body text and UI.
    '--color-accent': ['#046A38', '#86BC25'],
    '--color-on-accent': ['#FFFFFF', '#000000'],
    '--color-text-accent': ['#046A38', '#86BC25'],
    '--color-icon-accent': ['#046A38', '#86BC25'],
    '--color-accent-muted': ['#046A3826', '#86BC2533'],
    '--focus-outline-color': 'var(--color-accent)',

    '--color-background-body': ['#FFFFFF', '#000000'],
    '--color-background-surface': ['#FFFFFF', '#0A0A0A'],
    '--color-background-card': ['#FFFFFF', '#0A0A0A'],
    '--color-background-popover': ['#FFFFFF', '#141414'],
    '--color-background-muted': ['#0000000F', '#FFFFFF0F'],

    '--color-text-primary': ['#000000', '#FFFFFF'],
    '--color-text-secondary': ['#525252', '#A3A3A3'],
    '--color-text-disabled': ['#A3A3A3', '#666666'],
    '--color-icon-primary': ['#000000', '#FFFFFF'],
    '--color-icon-secondary': ['#525252', '#A3A3A3'],
    '--color-icon-disabled': ['#A3A3A3', '#666666'],

    // Alpha raised from the original 12%/20-24% hairlines (~1.3:1 contrast)
    // to clear the WCAG 1.4.11 non-text 3:1 floor: this is the sole depth
    // cue for every Card (no shadows in this theme) and the Stepper's
    // connector line, so it has to actually be visible to do its job.
    '--color-border': ['#00000073', '#FFFFFF73'],
    '--color-border-emphasized': ['#00000099', '#FFFFFF99'],
    '--color-skeleton': ['#D4D4D4', '#4D4D4D'],
    '--color-track': ['#D4D4D4', '#4D4D4D'],
  },
});
