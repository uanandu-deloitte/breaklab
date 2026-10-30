import { useEffect, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { Card } from "@astryxdesign/core/Card";
import { CodeBlock } from "@astryxdesign/core/CodeBlock";
import { Grid } from "@astryxdesign/core/Grid";
import { VStack } from "@astryxdesign/core/VStack";
import { HStack } from "@astryxdesign/core/HStack";
import { Badge } from "@astryxdesign/core/Badge";
import { Link } from "@astryxdesign/core/Link";
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

interface DashboardView {
  Streak: number;
  LabsCompleted: number;
  VariantsSeen: number;
  Active: ActiveLabView | null;
  Recent: RecentPostmortemView[] | null;
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

function DashboardContent() {
  const [view, setView] = useState<DashboardView | null>(null);
  const [error, setError] = useState<string | null>(null);

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
  if (!view) return <Text type="supporting">Loading…</Text>;

  return (
    <VStack gap={6}>
      <Heading level={1}>Dashboard</Heading>

      <Grid columns={{ minWidth: 180, max: 6 }} gap={3}>
        <StatCard label="Streak" value={`${view.Streak} day${view.Streak === 1 ? "" : "s"}`} />
        <StatCard label="Labs completed" value={view.LabsCompleted} />
        <StatCard label="Variants seen" value={view.VariantsSeen} />
        {view.DockerVersion && <StatCard label="Docker" value={view.DockerVersion} />}
        {view.Uptime && <StatCard label="Uptime" value={view.Uptime} />}
        {view.AvgResolutionTime && <StatCard label="Avg resolution time" value={view.AvgResolutionTime} />}
      </Grid>

      {view.Active && (
        <Card variant="blue">
          <VStack gap={2}>
            <HStack gap={2} vAlign="center">
              <Heading level={3}>
                <Link href={view.Active.Href}>
                  {view.Active.Lab}: {view.Active.Title}
                </Link>
              </Heading>
              <Badge label={view.Active.TicketID} variant="info" />
            </HStack>
            <Text type="supporting">
              {view.Active.ElapsedMinutes}m elapsed · {view.Active.HintsUsed}/{view.Active.HintsTotal} hints used
            </Text>
            {view.Active.ShellCmd && (
              <VStack gap={2}>
                <Heading level={4}>Get a shell</Heading>
                <CodeBlock code={view.Active.ShellCmd} language="bash" />
              </VStack>
            )}
            <CodeBlock code={view.Active.ReminderCmd} language="bash" />
          </VStack>
        </Card>
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
