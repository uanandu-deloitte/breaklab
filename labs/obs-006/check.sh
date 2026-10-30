#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

CONFIG=/etc/dashboards/overview.yaml
EXPORTED=/var/lib/metrics/exported.json

[ -f "$CONFIG" ] || fail "$CONFIG is missing"
[ -f "$EXPORTED" ] || fail "$EXPORTED is missing"

panel_metric=$(sed -n 's/^metric: *//p' "$CONFIG")
[ -n "$panel_metric" ] || fail "could not read the metric field from $CONFIG"

[ -n "${PANEL_METRIC:-}" ] || fail "PANEL_METRIC is empty"
[ "$PANEL_METRIC" = "$panel_metric" ] \
  || fail "PANEL_METRIC is '$PANEL_METRIC', but $CONFIG actually configures metric '$panel_metric'"

[ -n "${EXPORTED_METRIC:-}" ] || fail "EXPORTED_METRIC is empty"

jq -e --arg m "$EXPORTED_METRIC" 'has($m)' "$EXPORTED" >/dev/null 2>&1 \
  || fail "EXPORTED_METRIC '$EXPORTED_METRIC' is not actually a key in $EXPORTED"

[ "$EXPORTED_METRIC" != "$panel_metric" ] \
  || fail "EXPORTED_METRIC matches the panel's own metric name ('$panel_metric') — that's not the mismatch, find the exporter key that's actually different"

jq -e --arg m "$panel_metric" 'has($m)' "$EXPORTED" >/dev/null 2>&1 \
  && fail "panel metric '$panel_metric' is actually present in $EXPORTED — there's no rename gap here, re-check your answer"

rm -f /var/lib/dashboards/overview.rendered
render-dashboard >/dev/null 2>&1
rendered=$(cat /var/lib/dashboards/overview.rendered 2>/dev/null || echo "")

[ -n "${RENDER_RESULT:-}" ] || fail "RENDER_RESULT is empty"
[ "$RENDER_RESULT" = "$rendered" ] \
  || fail "RENDER_RESULT is '$RENDER_RESULT', but render-dashboard actually wrote '$rendered'"
echo "$rendered" | grep -qi 'no-data' \
  || fail "render-dashboard actually produced a real value ('$rendered'), not no-data — re-check your answer against the live render"

echo "PASS: correctly identified panel metric '$panel_metric' vs exported metric '$EXPORTED_METRIC' behind the silent no-data gap"
