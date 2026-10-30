#!/usr/bin/env bash
set -euo pipefail
# Reference fix: point worker.conf back at the real worker binary and drain
printf 'WORKER_BIN=/usr/local/libexec/worker.sh\n' > /etc/jobs/worker.conf
process-jobs
