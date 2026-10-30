#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

cd /workspace/project || fail "/workspace/project is missing"
[ -f .gitlab-ci.yml ] || fail ".gitlab-ci.yml is missing"

rm -f /workspace/dist/app.txt /workspace/deployed.txt

run_out=$(gitlab-ci-local --force-shell-executor 2>&1)
run_exit=$?
[ "$run_exit" -eq 0 ] || fail "gitlab-ci-local failed (exit $run_exit):
$run_out"

[ -f /workspace/dist/app.txt ] || fail "build job never produced /workspace/dist/app.txt"
[ -f /workspace/deployed.txt ] || fail "deploy job never ran -- /workspace/deployed.txt doesn't exist (it's still being skipped)"

content=$(cat /workspace/deployed.txt)
[ "$content" = "built" ] || fail "/workspace/deployed.txt has unexpected content: $content"

echo "PASS: both stages ran for real via gitlab-ci-local and deploy actually produced /workspace/deployed.txt"
