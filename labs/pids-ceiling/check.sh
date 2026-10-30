#!/usr/bin/env bash
# Runs on the host via docker cp, not docker exec: when the container is
# genuinely out of pids, exec itself can fail to fork (the same symptom
# this lab simulates), so the check must not depend on it. Submits a
# batch of fresh jobs rather than just one, since a single lucky fork
# can succeed even while the box is mostly out of room.
set -uo pipefail
n=10
prefix="probe-$$"
tmpfile=$(mktemp)
touch "$tmpfile"
for i in $(seq 1 "$n"); do
  docker cp "$tmpfile" "${BREAKLAB_CONTAINER}:/var/spool/jobs/${prefix}-${i}.job" >/dev/null 2>&1
done
rm -f "$tmpfile"

completed=0
for round in $(seq 1 15); do
  snap=$(mktemp -d)
  if docker cp "${BREAKLAB_CONTAINER}:/var/spool/jobs/done" "$snap/done" >/dev/null 2>&1; then
    completed=$(ls "$snap/done" 2>/dev/null | grep -c "^${prefix}-" || true)
  fi
  rm -rf "$snap"
  if [ "$completed" -ge "$n" ]; then
    break
  fi
  sleep 1
done

threshold=$((n * 7 / 10))
if [ "$completed" -ge "$threshold" ]; then
  echo "PASS: $completed/$n freshly submitted jobs completed"
  exit 0
else
  echo "FAIL: only $completed/$n freshly submitted jobs completed"
  exit 1
fi
