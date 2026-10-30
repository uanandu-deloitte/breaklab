#!/usr/bin/env bash
set -uo pipefail

GRACEFUL=/workspace/worker-graceful.marker
FORCED=/workspace/worker-forced.marker
ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

[ -f "$GRACEFUL" ] || fail "$GRACEFUL is missing — worker1 was not terminated gracefully with SIGTERM yet"
if pgrep -f 'worker1' >/dev/null; then
  fail "worker1 is still running — send it SIGTERM (plain kill) and give it a moment to exit"
fi

if [ -f "$FORCED" ]; then
  fail "$FORCED exists — a worker killed with SIGKILL must never get the chance to write its marker"
fi
if pgrep -f 'worker-forced.marker' >/dev/null; then
  fail "the second worker instance (started with /workspace/worker-forced.marker) is still running — kill -9 it"
fi

[ "${GRACEFUL_MARKER_EXISTS:-}" = "yes" ] || fail "GRACEFUL_MARKER_EXISTS must be set to \"yes\" in answers.env"
[ "${FORCED_MARKER_EXISTS:-}" = "no" ] || fail "FORCED_MARKER_EXISTS must be set to \"no\" in answers.env"

echo "PASS: SIGTERM let worker1 write its marker; SIGKILL gave the second instance no such chance"
