#!/usr/bin/env bash
set -euo pipefail
jq '.managed += ["/etc/motd.conf"]' /infra/state.json > /tmp/state.json && mv /tmp/state.json /infra/state.json
infra-apply
