#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore the spool directory to the worker user and drain
chown jobsrunner:jobsrunner /var/spool/jobs/done
chmod 755 /var/spool/jobs/done
process-jobs
