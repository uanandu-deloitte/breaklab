#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

DEPLOY=/opt/pipeline/deploy.sh
MARKER=/workspace/deploy-status.txt

# Known only to the lab's own test harness and this lab's README — never
# something the learner has to extract from inside the container.
CORRECT_TOKEN="brk-ci-7f2a9d41c8b3"
WRONG_TOKEN="brk-ci-definitely-wrong"
ORIGINAL_LITERAL="sk-fake-hardcoded-abc123"

[ -f "$DEPLOY" ] || fail "$DEPLOY is missing"

# (a) The original hardcoded literal must be gone from deploy.sh's source.
grep -qF "$ORIGINAL_LITERAL" "$DEPLOY" && fail "deploy.sh still contains the original hardcoded token literal — read it from \$DEPLOY_TOKEN instead"

# Replacing one hardcoded literal with the *real* token hardcoded in its
# place is still baking the real secret into the image — reject that too.
grep -qF "$CORRECT_TOKEN" "$DEPLOY" && fail "deploy.sh contains the real deploy token as a literal in its own source — it must only be supplied at run time via \$DEPLOY_TOKEN, never hardcoded"

# (b) Correct token via env -> deploy succeeds.
rm -f "$MARKER"
if ! DEPLOY_TOKEN="$CORRECT_TOKEN" bash "$DEPLOY" >/tmp/deploy-correct.log 2>&1; then
  fail "deploy.sh failed with the correct DEPLOY_TOKEN set (see /tmp/deploy-correct.log)"
fi
[ -f "$MARKER" ] && grep -qi "DEPLOYED" "$MARKER" || fail "deploy.sh with the correct token did not produce a DEPLOYED marker at $MARKER"

# (c) Wrong token via env -> deploy must NOT succeed (proves the check
# isn't trivially satisfiable by deleting the comparison logic).
rm -f "$MARKER"
DEPLOY_TOKEN="$WRONG_TOKEN" bash "$DEPLOY" >/tmp/deploy-wrong.log 2>&1
[ -f "$MARKER" ] && fail "deploy.sh with a WRONG token still produced a DEPLOYED marker — the check isn't real authentication"

# (c) No token at all -> deploy must NOT succeed.
rm -f "$MARKER"
( unset DEPLOY_TOKEN; bash "$DEPLOY" ) >/tmp/deploy-none.log 2>&1
[ -f "$MARKER" ] && fail "deploy.sh with NO token set still produced a DEPLOYED marker"

echo "PASS: deploy.sh reads DEPLOY_TOKEN at runtime, never hardcodes the secret, and only deploys with the correct token"
