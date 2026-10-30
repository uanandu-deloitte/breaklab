#!/usr/bin/env bash
set -uo pipefail

zombie_count() {
  ps -eo stat= | tr -d ' ' | grep -c '^Z'
}

first=$(zombie_count)
sleep 3
second=$(zombie_count)

# A correctly-fixed worker still shows a zombie or two at any instant
# (the reap loop runs once per cycle, so there's always a brief window
# before the most recent child gets collected) -- the real signal is
# that the count isn't climbing, and stays within a sane bound.
if [ "$first" -le 5 ] && [ "$second" -le "$((first + 1))" ]; then
  echo "PASS: zombie count isn't growing ($first, then $second)"
else
  echo "FAIL: zombie processes are still piling up ($first, then $second)"
  exit 1
fi
