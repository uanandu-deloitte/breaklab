#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

ANSWERS=/workspace/answers.env
[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"
# shellcheck disable=SC1090
source "$ANSWERS"

[ -n "${STAGED_CHANGE_COMMAND:-}" ] || fail "STAGED_CHANGE_COMMAND is empty"
echo "$STAGED_CHANGE_COMMAND" | grep -qi 'restore' || fail "STAGED_CHANGE_COMMAND doesn't look like git restore"

[ -n "${PRIVATE_COMMIT_COMMAND:-}" ] || fail "PRIVATE_COMMIT_COMMAND is empty"
echo "$PRIVATE_COMMIT_COMMAND" | grep -qi 'reset' || fail "PRIVATE_COMMIT_COMMAND doesn't look like git reset"

[ -n "${SHARED_COMMIT_COMMAND:-}" ] || fail "SHARED_COMMIT_COMMAND is empty"
echo "$SHARED_COMMIT_COMMAND" | grep -qi 'revert' || fail "SHARED_COMMIT_COMMAND doesn't look like git revert"

# --- Situation A: repo-staged ---
REPO_A=/workspace/repo-staged
[ -d "$REPO_A/.git" ] || fail "$REPO_A is missing"
[ -z "$(git -C "$REPO_A" status --porcelain)" ] || fail "$REPO_A working tree isn't clean -- the staged change is still there in some form"
grep -q "mode: safe" "$REPO_A/config.yaml" 2>/dev/null || fail "$REPO_A/config.yaml doesn't match the committed baseline (mode: safe)"
grep -q "mode: debug" "$REPO_A/config.yaml" 2>/dev/null && fail "$REPO_A/config.yaml still has the unwanted debug change"

# --- Situation B: repo-private ---
REPO_B=/workspace/repo-private
[ -d "$REPO_B/.git" ] || fail "$REPO_B is missing"
baseline_sha=$(cat /srv/.repo-private-baseline-sha 2>/dev/null) || fail "missing fixture record of repo-private's baseline commit"
head_sha=$(git -C "$REPO_B" rev-parse HEAD 2>/dev/null) || fail "$REPO_B has no HEAD"
[ "$head_sha" = "$baseline_sha" ] || fail "$REPO_B's HEAD ($head_sha) doesn't match its pre-bad-commit baseline ($baseline_sha) -- the bad commit should be gone from history entirely"
[ -z "$(git -C "$REPO_B" status --porcelain)" ] || fail "$REPO_B working tree isn't clean"
grep -q "BAD-COMMIT-MARKER" "$REPO_B/settings.yaml" 2>/dev/null && fail "$REPO_B/settings.yaml still has the bad commit's content"

# --- Situation C: repo-shared ---
REPO_C=/workspace/repo-shared
[ -d "$REPO_C/.git" ] || fail "$REPO_C is missing"
bad_sha=$(cat /srv/.repo-shared-bad-sha 2>/dev/null) || fail "missing fixture record of repo-shared's bad commit"
head_sha_c=$(git -C "$REPO_C" rev-parse HEAD 2>/dev/null) || fail "$REPO_C has no HEAD"
[ "$head_sha_c" != "$bad_sha" ] || fail "$REPO_C's HEAD is still the bad commit itself -- you need a new commit that undoes it, not a no-op"
git -C "$REPO_C" merge-base --is-ancestor "$bad_sha" HEAD \
  || fail "$REPO_C no longer contains the bad commit in its history -- it's already shared, so history must not be rewritten (that's what git reset would have done)"
grep -q "BAD-COMMIT-MARKER" "$REPO_C/settings.yaml" 2>/dev/null && fail "$REPO_C/settings.yaml still shows the bad commit's change -- it should be reverted"
[ -z "$(git -C "$REPO_C" status --porcelain)" ] || fail "$REPO_C working tree isn't clean"

echo "PASS: all three situations were fixed with the right command"
