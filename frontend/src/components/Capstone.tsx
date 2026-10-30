import { useEffect, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { VStack } from "@astryxdesign/core/VStack";
import { ProgressBar } from "@astryxdesign/core/ProgressBar";
import Shell from "./Shell";
import { LabDetailContent } from "./LabDetail";

// Mirrors internal/web's CapstoneView (GET /api/capstone).
interface StageStatus {
  ID: string;
  Title: string;
  Completed: number;
  Total: number;
}
interface CapstoneView {
  Found: boolean;
  Unlocked: boolean;
  Progress: StageStatus[] | null;
  Completed: number;
  Total: number;
  Lab: { ID: string };
}

function LockedProgress({ view }: { view: CapstoneView }) {
  return (
    <VStack gap={4}>
      <Heading level={1}>Capstone</Heading>
      <Text>Finish every stage below to unlock it.</Text>
      <ProgressBar label="Overall capstone progress" value={view.Completed} max={view.Total} hasValueLabel variant="accent" />
      {(view.Progress ?? []).map((stage) => (
        <VStack gap={1} key={stage.ID}>
          <Text type="supporting">{stage.Title}</Text>
          <ProgressBar
            label={`${stage.Title} progress`}
            isLabelHidden
            value={stage.Completed}
            max={stage.Total}
            hasValueLabel
          />
        </VStack>
      ))}
    </VStack>
  );
}

export default function Capstone() {
  const [view, setView] = useState<CapstoneView | null>(null);
  const [notFound, setNotFound] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch("/api/capstone", { credentials: "include" })
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
      <Shell>
        <Text color="accent">No capstone is available yet.</Text>
      </Shell>
    );
  if (error)
    return (
      <Shell>
        <Text color="accent">Couldn't load the capstone: {error}</Text>
      </Shell>
    );
  if (!view)
    return (
      <Shell>
        <Text type="supporting">Loading…</Text>
      </Shell>
    );

  // The unlocked lab brings its own Layout and padding (LabDetailContent,
  // shared with /labs/{id}); the locked progress view is plain content
  // that still needs Shell's own inset.
  return (
    <Shell contentPadding={view.Unlocked ? 0 : 6}>
      {view.Unlocked ? <LabDetailContent fixedLabID={view.Lab.ID} /> : <LockedProgress view={view} />}
    </Shell>
  );
}
