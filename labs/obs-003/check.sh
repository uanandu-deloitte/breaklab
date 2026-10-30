#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

RULES=/etc/alerts/rules.yaml
METRICS=/var/lib/metrics/current.json

REQUIRED_METRIC=error_rate
REQUIRED_OP=">"
REQUIRED_THRESHOLD=5

[ -s "$RULES" ] || fail "/etc/alerts/rules.yaml doesn't exist yet (or is empty) — author the rule"

metric=$(sed -n 's/^metric: *//p' "$RULES")
op=$(sed -n 's/^op: *//p' "$RULES" | tr -d '"')
threshold=$(sed -n 's/^threshold: *//p' "$RULES")

[ "$metric" = "$REQUIRED_METRIC" ] || fail "rules.yaml metric is '$metric', expected '$REQUIRED_METRIC'"
[ "$op" = "$REQUIRED_OP" ] || fail "rules.yaml op is '$op', expected '$REQUIRED_OP'"
[ "$threshold" = "$REQUIRED_THRESHOLD" ] || fail "rules.yaml threshold is '$threshold', expected '$REQUIRED_THRESHOLD'"

value=$(jq -r --arg m "$metric" '.[$m]' "$METRICS" 2>/dev/null)
[ -n "$value" ] && [ "$value" != "null" ] || fail "metric '$metric' isn't present in $METRICS"

awk -v v="$value" -v t="$threshold" 'BEGIN { exit !(v > t) }' \
  || fail "current ${metric}=${value} does not actually breach > ${threshold} — can't confirm a genuine fire"

: > /var/log/alerts.log
alert-evaluator

grep -q "${metric}=${value} breached ${op} ${threshold}" /var/log/alerts.log \
  || fail "alert-evaluator did not fire a fresh alert for ${metric}=${value} breaching ${op} ${threshold}"

echo "PASS: rule fires for the live ${metric}=${value} breach"
