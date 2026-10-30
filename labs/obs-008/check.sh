#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

cd /infra || fail "/infra is missing"
[ -f rules.yml ] || fail "/infra/rules.yml is missing"
[ -f rules_test.yml ] || fail "/infra/rules_test.yml is missing"

out=$(promtool test rules rules_test.yml 2>&1)
exit_code=$?

[ "$exit_code" -eq 0 ] || fail "promtool test rules failed:
$out"

echo "PASS: HighErrorRate fires exactly as rules_test.yml expects"
