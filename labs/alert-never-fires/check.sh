#!/usr/bin/env bash
set -uo pipefail
: > /var/log/alerts.log
alert-evaluator
if grep -q "disk_used_pct=97" /var/log/alerts.log 2>/dev/null; then
  echo "PASS: alert fired for the ongoing disk_used_pct=97 breach"
else
  echo "FAIL: no alert fired for the ongoing disk_used_pct=97 breach"
  exit 1
fi
