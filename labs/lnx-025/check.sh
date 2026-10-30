#!/usr/bin/env bash
set -uo pipefail
fail() { echo "FAIL: $1"; exit 1; }

SCRIPT=/opt/app/backup.sh
[ -x "$SCRIPT" ] || fail "/opt/app/backup.sh is missing or not executable"

rm -f /opt/backups/backup-*.tar.gz 2>/dev/null || true

out=$("$SCRIPT" 2>&1)
code=$?
if [ "$code" -ne 0 ]; then
  fail "backup.sh exited non-zero on a healthy source: $out"
fi

archive=$(ls -t /opt/backups/backup-*.tar.gz 2>/dev/null | head -1)
[ -n "$archive" ] || fail "no backup archive was created under /opt/backups"

if ! tar -tf "$archive" | grep -q "report draft.txt"; then
  fail "the archive doesn't contain the file with a space in its name"
fi

echo "$out" | grep -qi "succeeded" || fail "backup.sh didn't report success on a genuine success"

# Now prove it fails safely: hide the source directory and rerun.
mv /data/to-backup /data/to-backup.hidden
fail_out=$("$SCRIPT" 2>&1)
fail_code=$?
mv /data/to-backup.hidden /data/to-backup

if [ "$fail_code" -eq 0 ]; then
  fail "backup.sh exited 0 even though the source directory was missing"
fi
if echo "$fail_out" | grep -qi "succeeded"; then
  fail "backup.sh printed a success message even though the backup failed"
fi

echo "PASS: backup.sh backs up correctly, including the space-named file, and fails safely when the source is missing"
