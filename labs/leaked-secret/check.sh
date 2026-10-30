#!/usr/bin/env bash
set -uo pipefail
cd /srv/repo || { echo "FAIL: repo missing"; exit 1; }
if git log --all -p 2>/dev/null | grep "sk-FAKE" >/dev/null; then
  echo "FAIL: the secret is still present somewhere in history"
  exit 1
else
  echo "PASS: the secret is gone from all history"
fi
