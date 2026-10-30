import { useEffect, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { VStack } from "@astryxdesign/core/VStack";
import { CodeBlock } from "@astryxdesign/core/CodeBlock";
import Shell from "./Shell";

// Mirrors internal/web's PostmortemPageView (GET /api/postmortems/{slug})
// and internal/postmortem's Block/DiffLine shapes.
interface DiffLine {
  Sign: string;
  Text: string;
}
interface Block {
  Kind: number; // 0 = paragraph, 1 = code, 2 = diff (postmortem.BlockKind)
  Text: string;
  Diff: DiffLine[] | null;
}
interface SectionView {
  Name: string;
  Blocks: Block[] | null;
}
interface PostmortemPageView {
  Lab: string;
  Variant: string;
  Date: string;
  TicketID: string;
  Sections: SectionView[] | null;
}

function slugFromPath(): string {
  if (typeof window === "undefined") return "";
  const parts = window.location.pathname.split("/").filter(Boolean);
  return parts[1] ?? ""; // "/postmortems/<slug>"
}

function BlockView({ block }: { block: Block }) {
  if (block.Kind === 1) {
    return <CodeBlock code={block.Text} language="plaintext" />;
  }
  if (block.Kind === 2) {
    const diffText = (block.Diff ?? []).map((l) => l.Sign + l.Text).join("\n");
    return <CodeBlock code={diffText} language="plaintext" isDiff />;
  }
  return <Text as="p">{block.Text}</Text>;
}

function PostmortemDetailContent() {
  const [slug, setSlug] = useState("");
  const [view, setView] = useState<PostmortemPageView | null>(null);
  const [notFound, setNotFound] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => setSlug(slugFromPath()), []);

  useEffect(() => {
    if (!slug) return;
    fetch(`/api/postmortems/${slug}`, { credentials: "include" })
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
  }, [slug]);

  if (notFound) return <Text color="accent">Postmortem not found.</Text>;
  if (error) return <Text color="accent">Couldn't load this postmortem: {error}</Text>;
  if (!view) return <Text type="supporting">Loading…</Text>;

  return (
    <VStack gap={4}>
      <Heading level={1}>{view.Lab}</Heading>
      <Text type="supporting">
        {view.Variant} · {view.Date}
        {view.TicketID && ` · ${view.TicketID}`}
      </Text>
      {(view.Sections ?? []).map((section) => (
        <VStack gap={2} key={section.Name}>
          <Heading level={2}>{section.Name}</Heading>
          {(section.Blocks ?? []).map((block, i) => <BlockView key={i} block={block} />)}
        </VStack>
      ))}
    </VStack>
  );
}

export default function PostmortemDetail() {
  return (
    <Shell>
      <PostmortemDetailContent />
    </Shell>
  );
}
