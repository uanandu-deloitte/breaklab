#!/usr/bin/env bash
set -euo pipefail
# Oops: "security hardening" tightened this service's own umask so hard
# it also locks group members out of its shared export directory.
printf 'umask 077\n' > /etc/default/producer
