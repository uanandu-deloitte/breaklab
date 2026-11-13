import type { ReactNode } from "react";
import { Theme } from "@astryxdesign/core/theme";
import { AppShell } from "@astryxdesign/core/AppShell";
import { SideNav, SideNavItem, SideNavHeading, SideNavSection } from "@astryxdesign/core/SideNav";
import { themes, defaultThemeSlug } from "../astryx-themes";
import BrandMark from "./BrandMark";

// One small outline glyph per destination (16x16, 1.5px stroke, matching the
// design system's own fallback-icon style) — the built-in icon registry has
// nothing that reads as "roadmap" or "random incident", so these are drawn
// by hand rather than reached-for-but-wrong stand-ins.
function iconProps(extra?: Record<string, unknown>) {
  return {
    width: 16,
    height: 16,
    viewBox: "0 0 24 24",
    fill: "none",
    stroke: "currentColor",
    strokeWidth: 1.5,
    strokeLinecap: "round" as const,
    strokeLinejoin: "round" as const,
    "aria-hidden": true,
    ...extra,
  };
}

function DashboardIcon() {
  return (
    <svg {...iconProps()}>
      <rect x="3" y="3" width="8" height="8" rx="1.5" />
      <rect x="13" y="3" width="8" height="5" rx="1.5" />
      <rect x="13" y="10" width="8" height="11" rx="1.5" />
      <rect x="3" y="13" width="8" height="8" rx="1.5" />
    </svg>
  );
}

function PhasesIcon() {
  return (
    <svg {...iconProps()}>
      <path d="M4 20V14M12 20V9M20 20V4" />
      <circle cx="4" cy="12" r="1.75" fill="currentColor" stroke="none" />
      <circle cx="12" cy="7" r="1.75" fill="currentColor" stroke="none" />
      <circle cx="20" cy="2" r="1.75" fill="currentColor" stroke="none" />
    </svg>
  );
}

function ChallengesIcon() {
  return (
    <svg {...iconProps()}>
      <circle cx="12" cy="12" r="8" />
      <circle cx="12" cy="12" r="3.5" />
    </svg>
  );
}

function TopologyIcon() {
  return (
    <svg {...iconProps()}>
      <circle cx="12" cy="4.5" r="2" />
      <circle cx="5" cy="18" r="2" />
      <circle cx="19" cy="18" r="2" />
      <path d="M12 6.5V12M12 12 6.3 16.6M12 12l5.7 4.6" />
    </svg>
  );
}

function PostmortemsIcon() {
  return (
    <svg {...iconProps()}>
      <path d="M6 3h9l4 4v14a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1Z" />
      <path d="M9 11h7M9 15h7M9 7h3" />
    </svg>
  );
}

function ChaosIcon() {
  return (
    <svg {...iconProps()}>
      <path d="M13 2 4 14h6l-1 8 9-12h-6z" strokeLinejoin="round" fill="none" />
    </svg>
  );
}

function CapstoneIcon() {
  return (
    <svg {...iconProps()}>
      <path d="M5 3v18" />
      <path d="M5 4h12l-3 4 3 4H5" strokeLinejoin="round" />
    </svg>
  );
}

const NAV_SECTIONS: { title: string; items: { href: string; label: string; icon: ReactNode }[] }[] = [
  {
    title: "Practice",
    items: [
      { href: "/", label: "Dashboard", icon: <DashboardIcon /> },
      { href: "/phases", label: "Phases", icon: <PhasesIcon /> },
      { href: "/challenges", label: "Challenges", icon: <ChallengesIcon /> },
      { href: "/chaos", label: "Chaos", icon: <ChaosIcon /> },
    ],
  },
  {
    title: "Review",
    items: [
      { href: "/topology", label: "Topology", icon: <TopologyIcon /> },
      { href: "/postmortems", label: "Postmortems", icon: <PostmortemsIcon /> },
      { href: "/capstone", label: "Capstone", icon: <CapstoneIcon /> },
    ],
  },
];

function currentPath(): string {
  if (typeof window === "undefined") return ""; // SSG prerender has no window
  return window.location.pathname;
}

// A drilled-down page (/phases/phase-1, /postmortems/some-slug) should still
// light up its section in the nav — exact-match alone left the sidebar
// blank on every detail page, the one place wayfinding matters most.
function isActive(href: string, path: string): boolean {
  if (href === "/") return path === "/";
  return path === href || path.startsWith(`${href}/`);
}

// Shell wraps every page in the app's theme, side nav, and content
// padding. contentPadding=0 is for pages (lab detail, capstone) that
// build their own Layout with its own padding, matching Astryx's
// "AppShell full-bleed + nested Layout keeps its own inset" pattern.
export default function Shell({
  children,
  contentPadding = 6,
}: {
  children: ReactNode;
  contentPadding?: 0 | 0.5 | 1 | 1.5 | 2 | 3 | 4 | 5 | 6 | 8 | 10;
}) {
  const path = currentPath();
  return (
    <Theme theme={themes[defaultThemeSlug]}>
      <AppShell
        contentPadding={contentPadding}
        sideNav={
          <SideNav header={<SideNavHeading icon={<BrandMark />} heading="BreakLab" headingHref="/" />}>
            {NAV_SECTIONS.map((section) => (
              <SideNavSection key={section.title} title={section.title}>
                {section.items.map((item) => (
                  <SideNavItem
                    key={item.href}
                    label={item.label}
                    href={item.href}
                    icon={item.icon}
                    isSelected={isActive(item.href, path)}
                  />
                ))}
              </SideNavSection>
            ))}
          </SideNav>
        }
      >
        {children}
      </AppShell>
    </Theme>
  );
}
