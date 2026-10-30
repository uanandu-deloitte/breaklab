#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

check_one() {
  local var="$1" want="$2" got
  got="${!var:-}"
  [ -n "$got" ] || fail "$var is empty"
  [ "$got" = "$want" ] || fail "$var should be '$want', got '$got'"
}

check_one ENDPOINT_1 dns-failure
check_one ENDPOINT_2 refused
check_one ENDPOINT_3 timeout
check_one ENDPOINT_4 404
check_one ENDPOINT_5 500

echo "PASS: all 5 endpoints classified correctly"
