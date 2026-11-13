import { useEffect, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { VStack } from "@astryxdesign/core/VStack";
import { HStack } from "@astryxdesign/core/HStack";
import { Badge } from "@astryxdesign/core/Badge";
import { Button } from "@astryxdesign/core/Button";
import { Card } from "@astryxdesign/core/Card";
import { Link } from "@astryxdesign/core/Link";
import { ProgressBar } from "@astryxdesign/core/ProgressBar";
import { EmptyState } from "@astryxdesign/core/EmptyState";
import { Stepper, Step } from "@astryxdesign/core/Stepper";
import { LabCard, labHref } from "./LabCard";
import { LoadingState } from "./LoadingState";
import Shell from "./Shell";

// Mirrors internal/web/phases's PhaseDetailView (GET /api/phases/{id}),
// nesting internal/web/roadmap's RoadmapStageView/LabRowView.
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
  // Non-empty only when this stage is one of a set of alternatives (see
  // catalog.Stage.Track) — a learner follows one, not all of them.
  Track: string;
}

interface PhaseDetailView {
  ID: string;
  Title: string;
  Blurb: string;
  Stages: RoadmapStageView[] | null;
}

// A roadmap entry is either one ordinary stage, or a run of adjacent
// stages sharing a Track — alternatives rendered as one branching node.
// Mirrors catalog.GroupByTrack (internal/catalog/track.go) client-side.
type RoadmapEntry = { kind: "stage"; stage: RoadmapStageView } | { kind: "track"; track: string; stages: RoadmapStageView[] };

function groupByTrack(stages: RoadmapStageView[]): RoadmapEntry[] {
  const entries: RoadmapEntry[] = [];
  let i = 0;
  while (i < stages.length) {
    const track = stages[i].Track;
    if (!track) {
      entries.push({ kind: "stage", stage: stages[i] });
      i += 1;
      continue;
    }
    let j = i + 1;
    while (j < stages.length && stages[j].Track === track) j += 1;
    entries.push({ kind: "track", track, stages: stages.slice(i, j) });
    i = j;
  }
  return entries;
}

function entryKey(entry: RoadmapEntry): string {
  return entry.kind === "stage" ? entry.stage.ID : `track:${entry.track}`;
}

function findUpNext(stages: RoadmapStageView[]): LabRowView | null {
  for (const stage of stages) {
    for (const lab of stage.Labs) {
      if (lab.IsUpNext) return lab;
    }
  }
  return null;
}

// Maps the flat "which stage is up next" position onto its containing
// entry, so a stage inside a track group still advances the Stepper's
// shared activeStep rather than needing its own index space.
function activeEntryIndex(entries: RoadmapEntry[], stages: RoadmapStageView[]): number {
  const upNextStageID = stages.find((stage) => stage.Labs.some((lab) => lab.IsUpNext))?.ID;
  if (upNextStageID) {
    const index = entries.findIndex((entry) =>
      entry.kind === "stage" ? entry.stage.ID === upNextStageID : entry.stages.some((s) => s.ID === upNextStageID),
    );
    if (index !== -1) return index;
  }
  const allComplete = stages.every((stage) => stage.Labs.every((lab) => lab.Status === "complete"));
  return allComplete ? entries.length : 0;
}

function SkipButton({ onSkip }: { onSkip: () => void }) {
  return <Button label="Already know this — skip" variant="ghost" size="sm" clickAction={onSkip} />;
}

function LabCardList({ labs, backHref, backLabel }: { labs: LabRowView[]; backHref: string; backLabel: string }) {
  return (
    <>
      {labs.map((lab) => (
        <LabCard key={lab.ID} lab={lab} backHref={backHref} backLabel={backLabel} />
      ))}
    </>
  );
}

function phaseIDFromPath(): string {
  if (typeof window === "undefined") return ""; // SSG prerender has no window
  const parts = window.location.pathname.split("/").filter(Boolean);
  return parts[1] ?? ""; // "/phases/<id>"
}

function PhaseDetailContent() {
  const [phaseID, setPhaseID] = useState("");
  const [view, setView] = useState<PhaseDetailView | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [notFound, setNotFound] = useState(false);
  // A set, not a single id, because a track group and one of its member
  // stages are both "expanded" at once when drilling in two levels deep.
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set());

  function toggleExpanded(id: string) {
    setExpandedIds((current) => {
      const next = new Set(current);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }

  useEffect(() => {
    setPhaseID(phaseIDFromPath());
  }, []);

  function load() {
    if (!phaseID) return;
    fetch(`/api/phases/${phaseID}`, { credentials: "include" })
      .then((res) => {
        if (res.status === 404) {
          setNotFound(true);
          return null;
        }
        if (!res.ok) throw new Error(`server returned ${res.status}`);
        return res.json();
      })
      .then((data) => {
        if (data) setView(data);
      })
      .catch((err) => setError(String(err)));
  }

  useEffect(load, [phaseID]);

  async function handleSkip(stageID: string) {
    const res = await fetch(`/api/phases/stages/${stageID}/skip`, {
      method: "POST",
      credentials: "include",
    });
    if (res.ok) load();
  }

  if (notFound) return <Text color="accent">That phase doesn't exist.</Text>;
  if (error) return <Text color="accent">Couldn't load this phase: {error}</Text>;
  if (!view) return <LoadingState label="Loading phase…" rows={3} />;

  const stages = view.Stages ?? [];
  const entries = groupByTrack(stages);
  const labsTotal = stages.reduce((sum, s) => sum + s.Labs.length, 0);
  const labsCompleted = stages.reduce((sum, s) => sum + s.Labs.filter((l) => l.Status === "complete").length, 0);
  const upNext = findUpNext(stages);
  const activeStep = activeEntryIndex(entries, stages);
  const backHref = `/phases/${phaseID}`;
  const backLabel = `Back to ${view.Title}`;

  return (
    <VStack gap={6}>
      <VStack gap={2}>
        <Heading level={1}>{view.Title}</Heading>
        {view.Blurb && <Text type="supporting">{view.Blurb}</Text>}
        {labsTotal > 0 && (
          <ProgressBar
            label={`${view.Title} progress`}
            value={labsCompleted}
            max={labsTotal}
            hasValueLabel
            formatValueLabel={() => `${labsCompleted}/${labsTotal} labs complete`}
          />
        )}
      </VStack>

      {upNext && (
        <Card variant="blue">
          <HStack gap={2} justify="between" vAlign="center">
            <VStack gap={1}>
              <Text type="supporting" color="primary">
                Continue where you left off
              </Text>
              <Heading level={3}>
                <Link href={labHref(upNext.ID, backHref, backLabel)}>{upNext.Title}</Link>
              </Heading>
            </VStack>
          </HStack>
        </Card>
      )}

      {stages.length === 0 ? (
        <EmptyState
          title="No stages yet"
          description="This phase doesn't have any curriculum built out yet. Check back once it does."
        />
      ) : (
        <Stepper
          orientation="vertical"
          indicatorPosition="on-track"
          density="spacious"
          activeStep={activeStep}
          label={`${view.Title}'s stages`}
          onStepClick={(index) => {
            const entry = entries[index];
            if (entry) toggleExpanded(entryKey(entry));
          }}
        >
          {entries.map((entry, index) => {
            if (entry.kind === "stage") {
              const stage = entry.stage;
              const isExpanded = expandedIds.has(stage.ID);
              return (
                <Step
                  key={stage.ID}
                  step={index}
                  label={stage.Title}
                  description={stage.Blurb || undefined}
                  endContent={
                    <HStack gap={2} vAlign="center">
                      {index === 0 && <Badge label="Start here" variant="info" />}
                      <Text type="supporting">{stage.Labs.length} labs</Text>
                      <Text type="supporting">|</Text>
                      <Text type="label" color="accent">
                        {isExpanded ? "Collapse" : "Expand"}
                      </Text>
                    </HStack>
                  }
                >
                  {isExpanded && (
                    <VStack gap={3}>
                      <HStack justify="between" vAlign="center">
                        <Text type="label" color="secondary">
                          Labs
                        </Text>
                        {stage.ShowSkip && <SkipButton onSkip={() => handleSkip(stage.ID)} />}
                      </HStack>
                      <LabCardList labs={stage.Labs} backHref={backHref} backLabel={backLabel} />
                    </VStack>
                  )}
                </Step>
              );
            }

            const groupKey = entryKey(entry);
            const isExpanded = expandedIds.has(groupKey);
            const totalLabs = entry.stages.reduce((sum, s) => sum + s.Labs.length, 0);
            return (
              <Step
                key={groupKey}
                step={index}
                label={entry.track}
                description={`Pick one of these ${entry.stages.length} paths — finishing (or skipping) it unlocks the rest of the roadmap.`}
                endContent={
                  <HStack gap={2} vAlign="center">
                    {index === 0 && <Badge label="Start here" variant="info" />}
                    <Text type="supporting">{totalLabs} labs</Text>
                    <Text type="supporting">|</Text>
                    <Text type="label" color="accent">
                      {isExpanded ? "Collapse" : "Expand"}
                    </Text>
                  </HStack>
                }
              >
                {isExpanded && (
                  <VStack gap={3}>
                    {entry.stages.map((member) => {
                      const memberExpanded = expandedIds.has(member.ID);
                      const completed = member.Labs.filter((l) => l.Status === "complete").length;
                      return (
                        <Card key={member.ID}>
                          <VStack gap={2}>
                            <HStack justify="between" vAlign="start">
                              <VStack gap={1}>
                                <Heading level={4}>{member.Title}</Heading>
                                {member.Blurb && <Text type="supporting">{member.Blurb}</Text>}
                              </VStack>
                              <Text type="supporting">
                                {completed}/{member.Labs.length} labs
                              </Text>
                            </HStack>
                            <HStack gap={2} vAlign="center">
                              {member.ShowSkip && <SkipButton onSkip={() => handleSkip(member.ID)} />}
                              <Button
                                label={memberExpanded ? "Collapse" : "Expand"}
                                variant="ghost"
                                size="sm"
                                clickAction={() => toggleExpanded(member.ID)}
                              />
                            </HStack>
                            {memberExpanded && (
                              <VStack gap={3}>
                                <LabCardList labs={member.Labs} backHref={backHref} backLabel={backLabel} />
                              </VStack>
                            )}
                          </VStack>
                        </Card>
                      );
                    })}
                  </VStack>
                )}
              </Step>
            );
          })}
        </Stepper>
      )}
    </VStack>
  );
}

export default function PhaseDetail() {
  return (
    <Shell>
      <PhaseDetailContent />
    </Shell>
  );
}
