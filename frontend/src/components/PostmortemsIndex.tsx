import { useEffect, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { VStack } from "@astryxdesign/core/VStack";
import { HStack } from "@astryxdesign/core/HStack";
import { Badge } from "@astryxdesign/core/Badge";
import { Selector } from "@astryxdesign/core/Selector";
import { List, ListItem } from "@astryxdesign/core/List";
import { EmptyState } from "@astryxdesign/core/EmptyState";
import { Button } from "@astryxdesign/core/Button";
import Shell from "./Shell";

// Mirrors internal/web's PostmortemsIndexView (GET /api/postmortems).
interface PostmortemRowView {
  Lab: string;
  StageTitle: string;
  Difficulty: number;
  Date: string;
  Status: string; // "draft" or "published"
  TicketID: string;
  Href: string;
}
interface PostmortemsIndexView {
  Rows: PostmortemRowView[] | null;
  Stages: { ID: string; Title: string }[];
}

function PostmortemsIndexContent() {
  const [view, setView] = useState<PostmortemsIndexView | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [stage, setStage] = useState("");
  const [difficulty, setDifficulty] = useState("");
  const [status, setStatus] = useState("");

  useEffect(() => {
    const params = new URLSearchParams();
    if (stage) params.set("stage", stage);
    if (difficulty) params.set("difficulty", difficulty);
    if (status) params.set("status", status);
    fetch(`/api/postmortems?${params}`, { credentials: "include" })
      .then((res) => {
        if (!res.ok) throw new Error(`server returned ${res.status}`);
        return res.json();
      })
      .then(setView)
      .catch((err) => setError(String(err)));
  }, [stage, difficulty, status]);

  if (error) return <Text color="accent">Couldn't load postmortems: {error}</Text>;
  if (!view) return <Text type="supporting">Loading…</Text>;

  return (
    <VStack gap={4}>
      <Heading level={1}>Postmortems</Heading>

      <HStack gap={3} wrap="wrap">
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
          label="Status"
          isLabelHidden
          placeholder="All statuses"
          hasClear
          value={status}
          onChange={(v) => setStatus(v ?? "")}
          options={[
            { value: "draft", label: "Draft" },
            { value: "published", label: "Published" },
          ]}
        />
      </HStack>

      {view.Rows && view.Rows.length > 0 ? (
        <List hasDividers edge="outset">
          {view.Rows.map((row) => (
            <ListItem
              key={row.Href}
              label={row.Lab}
              description={`${row.StageTitle} · difficulty ${row.Difficulty} · ${row.Date}`}
              href={row.Href}
              endContent={<Badge label={row.Status} variant={row.Status === "published" ? "success" : "neutral"} />}
            />
          ))}
        </List>
      ) : stage || difficulty || status ? (
        <EmptyState
          title="No postmortems match these filters"
          description="Try a different stage, difficulty, or status."
          actions={
            <Button
              label="Clear filters"
              variant="secondary"
              clickAction={() => {
                setStage("");
                setDifficulty("");
                setStatus("");
              }}
            />
          }
        />
      ) : (
        <EmptyState
          title="No postmortems yet"
          description="Complete a lab and write its postmortem — Summary, Detection, Root cause, Fix, Prevention — to see it here."
        />
      )}
    </VStack>
  );
}

export default function PostmortemsIndex() {
  return (
    <Shell>
      <PostmortemsIndexContent />
    </Shell>
  );
}
