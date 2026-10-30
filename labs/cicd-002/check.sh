#!/usr/bin/env bash
set -uo pipefail

RUN_SH=/opt/pipeline/run.sh
TEST_SH=/opt/pipeline/test.sh
MARKER=/workspace/bad-deploy-happened
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$RUN_SH" ] || fail "/opt/pipeline/run.sh is missing — fixture was tampered with"
[ -f "$TEST_SH" ] || fail "/opt/pipeline/test.sh is missing — fixture was tampered with"

content=$(cat "$RUN_SH")

has_errexit() {
  echo "$content" | grep -qE '(^|[^a-zA-Z])set[[:space:]]+-[a-zA-Z]*e'
}

# The original bug: build.sh and/or test.sh followed by a bare ';' instead
# of '&&', with no errexit either — a failure in that stage has no effect
# on whether the next one runs.
has_bad_semicolon() {
  echo "$content" | grep -qE '(build|test)\.sh[[:space:]]*;'
}

if has_bad_semicolon && ! has_errexit; then
  fail "run.sh still chains build.sh/test.sh with a bare ';' and has no 'set -e' — a failing stage can't halt the next one this way"
fi

# The fix: both build.sh and test.sh followed by '&&' (so the next stage
# only runs if the previous one succeeded), or 'set -e' near the top.
and_count=$(echo "$content" | grep -oE '(build|test)\.sh[[:space:]]*&&' | wc -l | tr -d ' ')
if ! has_errexit && [ "$and_count" -lt 2 ]; then
  fail "run.sh doesn't chain both build.sh and test.sh with '&&' (and doesn't use 'set -e' either) — a failing stage still won't halt the pipeline"
fi

# Sanity: test.sh must still genuinely fail on its own. This lab is about
# run.sh's reaction to a real failure, not about making test.sh pass.
if "$TEST_SH" >/dev/null 2>&1; then
  fail "test.sh no longer fails on its own — it's fixture data for this lab and must keep failing; fix run.sh's handling of the failure instead"
fi

# Dynamic check: remove any stale marker, re-run run.sh fresh, and confirm
# the marker is genuinely absent afterward — not just absent by accident
# of leftover state.
rm -f "$MARKER"
"$RUN_SH" >/tmp/cicd-002-run-output.log 2>&1
run_exit=$?

if [ -f "$MARKER" ]; then
  fail "$MARKER exists after a fresh run.sh run — deploy.sh still ran even though test.sh failed"
fi

[ "$run_exit" -ne 0 ] || fail "run.sh exited 0 even though test.sh failed — the pipeline should report failure, not success"

echo "PASS: run.sh halts after test.sh fails, and deploy.sh never runs on a fresh invocation"
