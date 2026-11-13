import { useEffect, useMemo, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { Card } from "@astryxdesign/core/Card";
import { ClickableCard } from "@astryxdesign/core/ClickableCard";
import { CodeBlock } from "@astryxdesign/core/CodeBlock";
import { Grid } from "@astryxdesign/core/Grid";
import { VStack } from "@astryxdesign/core/VStack";
import { HStack } from "@astryxdesign/core/HStack";
import { Badge } from "@astryxdesign/core/Badge";
import { Link } from "@astryxdesign/core/Link";
import { Banner } from "@astryxdesign/core/Banner";
import { TabList, Tab } from "@astryxdesign/core/TabList";
import { ProgressBar } from "@astryxdesign/core/ProgressBar";
import { EmptyState } from "@astryxdesign/core/EmptyState";
import { VisuallyHidden } from "@astryxdesign/core/VisuallyHidden";
import { LoadingState } from "./LoadingState";
import Shell from "./Shell";

// Mirrors internal/web's DashboardView (GET /api/dashboard) field for
// field. Sidebar is omitted server-side; everything else is the exact
// data dashboardHandler used to render server-side.
interface ActiveLabView {
  Lab: string;
  Title: string;
  Href: string;
  TicketID: string;
  ElapsedMinutes: number;
  HintsUsed: number;
  HintsTotal: number;
  StageTitle: string;
  ShellCmd: string;
  ReminderCmd: string;
}

interface RecentPostmortemView {
  Lab: string;
  TicketID: string;
  Date: string;
  Href: string;
}

// Mirrors internal/web/dashboard's StageCardView.
interface StageCardView {
  ID: string;
  Title: string;
  PhaseID: string;
  LabsCompleted: number;
  LabsTotal: number;
  Complete: boolean;
}

// Mirrors internal/web/phases's PhaseSummaryView.
interface PhaseSummaryView {
  ID: string;
  Title: string;
  Blurb: string;
  LabsCompleted: number;
  LabsTotal: number;
}

interface DashboardView {
  Streak: number;
  LabsCompleted: number;
  VariantsSeen: number;
  Active: ActiveLabView | null;
  Recent: RecentPostmortemView[] | null;
  Stages: StageCardView[] | null;
  Phases: PhaseSummaryView[] | null;
  MonthlyEngagement: Record<string, number> | null;
  MostEngagedDay: string | null;
  DockerVersion: string;
  Uptime: string;
  AvgResolutionTime: string;
  PostmortemDrafts: number;
  PostmortemPublished: number;
}

function StatCard({ label, value }: { label: string; value: string | number }) {
  return (
    <Card>
      <VStack gap={1}>
        <Text type="supporting">{label}</Text>
        <Heading level={2}>{String(value)}</Heading>
      </VStack>
    </Card>
  );
}

function ActiveLabHero({ active }: { active: ActiveLabView }) {
  return (
    <Card variant="blue">
      <VStack gap={2}>
        <Text type="supporting">Continue where you left off</Text>
        <HStack gap={2} vAlign="center">
          <Heading level={3}>
            <Link href={active.Href}>
              {active.Lab}: {active.Title}
            </Link>
          </Heading>
          <Badge label={active.TicketID} variant="info" />
        </HStack>
        <Text type="supporting">
          {active.ElapsedMinutes}m elapsed · {active.HintsUsed}/{active.HintsTotal} hints used
        </Text>
        {active.ShellCmd && (
          <VStack gap={2}>
            <Heading level={4}>Get a shell</Heading>
            <CodeBlock code={active.ShellCmd} language="bash" />
          </VStack>
        )}
        <CodeBlock code={active.ReminderCmd} language="bash" />
      </VStack>
    </Card>
  );
}

function StageCard({ stage }: { stage: StageCardView }) {
  return (
    <ClickableCard label={stage.Title} href={`/phases/${stage.PhaseID}`}>
      <VStack gap={2}>
        <Heading level={4}>{stage.Title}</Heading>
        <ProgressBar
          label={`${stage.Title} progress`}
          value={stage.LabsCompleted}
          max={Math.max(stage.LabsTotal, 1)}
          hasValueLabel
          formatValueLabel={() => `${stage.LabsCompleted}/${stage.LabsTotal} labs`}
          variant={stage.Complete ? "success" : "accent"}
        />
      </VStack>
    </ClickableCard>
  );
}

const DAY_FORMAT = new Intl.DateTimeFormat(undefined, { month: "long", day: "numeric" });

function engagementVariant(count: number): "muted" | "blue" | "cyan" {
  if (count <= 0) return "muted";
  if (count <= 2) return "blue";
  return "cyan";
}

function MonthlyEngagement({
  engagement,
  mostEngagedDay,
}: {
  engagement: Record<string, number>;
  mostEngagedDay: string | null;
}) {
  const now = useMemo(() => new Date(), []);
  const days = useMemo(() => {
    const year = now.getFullYear();
    const month = now.getMonth();
    const daysInMonth = new Date(year, month + 1, 0).getDate();
    return Array.from({ length: daysInMonth }, (_, i) => {
      const day = i + 1;
      const key = `${year}-${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
      return { key, day, count: engagement[key] ?? 0 };
    });
  }, [now, engagement]);

  return (
    <VStack gap={2}>
      <HStack gap={2} vAlign="center">
        <Heading level={3}>Monthly engagement</Heading>
      </HStack>
      {/* Fixed container width + rowHeight makes each 1fr column track
          exactly as tall as it is wide, so every cell below (which
          stretches to fill its track, Grid's default) renders as an
          actual square instead of a short rectangle. */}
      <Grid columns={7} gap={1} width={136} rowHeight={16}>
        {days.map((d) => (
          <Card key={d.key} variant={engagementVariant(d.count)} padding={0}>
            <VisuallyHidden>
              {d.key}: {d.count} lesson{d.count === 1 ? "" : "s"} completed
            </VisuallyHidden>
          </Card>
        ))}
      </Grid>
      {mostEngagedDay ? (
        <Text type="supporting">
          Most engaged day: {DAY_FORMAT.format(new Date(mostEngagedDay + "T00:00:00"))} ·{" "}
          {engagement[mostEngagedDay]} lesson{engagement[mostEngagedDay] === 1 ? "" : "s"} completed
        </Text>
      ) : (
        <Text type="supporting">No activity yet this month.</Text>
      )}
    </VStack>
  );
}

function DashboardContent() {
  const [view, setView] = useState<DashboardView | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [tab, setTab] = useState<"in-progress" | "completed">("in-progress");

  useEffect(() => {
    fetch("/api/dashboard", { credentials: "include" })
      .then((res) => {
        if (!res.ok) throw new Error(`server returned ${res.status}`);
        return res.json();
      })
      .then(setView)
      .catch((err) => setError(String(err)));
  }, []);

  if (error) return <Text color="accent">Couldn't load the dashboard: {error}</Text>;
  if (!view) return <LoadingState label="Loading dashboard…" rows={3} />;

  const stages = view.Stages ?? [];
  const visibleStages = stages.filter((s) =>
    tab === "completed" ? s.Complete : !s.Complete && s.LabsCompleted > 0,
  );
  const emptyPhases = (view.Phases ?? []).filter((p) => p.LabsTotal === 0);

  return (
    <VStack gap={6}>
      <Heading level={1}>Dashboard</Heading>

      {view.Active ? (
        <ActiveLabHero active={view.Active} />
      ) : (
        <EmptyState
          title="No lab in progress"
          description="Start a lab from a phase page or the challenge browser to see it here."
        />
      )}

      <VStack gap={2}>
        <Heading level={3}>Overall activity</Heading>
        <Grid columns={{ minWidth: 180, max: 6 }} gap={3}>
          <StatCard label="Streak" value={`${view.Streak} day${view.Streak === 1 ? "" : "s"}`} />
          <StatCard label="Labs completed" value={view.LabsCompleted} />
          <StatCard label="Variants seen" value={view.VariantsSeen} />
          {view.DockerVersion && <StatCard label="Docker" value={view.DockerVersion} />}
          {view.Uptime && <StatCard label="Uptime" value={view.Uptime} />}
          {view.AvgResolutionTime && <StatCard label="Avg resolution time" value={view.AvgResolutionTime} />}
        </Grid>
      </VStack>

      {emptyPhases.map((phase) => (
        <Banner
          key={phase.ID}
          status="info"
          title={`${phase.Title} — coming soon`}
          description={phase.Blurb}
          endContent={<Link href={`/phases/${phase.ID}`}>View phase</Link>}
        />
      ))}

      <VStack gap={2}>
        <TabList value={tab} onChange={(v) => setTab(v as "in-progress" | "completed")}>
          <Tab value="in-progress" label="In Progress" />
          <Tab value="completed" label="Completed" />
        </TabList>
        {visibleStages.length === 0 ? (
          <EmptyState
            title={tab === "completed" ? "No completed stages yet" : "No stages in progress"}
            description="Pick a phase to start a stage's labs."
            isCompact
          />
        ) : (
          <Grid columns={{ minWidth: 220, max: 4 }} gap={3}>
            {visibleStages.map((stage) => (
              <StageCard key={stage.ID} stage={stage} />
            ))}
          </Grid>
        )}
      </VStack>

      {view.MonthlyEngagement && (
        <MonthlyEngagement engagement={view.MonthlyEngagement} mostEngagedDay={view.MostEngagedDay} />
      )}

      {view.Recent && view.Recent.length > 0 && (
        <VStack gap={2}>
          <Heading level={3}>Recent postmortems</Heading>
          {view.Recent.map((pm) => (
            <HStack key={pm.TicketID} gap={2} justify="between">
              <Link href={pm.Href}>{pm.Lab}</Link>
              <Text type="supporting">{pm.Date}</Text>
            </HStack>
          ))}
        </VStack>
      )}
    </VStack>
  );
}

export default function Dashboard() {
  return (
    <Shell>
      <DashboardContent />
    </Shell>
  );
}
