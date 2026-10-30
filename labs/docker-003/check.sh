#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

live_cmdline=$(tr '\0' ' ' < /proc/1/cmdline | sed -e 's/[[:space:]]*$//')

[ -n "${PID1_COMMAND:-}" ] || fail "PID1_COMMAND is empty"
got=$(echo "$PID1_COMMAND" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')
[ "$got" = "$live_cmdline" ] || fail "PID1_COMMAND ($got) doesn't match the live /proc/1/cmdline ($live_cmdline)"

[ -n "${IS_SHELL_FORM:-}" ] || fail "IS_SHELL_FORM is empty"
echo "$IS_SHELL_FORM" | grep -qi '^yes$' || fail "IS_SHELL_FORM should be yes — this lab's Dockerfile uses CMD's shell form, so PID 1 is /bin/sh wrapping sleep, not sleep itself"

# Ground truth: PID 1 really is a shell (this lab's Dockerfile uses CMD's
# shell form on purpose), and the real command runs as a separate child.
pid1_comm=$(ps -p 1 -o comm= | tr -d '[:space:]')
case "$pid1_comm" in
  sh|bash|dash) : ;;
  *) fail "PID 1's comm ($pid1_comm) isn't a shell — did the Dockerfile's CMD change to exec form?" ;;
esac

ps -eo pid,ppid,comm | awk '$2 == 1 && $3 == "sleep" {found=1} END {exit !found}' \
  || fail "no 'sleep' process found as a direct child of PID 1 — expected the shell (PID 1) to have forked it"

echo "PASS: answers.env correctly identifies PID 1 as a shell wrapping the real command"
