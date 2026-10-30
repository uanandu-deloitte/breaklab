#!/usr/bin/env bash
set -uo pipefail
cd /srv/repo || { echo "FAIL: repo missing"; exit 1; }
if git log main --oneline 2>/dev/null | grep "RELEASE: v1.2.0" >/dev/null && grep -q "v1.2.0" RELEASE.md 2>/dev/null; then
  echo "PASS: the release commit is on main"
else
  echo "FAIL: the release commit is still missing from main"
  exit 1
fi
