#!/usr/bin/env bash
set -uo pipefail
cd /srv/repo || { echo "FAIL: repo missing"; exit 1; }
if bash healthcheck.sh >/tmp/hc.out 2>&1; then
  echo "PASS: healthcheck passes on main"
else
  echo "FAIL: healthcheck still fails on main"
  exit 1
fi
