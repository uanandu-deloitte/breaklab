#!/usr/bin/env bash
set -uo pipefail
live=$(cat /var/app/data/marker.txt 2>/dev/null)
backed_up=$(tar xzfO /var/backups/latest.tar.gz ./marker.txt 2>/dev/null)
if [ -n "$live" ] && [ "$live" = "$backed_up" ]; then
  echo "PASS: backup is current ($backed_up)"
else
  echo "FAIL: backup has '${backed_up:-nothing}', live data is '${live:-nothing}'"
  exit 1
fi
