import type { ReactNode } from "react";
import { Theme } from "@astryxdesign/core/theme";
import { AppShell } from "@astryxdesign/core/AppShell";
import { SideNav, SideNavItem, SideNavHeading } from "@astryxdesign/core/SideNav";
import { themes, defaultThemeSlug } from "../astryx-themes";
import BrandMark from "./BrandMark";

const NAV_ITEMS = [
  { href: "/", label: "Dashboard" },
  { href: "/roadmap", label: "Roadmap" },
  { href: "/challenges", label: "Challenges" },
  { href: "/topology", label: "Topology" },
  { href: "/postmortems", label: "Postmortems" },
  { href: "/chaos", label: "Chaos" },
  { href: "/capstone", label: "Capstone" },
];

function currentPath(): string {
  if (typeof window === "undefined") return ""; // SSG prerender has no window
  return window.location.pathname;
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
            {NAV_ITEMS.map((item) => (
              <SideNavItem key={item.href} label={item.label} href={item.href} isSelected={path === item.href} />
            ))}
          </SideNav>
        }
      >
        {children}
      </AppShell>
    </Theme>
  );
}
