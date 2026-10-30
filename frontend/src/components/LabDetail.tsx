import { useEffect, useState } from "react";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { VStack } from "@astryxdesign/core/VStack";
import { HStack } from "@astryxdesign/core/HStack";
import { Badge } from "@astryxdesign/core/Badge";
import { Button } from "@astryxdesign/core/Button";
import { Card } from "@astryxdesign/core/Card";
import { CodeBlock } from "@astryxdesign/core/CodeBlock";
import { TextArea } from "@astryxdesign/core/TextArea";
import { List, ListItem } from "@astryxdesign/core/List";
import { Link } from "@astryxdesign/core/Link";
import { Icon } from "@astryxdesign/core/Icon";
import { Markdown } from "@astryxdesign/core/Markdown";
import { Layout, LayoutHeader, LayoutContent, LayoutPanel } from "@astryxdesign/core/Layout";
import Shell from "./Shell";

// Mirrors internal/web's LabDetailView (GET /api/labs/{id}).
interface PortNodeView {
  HostPort: string;
  ContainerPort: string;
  OpenHref: string;
}
interface HintSlotView {
  Revealed: boolean;
  Cost: number;
}
interface GuidedStepView {
  Do: string;
  Why: string;
}
interface PracticeQuestionView {
  Prompt: string;
  Answer: string;
}
interface LabDetailView {
  ID: string;
  Title: string;
  StageTitle: string;
  Status: string;
  Difficulty: number;
  EstimatedMinutes: number;
  Tags: string[] | null;
  Ticket: string;
  TicketID: string;
  Ports: PortNodeView[] | null;
  HintsUsed: number;
  HintSlots: HintSlotView[] | null;
  DiagnosticLog: string;
  Command: string;
  Active: boolean;
  ShellCmd: string;
  CanLaunch: boolean;
  BlockingLab: string;
  Summary: string;
  Learn: string;
  Lesson: string; // pre-rendered HTML (internal/markdown), distinct from the plain-text Learn field — a lab sets one or the other
  Guided: GuidedStepView[] | null;
  Practice: { Task: string; Questions: PracticeQuestionView[] | null };
  Hypothesis: string;
  LockedOnTitles: string[] | null;
  ChallengeMode: boolean;
}

function labIDFromPath(): string {
  if (typeof window === "undefined") return ""; // SSG prerender has no window
  const parts = window.location.pathname.split("/").filter(Boolean);
  return parts[1] ?? ""; // "/labs/<id>"
}

// ChallengeMode never gates the Hint action server-side — the API
// would still reveal a hint if asked — it only hides this page's own
// Hints section and its reveal button, same as the old html/template
// page did.
function challengeModeFromQuery(): boolean {
  if (typeof window === "undefined") return false;
  return new URLSearchParams(window.location.search).get("mode") === "challenge";
}

// LabDetailContent renders one lab's detail + actions. fixedLabID, when
// set, skips the path-derived id entirely — the capstone page's own
// lab isn't at /labs/{id}, so it passes its id directly (Capstone.tsx).
function LabDetailContent({ fixedLabID }: { fixedLabID?: string }) {
  const [labID, setLabID] = useState(fixedLabID ?? "");
  const [view, setView] = useState<LabDetailView | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [notFound, setNotFound] = useState(false);
  const [checkMsg, setCheckMsg] = useState<string | null>(null);
  const [checkPassed, setCheckPassed] = useState(false);
  const [revealedHint, setRevealedHint] = useState<string | null>(null);
  const [hypothesis, setHypothesis] = useState("");
  const [actionError, setActionError] = useState<string | null>(null);
  const [liveLogHTML, setLiveLogHTML] = useState("");
  const [pendingAction, setPendingAction] = useState<"start" | "check" | "hint" | "stop" | null>(null);

  useEffect(() => {
    if (!fixedLabID) setLabID(labIDFromPath());
  }, [fixedLabID]);

  // Diagnostic log stream (web-ui spec): a live-updating view of the
  // active session's own container log, scoped to this one lab's
  // detail page. Each SSE line arrives pre-escaped and wrapped in a
  // <span class="log-line"> by the server specifically so it's safe to
  // insert as HTML (never as the attacker-influenceable raw text it
  // actually is) — appending it via dangerouslySetInnerHTML preserves
  // that same guarantee, it does not bypass it.
  // Mirrors lab.go's isMystery: only an incident-style lab (an
  // untyped/legacy lab, or one explicitly typed "incident") has a
  // process actively writing to stdout for this to stream. A guided or
  // practice lab's container just sits idle (sleep infinity) around a
  // static fixture file, so there's nothing to follow — opening a
  // connection for one would just sit empty forever.
  const isIncidentLab = !view || view.Type === "" || view.Type === "incident";

  useEffect(() => {
    if (!labID || !view?.Active || !isIncidentLab) return;
    setLiveLogHTML("");
    const source = new EventSource(`/labs/${labID}/logs/stream`);
    source.onmessage = (event) => {
      setLiveLogHTML((prev) => prev + event.data + "\n");
    };
    return () => source.close(); // stops following logs when leaving the page
  }, [labID, view?.Active, isIncidentLab]);

  function load() {
    if (!labID) return;
    const mode = challengeModeFromQuery() ? "?mode=challenge" : "";
    fetch(`/api/labs/${labID}${mode}`, { credentials: "include" })
      .then((res) => {
        if (res.status === 404) {
          setNotFound(true);
          return null;
        }
        if (!res.ok) throw new Error(`server returned ${res.status}`);
        return res.json();
      })
      .then((data) => {
        if (data) {
          setView(data);
          setHypothesis(data.Hypothesis ?? "");
        }
      })
      .catch((err) => setError(String(err)));
  }

  useEffect(load, [labID]);

  async function postAction(path: string, body?: unknown) {
    setActionError(null);
    const res = await fetch(`/api/labs/${labID}/${path}`, {
      method: "POST",
      credentials: "include",
      headers: body ? { "Content-Type": "application/json" } : undefined,
      body: body ? JSON.stringify(body) : undefined,
    });
    if (!res.ok) {
      const payload = await res.json().catch(() => ({ error: `server returned ${res.status}` }));
      setActionError(payload.error ?? "action failed");
      return null;
    }
    return res.status === 204 ? true : res.json();
  }

  async function handleStart(confirm: boolean) {
    setPendingAction("start");
    try {
      if (await postAction("start", { confirm })) load();
    } finally {
      setPendingAction(null);
    }
  }
  async function handleCheck() {
    setPendingAction("check");
    try {
      const result = await postAction("check");
      if (result) {
        setCheckMsg(result.Message);
        setCheckPassed(result.Passed);
        load();
      }
    } finally {
      setPendingAction(null);
    }
  }
  async function handleHint() {
    setPendingAction("hint");
    try {
      const result = await postAction("hint");
      if (result) {
        setRevealedHint(result.hint);
        load();
      }
    } finally {
      setPendingAction(null);
    }
  }
  async function handleStop() {
    setPendingAction("stop");
    try {
      if (await postAction("stop")) load();
    } finally {
      setPendingAction(null);
    }
  }
  async function handleSaveHypothesis() {
    await postAction("hypothesis", { hypothesis });
  }

  if (notFound) return <Text color="accent">Lab not found.</Text>;
  if (error) return <Text color="accent">Couldn't load this lab: {error}</Text>;
  if (!view) return <Text type="supporting">Loading…</Text>;

  return (
    <Layout
      padding={6}
      header={
        <LayoutHeader hasDivider>
          <VStack gap={3}>
            <Link href="/challenges" type="supporting" color="secondary">
              <HStack as="span" gap={1} vAlign="center">
                <Icon icon="chevronLeft" size="sm" />
                Back to challenges
              </HStack>
            </Link>
            <HStack gap={4} vAlign="start" hAlign="between" wrap="wrap">
              <VStack gap={1}>
                <HStack gap={2} vAlign="center">
                  <Heading level={1}>{view.Title}</Heading>
                  <Badge label={view.Status} />
                </HStack>
                <Text type="supporting">
                  {view.StageTitle} · difficulty {view.Difficulty} · ~{view.EstimatedMinutes}m
                  {view.TicketID && ` · ${view.TicketID}`}
                </Text>
              </VStack>
              <HStack gap={2} wrap="wrap">
                {!view.Active && view.CanLaunch && (
                  <Button label="Start" variant="primary" isLoading={pendingAction === "start"} clickAction={() => handleStart(false)} />
                )}
                {!view.Active && !view.CanLaunch && view.BlockingLab && (
                  <Button
                    label={`Stop ${view.BlockingLab} and start this instead`}
                    variant="secondary"
                    isLoading={pendingAction === "start"}
                    clickAction={() => handleStart(true)}
                  />
                )}
                {view.Active && <Button label="Run check" variant="primary" isLoading={pendingAction === "check"} clickAction={handleCheck} />}
                {view.Active && !view.ChallengeMode && (
                  <Button label="Reveal next hint" variant="secondary" isLoading={pendingAction === "hint"} clickAction={handleHint} />
                )}
                {view.Active && <Button label="Stop" variant="secondary" isLoading={pendingAction === "stop"} clickAction={handleStop} />}
              </HStack>
            </HStack>
            {pendingAction === "start" && (
              <Card variant="blue">
                <Text>Starting lab… building the container image. The first run can take a few minutes while the base image and packages download — this is normal, not stuck.</Text>
              </Card>
            )}
            {actionError && (
              <Card variant="red">
                <Text>{actionError}</Text>
              </Card>
            )}
          </VStack>
        </LayoutHeader>
      }
      end={
        <LayoutPanel hasDivider width={320} role="complementary" label="Lab details">
          <VStack gap={6}>
            <VStack gap={2}>
              <Heading level={3}>Ticket</Heading>
              <Markdown>{view.Ticket}</Markdown>
            </VStack>

            {!view.CanLaunch && view.LockedOnTitles && view.LockedOnTitles.length > 0 && (
              <Card variant="orange">
                <Text>Finish first: {view.LockedOnTitles.join(", ")}</Text>
              </Card>
            )}

            {view.Active && !view.ChallengeMode && view.HintSlots && view.HintSlots.length > 0 && (
              <VStack gap={2}>
                <Heading level={3}>
                  Hints ({view.HintsUsed}/{view.HintSlots.length} used)
                </Heading>
                <List hasDividers edge="outset">
                  {view.HintSlots.map((slot, i) => (
                    <ListItem key={i} label={slot.Revealed ? `Hint ${i + 1} (revealed)` : `Hint ${i + 1} — ${slot.Cost} pts`} />
                  ))}
                </List>
              </VStack>
            )}

            <VStack gap={2}>
              <TextArea
                label="Your hypothesis"
                description="What do you think is causing this, before you run the check?"
                value={hypothesis}
                onChange={setHypothesis}
                rows={3}
              />
              <Button label="Save hypothesis" variant="secondary" size="sm" clickAction={handleSaveHypothesis} />
            </VStack>
          </VStack>
        </LayoutPanel>
      }
      content={
        <LayoutContent>
          <VStack gap={4}>
            {view.Summary && <Text>{view.Summary}</Text>}

            {view.Active && <CodeBlock code={view.Command} language="bash" />}

            {checkMsg && (
              <Card variant={checkPassed ? "green" : "red"}>
                <Text>{checkMsg}</Text>
              </Card>
            )}
            {revealedHint && (
              <Card variant="yellow">
                <Text>{revealedHint}</Text>
              </Card>
            )}

            {view.Active && view.DiagnosticLog && (
              <Card variant="gray">
                <VStack gap={2}>
                  <Heading level={4}>Container log tail</Heading>
                  <CodeBlock code={view.DiagnosticLog} language="plaintext" hasLanguageLabel={false} isWrapped />
                </VStack>
              </Card>
            )}

            {view.Active && isIncidentLab && (
              <Card variant="gray">
                <VStack gap={2}>
                  <Heading level={4}>Live log</Heading>
                  {/* No Astryx component accepts raw HTML children; this is one of
                      two deliberate exceptions on this page (see Lesson below) —
                      injecting the SSE stream's already-server-escaped <span>
                      lines (see the effect above for why that's still safe). */}
                  <div dangerouslySetInnerHTML={{ __html: liveLogHTML }} />
                </VStack>
              </Card>
            )}

            {view.Active && view.ShellCmd && (
              <VStack gap={2}>
                <Heading level={4}>Get a shell</Heading>
                <CodeBlock code={view.ShellCmd} language="bash" />
              </VStack>
            )}

            {view.Lesson && (
              <VStack gap={2}>
                <Heading level={3}>Learn</Heading>
                {/* Pre-rendered HTML from internal/markdown (lesson.md), not raw
                    markdown — Astryx's Markdown component takes markdown source,
                    not HTML, so this is the one other deliberate
                    dangerouslySetInnerHTML exception alongside the live log tail
                    above, for the same reason: already-safe, server-rendered markup. */}
                <div dangerouslySetInnerHTML={{ __html: view.Lesson }} />
              </VStack>
            )}
            {view.Learn && (
              <VStack gap={2}>
                <Heading level={3}>Learn</Heading>
                <Text as="p">{view.Learn}</Text>
              </VStack>
            )}

            {view.Guided && view.Guided.length > 0 && (
              <VStack gap={2}>
                <Heading level={3}>Guided steps</Heading>
                <List hasDividers edge="outset">
                  {view.Guided.map((step, i) => (
                    <ListItem key={i} label={step.Do} description={step.Why} />
                  ))}
                </List>
              </VStack>
            )}

            {view.Practice?.Questions && view.Practice.Questions.length > 0 && (
              <VStack gap={2}>
                <Heading level={3}>Practice</Heading>
                {view.Practice.Task && <Text>{view.Practice.Task}</Text>}
                <List hasDividers edge="outset">
                  {view.Practice.Questions.map((q, i) => (
                    <ListItem key={i} label={q.Prompt} description={q.Answer} />
                  ))}
                </List>
              </VStack>
            )}
          </VStack>
        </LayoutContent>
      }
    />
  );
}

// Exported unwrapped for Capstone.tsx, which embeds this inside its own
// Shell (AppShell can't nest inside another AppShell).
export { LabDetailContent };

export default function LabDetail({ fixedLabID }: { fixedLabID?: string }) {
  return (
    <Shell contentPadding={0}>
      <LabDetailContent fixedLabID={fixedLabID} />
    </Shell>
  );
}
