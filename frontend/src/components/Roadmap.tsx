import { useEffect, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { VStack } from "@astryxdesign/core/VStack";
import { HStack } from "@astryxdesign/core/HStack";
import { Badge } from "@astryxdesign/core/Badge";
import { Button } from "@astryxdesign/core/Button";
import { TreeList } from "@astryxdesign/core/TreeList";
import type { TreeListItemData } from "@astryxdesign/core/TreeList";
import Shell from "./Shell";

// Mirrors internal/web's RoadmapView (GET /api/roadmap).
interface LabRowView {
  ID: string;
  Title: string;
  Summary: string;
  Status: string;
  Difficulty: number;
  EstimatedMinutes: number;
  ReviewDue: boolean;
  IsUpNext: boolean;
}

interface RoadmapStageView {
  ID: string;
  Title: string;
  Blurb: string;
  Index: number;
  Labs: LabRowView[];
  ShowSkip: boolean;
}

interface RoadmapView {
  Stages: RoadmapStageView[];
}

const statusVariant: Record<string, "neutral" | "info" | "success" | "warning"> = {
  available: "info",
  "in-progress": "warning",
  complete: "success",
  locked: "neutral",
};

function labItem(lab: LabRowView): TreeListItemData {
  return {
    id: lab.ID,
    label: lab.Title,
    description: lab.Summary,
    href: `/labs/${lab.ID}`,
    isSelected: lab.IsUpNext,
    endContent: (
      <HStack gap={2} vAlign="center">
        <Badge label={lab.Status} variant={statusVariant[lab.Status] ?? "neutral"} />
        <Badge label={`difficulty ${lab.Difficulty}`} />
        {lab.ReviewDue && <Badge label="Review due" variant="warning" />}
        {lab.IsUpNext && <Badge label="Up next" variant="info" />}
      </HStack>
    ),
  };
}

function stageItem(stage: RoadmapStageView, onSkip: (id: string) => Promise<void>): TreeListItemData {
  const completed = stage.Labs.filter((l) => l.Status === "complete").length;
  return {
    id: stage.ID,
    label: `Stage ${stage.Index} · ${stage.Title}`,
    description: stage.Blurb || undefined,
    isExpanded: true,
    children: stage.Labs.map(labItem),
    endContent: (
      <HStack gap={2} vAlign="center">
        <Text type="supporting">
          {completed}/{stage.Labs.length}
        </Text>
        {stage.ShowSkip && (
          <Button
            label="Already know this — skip"
            variant="ghost"
            size="sm"
            clickAction={() => onSkip(stage.ID)}
          />
        )}
      </HStack>
    ),
  };
}

function RoadmapContent() {
  const [view, setView] = useState<RoadmapView | null>(null);
  const [error, setError] = useState<string | null>(null);

  function load() {
    fetch("/api/roadmap", { credentials: "include" })
      .then((res) => {
        if (!res.ok) throw new Error(`server returned ${res.status}`);
        return res.json();
      })
      .then(setView)
      .catch((err) => setError(String(err)));
  }

  useEffect(load, []);

  async function handleSkip(stageID: string) {
    const res = await fetch(`/api/roadmap/stages/${stageID}/skip`, {
      method: "POST",
      credentials: "include",
    });
    if (res.ok) load();
  }

  if (error) return <Text color="accent">Couldn't load the roadmap: {error}</Text>;
  if (!view) return <Text type="supporting">Loading…</Text>;

  return (
    <VStack gap={4}>
      <Heading level={1}>Roadmap</Heading>
      <TreeList
        header={<Text type="body">Curriculum roadmap, by stage</Text>}
        variant="lineGuides"
        density="spacious"
        items={view.Stages.map((stage) => stageItem(stage, handleSkip))}
      />
    </VStack>
  );
}

export default function Roadmap() {
  return (
    <Shell>
      <RoadmapContent />
    </Shell>
  );
}
