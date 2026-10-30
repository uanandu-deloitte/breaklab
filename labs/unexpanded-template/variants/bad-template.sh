#!/usr/bin/env bash
set -euo pipefail
# Oops: an edit to the template mistyped the variable name
sed -i 's/\${QUEUE}/\${QUEUE_NAME}/' /infra/templates/worker.env.tmpl
infra-apply-worker
