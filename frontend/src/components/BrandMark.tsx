// BreakLab's mark: a badge in the theme's accent color (Deloitte green)
// with a fault-line bolt in the matching on-accent tone, echoing the
// break/fix premise of the labs themselves. Both fills reference theme
// tokens directly, so the mark follows the active theme/mode instead of
// a hardcoded color.
export default function BrandMark() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true">
      <rect width="24" height="24" rx="6" fill="var(--color-accent)" />
      <path d="M12.5 2 5 13h5.5L9.5 22 19 10h-5.5z" fill="var(--color-on-accent)" />
    </svg>
  );
}
