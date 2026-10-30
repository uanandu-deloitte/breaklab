import { useEffect, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { VStack } from "@astryxdesign/core/VStack";
import { HStack } from "@astryxdesign/core/HStack";
import { Badge } from "@astryxdesign/core/Badge";
import { List, ListItem } from "@astryxdesign/core/List";
import Shell from "./Shell";

// Mirrors internal/web's ChaosView (GET /api/chaos).
interface LabRowView {
  ID: string;
  Title: string;
  Summary: string;
  Status: string;
  Difficulty: number;
}
interface ChaosConceptView {
  ID: string;
  Title: string;
  Blurb: string;
  Labs: LabRowView[] | null;
}
interface ChaosView {
  Concepts: ChaosConceptView[] | null;
}

const statusVariant: Record<string, "neutral" | "info" | "success" | "warning"> = {
  available: "info",
  "in-progress": "warning",
  complete: "success",
  locked: "neutral",
};

function ChaosContent() {
  const [view, setView] = useState<ChaosView | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch("/api/chaos", { credentials: "include" })
      .then((res) => {
        if (!res.ok) throw new Error(`server returned ${res.status}`);
        return res.json();
      })
      .then(setView)
      .catch((err) => setError(String(err)));
  }, []);

  if (error) return <Text color="accent">Couldn't load chaos labs: {error}</Text>;
  if (!view) return <Text type="supporting">Loading…</Text>;

  return (
    <VStack gap={6}>
      <Heading level={1}>Chaos engineering</Heading>
      {(view.Concepts ?? []).map((concept) => (
        <VStack gap={2} key={concept.ID}>
          <Heading level={2}>{concept.Title}</Heading>
          <Text type="supporting">{concept.Blurb}</Text>
          <List hasDividers edge="outset">
            {(concept.Labs ?? []).map((lab) => (
              <ListItem
                key={lab.ID}
                label={lab.Title}
                description={lab.Summary}
                href={`/labs/${lab.ID}`}
                endContent={
                  <HStack gap={2}>
                    <Badge label={lab.Status} variant={statusVariant[lab.Status] ?? "neutral"} />
                    <Badge label={`difficulty ${lab.Difficulty}`} />
                  </HStack>
                }
              />
            ))}
          </List>
        </VStack>
      ))}
    </VStack>
  );
}

export default function Chaos() {
  return (
    <Shell>
      <ChaosContent />
    </Shell>
  );
}
