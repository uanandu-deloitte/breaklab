import { useEffect, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { VStack } from "@astryxdesign/core/VStack";
import { HStack } from "@astryxdesign/core/HStack";
import { Badge } from "@astryxdesign/core/Badge";
import { Button } from "@astryxdesign/core/Button";
import { Link } from "@astryxdesign/core/Link";
import { TextInput } from "@astryxdesign/core/TextInput";
import { Selector } from "@astryxdesign/core/Selector";
import { List, ListItem } from "@astryxdesign/core/List";
import { EmptyState } from "@astryxdesign/core/EmptyState";
import Shell from "./Shell";

// Mirrors internal/web's ChallengeView (GET /api/challenges).
interface LabRowView {
  ID: string;
  Title: string;
  Summary: string;
  StageTitle: string;
  Status: string;
  Difficulty: number;
  EstimatedMinutes: number;
  ReviewDue: boolean;
}

interface ActiveIncidentView {
  Lab: string;
  TicketID: string;
  Href: string;
}

interface ChallengeView {
  Rows: LabRowView[] | null;
  Stages: { ID: string; Title: string }[];
  Tags: string[] | null;
  ActiveIncident: ActiveIncidentView | null;
}

const statusVariant: Record<string, "neutral" | "info" | "success" | "warning"> = {
  available: "info",
  "in-progress": "warning",
  complete: "success",
  locked: "neutral",
};

function ChallengesContent() {
  const [view, setView] = useState<ChallengeView | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [query, setQuery] = useState("");
  const [stage, setStage] = useState("");
  const [difficulty, setDifficulty] = useState("");
  const [sortBy, setSortBy] = useState("");

  function load() {
    const params = new URLSearchParams();
    if (query) params.set("q", query);
    if (stage) params.set("stage", stage);
    if (difficulty) params.set("difficulty", difficulty);
    if (sortBy) params.set("sort", sortBy);
    fetch(`/api/challenges?${params}`, { credentials: "include" })
      .then((res) => {
        if (!res.ok) throw new Error(`server returned ${res.status}`);
        return res.json();
      })
      .then(setView)
      .catch((err) => setError(String(err)));
  }

  useEffect(load, [query, stage, difficulty, sortBy]);

  async function launchRandom() {
    const params = new URLSearchParams();
    if (query) params.set("q", query);
    if (stage) params.set("stage", stage);
    if (difficulty) params.set("difficulty", difficulty);
    const res = await fetch(`/api/challenges/random?${params}`, { method: "POST", credentials: "include" });
    const payload = await res.json();
    if (payload.lab) window.location.href = `/labs/${payload.lab}`;
  }

  if (error) return <Text color="accent">Couldn't load challenges: {error}</Text>;
  if (!view) return <Text type="supporting">Loading…</Text>;

  return (
    <VStack gap={4}>
      <Heading level={1}>Challenges</Heading>

      {view.ActiveIncident && (
        <HStack gap={2} vAlign="center">
          <Badge label="active incident" variant="warning" />
          <Text type="code">{view.ActiveIncident.TicketID}</Text>
          <Link href={view.ActiveIncident.Href}>{view.ActiveIncident.Lab}</Link>
        </HStack>
      )}

      <HStack gap={3} wrap="wrap" vAlign="start">
        <TextInput label="Search" isLabelHidden placeholder="Search by title or tag…" value={query} onChange={setQuery} width={240} />
        <Selector
          label="Stage"
          isLabelHidden
          placeholder="All stages"
          hasClear
          value={stage}
          onChange={(v) => setStage(v ?? "")}
          options={view.Stages.map((s) => ({ value: s.ID, label: s.Title }))}
        />
        <Selector
          label="Difficulty"
          isLabelHidden
          placeholder="All difficulties"
          hasClear
          value={difficulty}
          onChange={(v) => setDifficulty(v ?? "")}
          options={["1", "2", "3", "4", "5"].map((d) => ({ value: d, label: `difficulty ${d}` }))}
        />
        <Selector
          label="Sort"
          isLabelHidden
          placeholder="Recent activity"
          value={sortBy}
          onChange={(v) => setSortBy(v ?? "")}
          options={[
            { value: "", label: "Recent activity" },
            { value: "difficulty", label: "Difficulty" },
            { value: "alpha", label: "Alphabetical" },
          ]}
        />
        <Button label="Launch random" variant="secondary" clickAction={launchRandom} />
      </HStack>

      {view.Rows && view.Rows.length > 0 ? (
        <List hasDividers edge="outset">
          {view.Rows.map((lab) => (
            <ListItem
              key={lab.ID}
              label={lab.Title}
              description={lab.Summary}
              href={`/labs/${lab.ID}`}
              endContent={
                <HStack gap={2}>
                  <Badge label={lab.Status} variant={statusVariant[lab.Status] ?? "neutral"} />
                  <Badge label={`difficulty ${lab.Difficulty}`} />
                  <Badge label={`${lab.EstimatedMinutes}m`} />
                  {lab.ReviewDue && <Badge label="Review due" variant="warning" />}
                </HStack>
              }
            />
          ))}
        </List>
      ) : (
        <EmptyState
          title="No challenges match these filters"
          description="Try a different search, stage, or difficulty."
          actions={
            <Button
              label="Clear filters"
              variant="secondary"
              clickAction={() => {
                setQuery("");
                setStage("");
                setDifficulty("");
              }}
            />
          }
        />
      )}
    </VStack>
  );
}

export default function Challenges() {
  return (
    <Shell>
      <ChallengesContent />
    </Shell>
  );
}
