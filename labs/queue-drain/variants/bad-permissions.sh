#!/usr/bin/env bash
set -euo pipefail
# Oops: a permissions "hardening" pass locked the worker out of its own spool
chown root:root /var/spool/jobs/done
chmod 700 /var/spool/jobs/done
