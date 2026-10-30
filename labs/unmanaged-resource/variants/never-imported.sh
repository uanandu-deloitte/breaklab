#!/usr/bin/env bash
set -euo pipefail
# Oops: this resource's path was never added to state.json
jq '.managed -= ["/etc/motd.conf"]' /infra/state.json > /tmp/state.json && mv /tmp/state.json /infra/state.json
sed -i 's/^message:.*/message: Welcome v3/' /infra/resources.yaml
infra-apply
