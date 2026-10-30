#!/usr/bin/env bash
set -euo pipefail
# This quarter's message changed
sed -i 's/^message:.*/message: Q3 pricing notice v2/' /infra/site.yaml
# Simulate someone deploying it from an ordinary interactive shell,
# where /etc/bash.bashrc's legacy export quietly wins
bash -ic 'infra-apply'
