import { useEffect, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { VStack } from "@astryxdesign/core/VStack";
import { Badge } from "@astryxdesign/core/Badge";
import { Stepper, Step } from "@astryxdesign/core/Stepper";
import { LoadingState } from "./LoadingState";
import Shell from "./Shell";

// Mirrors internal/web/phases's PhaseSummaryView/PhaseListView (GET /api/phases).
interface PhaseSummaryView {
  ID: string;
  Title: string;
  Blurb: string;
  LabsCompleted: number;
  LabsTotal: number;
}

interface PhaseListView {
  Phases: PhaseSummaryView[] | null;
}

function activePhaseIndex(phases: PhaseSummaryView[]): number {
  const index = phases.findIndex((phase) => phase.LabsTotal > 0 && phase.LabsCompleted < phase.LabsTotal);
  if (index !== -1) return index;
  const allComplete = phases.every((phase) => phase.LabsTotal > 0 && phase.LabsCompleted >= phase.LabsTotal);
  return allComplete ? phases.length : 0;
}

function PhasesContent() {
  const [view, setView] = useState<PhaseListView | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch("/api/phases", { credentials: "include" })
      .then((res) => {
        if (!res.ok) throw new Error(`server returned ${res.status}`);
        return res.json();
      })
      .then(setView)
      .catch((err) => setError(String(err)));
  }, []);

  if (error) return <Text color="accent">Couldn't load the phases: {error}</Text>;
  if (!view) return <LoadingState label="Loading phases…" rows={2} />;

  const phases = view.Phases ?? [];
  const activeStep = activePhaseIndex(phases);

  return (
    <VStack gap={4}>
      <Heading level={1}>Phases</Heading>
      <Stepper
        orientation="vertical"
        indicatorPosition="on-track"
        density="spacious"
        activeStep={activeStep}
        label="Phases"
        onStepClick={(index) => {
          const phase = phases[index];
          if (phase && phase.LabsTotal > 0) window.location.href = `/phases/${phase.ID}`;
        }}
      >
        {phases.map((phase, index) => {
          const isEmpty = phase.LabsTotal === 0;
          return (
            <Step
              key={phase.ID}
              step={index}
              label={phase.Title}
              description={phase.Blurb || undefined}
              isDisabled={isEmpty}
              endContent={
                isEmpty ? (
                  <Badge label="Coming soon" variant="neutral" />
                ) : (
                  <>
                    {index === 0 && <Badge label="Start here" variant="info" />}
                    <Text type="supporting">
                      {phase.LabsCompleted}/{phase.LabsTotal} labs
                    </Text>
                  </>
                )
              }
            />
          );
        })}
      </Stepper>
    </VStack>
  );
}

export default function Phases() {
  return (
    <Shell>
      <PhasesContent />
    </Shell>
  );
}
