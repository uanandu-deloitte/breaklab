import { useEffect, useState } from "react";
import { Theme } from "@astryxdesign/core/theme";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { VStack } from "@astryxdesign/core/VStack";
import { Center } from "@astryxdesign/core/Center";
import { ProgressBar } from "@astryxdesign/core/ProgressBar";
import { List, ListItem } from "@astryxdesign/core/List";
import { themes, defaultThemeSlug } from "../astryx-themes";

// Mirrors internal/web's publicProfileView (GET /api/u/{name}).
interface CompletedLab {
  ID: string;
  Title: string;
  Date: string;
  Difficulty: number;
  MinutesToFix: number;
}
interface Stage {
  ID: string;
  Title: string;
  Completed: number;
  Total: number;
  Percent: number;
  Labs: CompletedLab[] | null;
}
interface PublicProfileView {
  Username: string;
  Stages: Stage[] | null;
  Totals: { LabsCompleted: number; VariantsSeen: number; Streak: number };
}

function usernameFromPath(): string {
  if (typeof window === "undefined") return ""; // SSG prerender has no window
  const parts = window.location.pathname.split("/").filter(Boolean);
  return parts[1] ?? ""; // "/u/<name>"
}

// Standalone public page: no Shell/nav — this is a share link for
// visitors who aren't using the app, not part of its own navigation.
function PublicProfileContent() {
  const [view, setView] = useState<PublicProfileView | null>(null);
  const [notFound, setNotFound] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const name = usernameFromPath();
    if (!name) return;
    fetch(`/api/u/${name}`, { credentials: "include" })
      .then((res) => {
        if (res.status === 404) {
          setNotFound(true);
          return null;
        }
        if (!res.ok) throw new Error(`server returned ${res.status}`);
        return res.json();
      })
      .then((data) => data && setView(data))
      .catch((err) => setError(String(err)));
  }, []);

  if (notFound)
    return (
      <Center axis="horizontal" padding={6}>
        <Text color="accent">That profile doesn't exist or isn't public.</Text>
      </Center>
    );
  if (error)
    return (
      <Center axis="horizontal" padding={6}>
        <Text color="accent">Couldn't load this profile: {error}</Text>
      </Center>
    );
  if (!view)
    return (
      <Center axis="horizontal" padding={6}>
        <Text type="supporting">Loading…</Text>
      </Center>
    );

  return (
    <Center axis="horizontal" maxWidth={720} padding={6}>
      <VStack gap={6}>
        <Heading level={1}>{view.Username}</Heading>
        <Text type="supporting">
          {view.Totals.LabsCompleted} labs completed · {view.Totals.VariantsSeen} variants seen · {view.Totals.Streak} day streak
        </Text>
        {(view.Stages ?? []).map((stage) => (
          <VStack gap={2} key={stage.ID}>
            <Heading level={3}>{stage.Title}</Heading>
            <ProgressBar
              label={`${stage.Title} progress`}
              isLabelHidden
              value={stage.Completed}
              max={stage.Total}
              hasValueLabel
            />
            {stage.Labs && stage.Labs.length > 0 && (
              <List hasDividers edge="outset">
                {stage.Labs.map((lab) => (
                  <ListItem key={lab.ID} label={lab.Title} description={`${lab.Date} · difficulty ${lab.Difficulty} · ${lab.MinutesToFix}m`} />
                ))}
              </List>
            )}
          </VStack>
        ))}
      </VStack>
    </Center>
  );
}

export default function PublicProfile() {
  return (
    <Theme theme={themes[defaultThemeSlug]}>
      <PublicProfileContent />
    </Theme>
  );
}
