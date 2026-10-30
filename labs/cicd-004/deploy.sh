#!/usr/bin/env bash
set -euo pipefail

# BUG: both the "submitted" token and the "expected" token are literals
# baked directly into this script. Anyone who can read this file can
# deploy — the comparison below always agrees with itself.
DEPLOY_TOKEN="sk-fake-hardcoded-abc123"
EXPECTED="sk-fake-hardcoded-abc123"

if [ "$DEPLOY_TOKEN" = "$EXPECTED" ]; then
  echo "DEPLOYED" > /workspace/deploy-status.txt
  echo "deploy authorized, shipping release"
else
  echo "DEPLOY DENIED: invalid token" >&2
  exit 1
fi
