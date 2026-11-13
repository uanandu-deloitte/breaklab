import { VStack } from "@astryxdesign/core/VStack";
import { Skeleton } from "@astryxdesign/core/Skeleton";
import { VisuallyHidden } from "@astryxdesign/core/VisuallyHidden";

// A shimmer placeholder shaped like the page it's standing in for (a title
// bar plus a stack of card-sized blocks), rather than a lone "Loading…"
// string floating on an empty page. The real label stays in the accessible
// tree via role="status"; the shapes themselves are aria-hidden decoration.
export function LoadingState({ label = "Loading…", rows = 3 }: { label?: string; rows?: number }) {
  return (
    <div role="status" aria-live="polite">
      <VisuallyHidden>{label}</VisuallyHidden>
      <VStack gap={6} aria-hidden="true">
        <VStack gap={2}>
          <Skeleton width={220} height={28} radius={2} />
          <Skeleton width={360} height={16} radius={2} />
        </VStack>
        <VStack gap={3}>
          {Array.from({ length: rows }, (_, i) => (
            <Skeleton key={i} height={88} radius={3} index={i} />
          ))}
        </VStack>
      </VStack>
    </div>
  );
}
