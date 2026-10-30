#!/usr/bin/env bash
set -uo pipefail

ANSWERS=/workspace/answers.env
fail() { echo "FAIL: $1"; exit 1; }

[ -f "$ANSWERS" ] || fail "/workspace/answers.env not found"

# shellcheck disable=SC1090
source "$ANSWERS"

[ -f /var/run/reader.pid ] || fail "/var/run/reader.pid is missing — the lab environment may not have started correctly"
pid=$(cat /var/run/reader.pid)
[ -d "/proc/$pid" ] || fail "reader process (pid $pid) is not running — restart the lab environment"

# Ground truth: find which path the reader's log-reading descriptor
# currently resolves to. It should be the rotated-away app.log.1, not the
# live app.log, since tail -f follows the descriptor it originally opened.
fd_target=""
for fd in /proc/"$pid"/fd/*; do
  tgt=$(readlink -f "$fd" 2>/dev/null || true)
  case "$tgt" in
    /var/log/app.log*) fd_target="$tgt" ;;
  esac
done
[ -n "$fd_target" ] || fail "could not find the reader's log-reading file descriptor under /proc/$pid/fd"

[ -f /var/log/app.log ] || fail "/var/log/app.log is missing"
[ -f "$fd_target" ] || fail "the file the reader has open ($fd_target) no longer exists"

current_inode=$(stat -c %i /var/log/app.log 2>/dev/null || echo "")
fd_inode=$(stat -c %i "$fd_target" 2>/dev/null || echo "")
[ -n "$current_inode" ] && [ -n "$fd_inode" ] || fail "could not stat the log files to compare inodes"

# Compare actual content, not just line counts: a rotated-in file can
# coincidentally have the same number of lines as the stale reader
# output while holding completely different data, so a count-only check
# would miss exactly the failure mode this lab is about.
ground_truth_stale="no"
diff -q /var/log/app.log /var/log/reader-output.log >/dev/null 2>&1 || ground_truth_stale="yes"

ground_truth_inodes_match="no"
[ "$current_inode" = "$fd_inode" ] && ground_truth_inodes_match="yes"

[ -n "${READER_STALE:-}" ] || fail "READER_STALE is empty"
echo "$READER_STALE" | grep -qi "^${ground_truth_stale}\$" \
  || fail "READER_STALE is '$READER_STALE', but reader-output.log's content ($(wc -l < /var/log/reader-output.log) line(s)) vs app.log's content ($(wc -l < /var/log/app.log) line(s)) says otherwise — expected '$ground_truth_stale'"

[ -n "${READER_FD_PATH:-}" ] || fail "READER_FD_PATH is empty"
[ "$READER_FD_PATH" = "$fd_target" ] \
  || fail "READER_FD_PATH is '$READER_FD_PATH', but the reader process actually has '$fd_target' open"

[ -n "${INODES_MATCH:-}" ] || fail "INODES_MATCH is empty"
echo "$INODES_MATCH" | grep -qi "^${ground_truth_inodes_match}\$" \
  || fail "INODES_MATCH is '$INODES_MATCH', expected '$ground_truth_inodes_match' (current app.log inode=$current_inode, reader's open file inode=$fd_inode)"

echo "PASS: correctly diagnosed the reader as stuck on the rotated-away file ($fd_target) instead of the current app.log"
