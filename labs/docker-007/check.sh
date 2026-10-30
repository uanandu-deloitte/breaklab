#!/usr/bin/env bash
# Runs on the host (run_in: host): docker logs can only be inspected from
# outside the container, so this check reads answers.env via `docker exec`
# and compares it against docker logs' real output and the real on-disk
# file, rather than trusting the learner's claims.
set -uo pipefail

container="${BREAKLAB_CONTAINER:?BREAKLAB_CONTAINER not set}"
fail() { echo "FAIL: $1"; exit 1; }

answers_raw=$(docker exec "$container" cat /workspace/answers.env 2>/dev/null) \
  || fail "/workspace/answers.env not found in the container"

tmpfile=$(mktemp)
printf '%s\n' "$answers_raw" > "$tmpfile"
# shellcheck disable=SC1090
source "$tmpfile"
rm -f "$tmpfile"

[ -n "${STDOUT_CAPTURED:-}" ] || fail "STDOUT_CAPTURED is empty"
[ -n "${FILE_LOG_CAPTURED:-}" ] || fail "FILE_LOG_CAPTURED is empty"
[ -n "${LOG_FILE_PATH:-}" ] || fail "LOG_FILE_PATH is empty"

echo "$STDOUT_CAPTURED" | grep -qi '^yes$' \
  || fail "STDOUT_CAPTURED should be yes — docker logs does capture PID 1's stdout"
echo "$FILE_LOG_CAPTURED" | grep -qi '^no$' \
  || fail "FILE_LOG_CAPTURED should be no — docker logs never shows a file written to disk inside the container"

# Ground truth: docker logs really does show the heartbeat (stdout) lines.
docker logs "$container" 2>&1 | grep -q '^heartbeat ' \
  || fail "docker logs shows no heartbeat lines from this container's stdout — is the app still running?"

# Ground truth: the on-disk log file genuinely exists with content docker
# logs never echoes, and the learner's path points at the real file.
docker exec "$container" test -f "$LOG_FILE_PATH" \
  || fail "LOG_FILE_PATH ($LOG_FILE_PATH) doesn't exist in the container"
docker exec "$container" grep -q '^file-only ' "$LOG_FILE_PATH" \
  || fail "LOG_FILE_PATH ($LOG_FILE_PATH) doesn't contain this app's file-only log lines"
docker logs "$container" 2>&1 | grep -q 'file-only' \
  && fail "docker logs is showing file-only content — fixture assumption broken"

echo "PASS: answers.env correctly distinguishes docker logs' stdout capture from this container's separate on-disk log file"
