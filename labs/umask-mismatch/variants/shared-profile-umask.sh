#!/usr/bin/env bash
set -euo pipefail
# Oops: a platform-wide "harden everything" profile snippet, meant for
# interactive shells, now also gets sourced before this service writes
# its export — silently overriding its own (fine) umask.
printf 'umask 077\n' > /etc/profile.d/10-harden.sh
