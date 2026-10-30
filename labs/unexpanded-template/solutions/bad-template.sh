#!/usr/bin/env bash
set -euo pipefail
# Reference fix: restore the template's variable name and re-apply
sed -i 's/\${QUEUE_NAME}/\${QUEUE}/' /infra/templates/worker.env.tmpl
infra-apply-worker
