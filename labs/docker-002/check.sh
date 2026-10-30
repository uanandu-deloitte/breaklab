#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
SAMPLE=/workspace/sample.Dockerfile
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$SAMPLE" ] || fail "/workspace/sample.Dockerfile is missing — fixture was tampered with"
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

trim() { echo "$1" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//'; }

[ -n "${CACHE_MISS_ON_SOURCE_CHANGE:-}" ] || fail "CACHE_MISS_ON_SOURCE_CHANGE is empty"
[ "$(trim "$CACHE_MISS_ON_SOURCE_CHANGE")" = "yes" ] || fail "CACHE_MISS_ON_SOURCE_CHANGE should be yes — COPY . . runs before npm ci, so any source change busts the dependency layer's cache too"

# The sample must actually still contain the cache-busting instruction
# this question is about, before grading the learner's answer against it.
grep -qF 'COPY . .' "$SAMPLE" || fail "sample.Dockerfile no longer contains 'COPY . .' — fixture was tampered with"
[ -n "${CACHE_BUSTING_INSTRUCTION:-}" ] || fail "CACHE_BUSTING_INSTRUCTION is empty"
[ "$(trim "$CACHE_BUSTING_INSTRUCTION")" = "COPY . ." ] || fail "CACHE_BUSTING_INSTRUCTION should name the exact instruction 'COPY . .' that runs before npm ci"

[ -n "${BUILDER_TOOLS_IN_FINAL_IMAGE:-}" ] || fail "BUILDER_TOOLS_IN_FINAL_IMAGE is empty"
[ "$(trim "$BUILDER_TOOLS_IN_FINAL_IMAGE")" = "no" ] || fail "BUILDER_TOOLS_IN_FINAL_IMAGE should be no — build-essential/python3 are only ever installed in the builder stage, and no COPY --from pulls them across"

# Derive the real set of paths the final stage copies from the builder
# stage directly from the fixture, rather than hardcoding it here.
expected_copies=$(grep 'COPY --from=builder' "$SAMPLE" \
  | awk '{print $NF}' | sed -e 's#^\./##' -e 's#/$##' | sort | tr '\n' ',' | sed 's/,$//')
[ -n "$expected_copies" ] || fail "couldn't find any COPY --from=builder lines in sample.Dockerfile — fixture was tampered with"

[ -n "${FINAL_STAGE_COPIES_FROM_BUILDER:-}" ] || fail "FINAL_STAGE_COPIES_FROM_BUILDER is empty"
got_copies=$(echo "${FINAL_STAGE_COPIES_FROM_BUILDER}" | tr ',' '\n' | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' | sort | tr '\n' ',' | sed 's/,$//')
[ "$got_copies" = "$expected_copies" ] || fail "FINAL_STAGE_COPIES_FROM_BUILDER should list exactly: $expected_copies (got: $got_copies)"

[ -n "${SAME_RUN_DELETE_AVOIDS_BLOAT:-}" ] || fail "SAME_RUN_DELETE_AVOIDS_BLOAT is empty"
[ "$(trim "$SAME_RUN_DELETE_AVOIDS_BLOAT")" = "yes" ] || fail "SAME_RUN_DELETE_AVOIDS_BLOAT should be yes — Block A's file never exists when that RUN's layer diff is captured"

[ -n "${LATER_RUN_DELETE_AVOIDS_BLOAT:-}" ] || fail "LATER_RUN_DELETE_AVOIDS_BLOAT is empty"
[ "$(trim "$LATER_RUN_DELETE_AVOIDS_BLOAT")" = "no" ] || fail "LATER_RUN_DELETE_AVOIDS_BLOAT should be no — Block B's file was already committed to an earlier layer before the later RUN deletes it"

echo "PASS: answers.env correctly reads the sample Dockerfile's cache, layering, and multi-stage behavior"
