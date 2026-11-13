import { ClickableCard } from "@astryxdesign/core/ClickableCard";
import { VStack } from "@astryxdesign/core/VStack";
import { HStack } from "@astryxdesign/core/HStack";
import { Heading } from "@astryxdesign/core/Heading";
import { Text } from "@astryxdesign/core/Text";
import { Badge } from "@astryxdesign/core/Badge";
import { Link } from "@astryxdesign/core/Link";
import { Icon } from "@astryxdesign/core/Icon";

export const labStatusVariant: Record<string, "neutral" | "info" | "success" | "warning"> = {
  available: "info",
  "in-progress": "warning",
  complete: "success",
  locked: "neutral",
};

function DifficultyIcon() {
  return (
    <svg
      width="14"
      height="14"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      aria-hidden="true"
    >
      <path d="M4 20V14M12 20V9M20 20V4" />
    </svg>
  );
}

export interface LabCardData {
  ID: string;
  Title: string;
  Status: string;
  Difficulty: number;
  EstimatedMinutes: number;
  ReviewDue: boolean;
  IsUpNext?: boolean;
}

export function labHref(labID: string, backHref?: string, backLabel?: string): string {
  if (!backHref) return `/labs/${labID}`;
  const params = new URLSearchParams({ back: backHref });
  if (backLabel) params.set("backLabel", backLabel);
  return `/labs/${labID}?${params}`;
}

export function LabCard({
  lab,
  stageLabel,
  backHref,
  backLabel,
}: {
  lab: LabCardData;
  stageLabel?: string;
  backHref?: string;
  backLabel?: string;
}) {
  const href = labHref(lab.ID, backHref, backLabel);
  return (
    <ClickableCard label={lab.Title} href={href}>
      <VStack gap={3}>
        <HStack gap={2} vAlign="center" wrap="wrap">
          <Badge label={lab.Status} variant={labStatusVariant[lab.Status] ?? "neutral"} />
          {lab.IsUpNext && <Badge label="Up next" variant="info" />}
          {lab.ReviewDue && <Badge label="Review due" variant="warning" />}
        </HStack>
        <Heading level={3}>
          <Link href={href}>{lab.Title}</Link>
        </Heading>
        <HStack gap={4} wrap="wrap">
          <HStack gap={1} vAlign="center">
            <DifficultyIcon />
            <Text type="supporting">Difficulty {lab.Difficulty}</Text>
          </HStack>
          <HStack gap={1} vAlign="center">
            <Icon icon="clock" size="sm" />
            <Text type="supporting">{lab.EstimatedMinutes} min</Text>
          </HStack>
          {stageLabel && <Text type="supporting">{stageLabel}</Text>}
        </HStack>
      </VStack>
    </ClickableCard>
  );
}
