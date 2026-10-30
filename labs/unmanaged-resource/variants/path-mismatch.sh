#!/usr/bin/env bash
set -euo pipefail
# Oops: this was "imported" with a trailing slash that doesn't match the
# exact path infra-apply actually targets
jq '.managed = ["/etc/motd.conf/"]' /infra/state.json > /tmp/state.json && mv /tmp/state.json /infra/state.json
sed -i 's/^message:.*/message: Welcome v3/' /infra/resources.yaml
infra-apply
