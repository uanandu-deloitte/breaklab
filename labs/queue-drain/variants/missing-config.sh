#!/usr/bin/env bash
set -euo pipefail
# Oops: a cleanup pass "fixed" a path and broke the worker config
sed -i 's#/usr/local/libexec/worker.sh#/usr/local/libexec/worker#' /etc/jobs/worker.conf
