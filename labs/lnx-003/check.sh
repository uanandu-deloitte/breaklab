#!/usr/bin/env bash
set -uo pipefail

TRANSCRIPT=/workspace/transcript.txt
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$TRANSCRIPT" ] || fail "/workspace/transcript.txt not found"

home=$(eval echo ~"$(whoami)")
required=(/etc /var/log /opt "$home")

last_line=0
for path in "${required[@]}"; do
  line=$(grep -nxF "$path" "$TRANSCRIPT" | head -1 | cut -d: -f1)
  if [ -z "$line" ]; then
    fail "transcript is missing $path"
  fi
  if [ "$line" -le "$last_line" ]; then
    fail "$path appears out of order in the transcript"
  fi
  last_line=$line
done

echo "PASS: transcript visits every required location in order"
