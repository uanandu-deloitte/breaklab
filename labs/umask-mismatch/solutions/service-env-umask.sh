#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore a sane umask for this service.
printf 'umask 022\n' > /etc/default/producer
