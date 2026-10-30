#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

cd /workspace/project || fail "/workspace/project is missing"
[ -f .github/workflows/ci.yml ] || fail ".github/workflows/ci.yml is missing"

rm -f /workspace/dist/app.txt /workspace/deployed.txt

run_out=$(act push -P ubuntu-latest=-self-hosted 2>&1)
run_exit=$?
[ "$run_exit" -eq 0 ] || fail "act push failed (exit $run_exit):
$run_out"

[ -f /workspace/dist/app.txt ] || fail "build job never produced /workspace/dist/app.txt"
[ -f /workspace/deployed.txt ] || fail "deploy job never ran -- /workspace/deployed.txt doesn't exist (it's still being skipped)"

content=$(cat /workspace/deployed.txt)
[ "$content" = "built" ] || fail "/workspace/deployed.txt has unexpected content: $content"

echo "PASS: both jobs ran for real via act and deploy actually produced /workspace/deployed.txt"
