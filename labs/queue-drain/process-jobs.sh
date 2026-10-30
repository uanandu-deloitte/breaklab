#!/usr/bin/env bash
set -euo pipefail
conf=/etc/jobs/worker.conf
if [ ! -r "$conf" ]; then
  echo "process-jobs: cannot read $conf" >&2
  exit 1
fi
. "$conf"
if [ -z "${WORKER_BIN:-}" ] || [ ! -x "$WORKER_BIN" ]; then
  echo "process-jobs: WORKER_BIN is not set to an executable (got '${WORKER_BIN:-}')" >&2
  exit 1
fi
shopt -s nullglob
for job in /var/spool/jobs/*.json; do
  su -s /bin/sh jobsrunner -c "$WORKER_BIN '$job'"
done
