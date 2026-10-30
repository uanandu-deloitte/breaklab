#!/usr/bin/env bash
set -euo pipefail
# Oops: an edit to the apply script dropped REPLICAS from its substitution list
cat > /usr/local/bin/infra-apply-worker <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
replicas=$(sed -n 's/^replicas: *//p' /infra/worker.yaml)
queue=$(sed -n 's/^queue: *//p' /infra/worker.yaml)
mkdir -p /etc/worker
VARS='${QUEUE}'
REPLICAS="$replicas" QUEUE="$queue" envsubst "$VARS" < /infra/templates/worker.env.tmpl > /etc/worker/worker.env
EOF
chmod 755 /usr/local/bin/infra-apply-worker
infra-apply-worker
