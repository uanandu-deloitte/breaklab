#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
[ -f /workspace/config.old ] || fail "/workspace/config.old is missing"
[ -f /workspace/config.new ] || fail "/workspace/config.new is missing"

# shellcheck disable=SC1090
source "$ANSWERS"

[ "${DATABASE_HOST_OLD:-}" = "db-old.internal" ] || fail "DATABASE_HOST_OLD should be \"db-old.internal\", got \"${DATABASE_HOST_OLD:-}\""
[ "${DATABASE_HOST_NEW:-}" = "db-new.internal" ] || fail "DATABASE_HOST_NEW should be \"db-new.internal\", got \"${DATABASE_HOST_NEW:-}\""

[ "${LOG_LEVEL_OLD:-}" = "debug" ] || fail "LOG_LEVEL_OLD should be \"debug\", got \"${LOG_LEVEL_OLD:-}\""
[ "${LOG_LEVEL_NEW:-}" = "info" ] || fail "LOG_LEVEL_NEW should be \"info\", got \"${LOG_LEVEL_NEW:-}\""

[ "${TIMEOUT_OLD:-}" = "30" ] || fail "TIMEOUT_OLD should be \"30\", got \"${TIMEOUT_OLD:-}\""
[ "${TIMEOUT_NEW:-}" = "60" ] || fail "TIMEOUT_NEW should be \"60\", got \"${TIMEOUT_NEW:-}\""

echo "PASS: all 3 changed keys and their old/new values are correctly recorded"
