// BreakLab's mark: a terracotta badge (the one "click me"/brand hue —
// see DESIGN.md's Signal Terracotta) with a white fault-line bolt,
// echoing the break/fix premise of the labs themselves.
export default function BrandMark() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true">
      <rect width="24" height="24" rx="6" fill="#c14e1f" />
      <path d="M12.5 2 5 13h5.5L9.5 22 19 10h-5.5z" fill="#fff" />
    </svg>
  );
}
