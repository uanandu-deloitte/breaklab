#!/usr/bin/env bash
set -euo pipefail

token=$(kubectl get configmap app-config -n default -o jsonpath='{.data.API_TOKEN}')

# Move the sensitive value into a proper Secret, same key.
kubectl create secret generic app-secrets -n default \
  --from-literal="API_TOKEN=${token}" >/dev/null

# Point the Deployment's API_TOKEN env var at the Secret instead of the ConfigMap.
kubectl patch deployment app -n default --type=json -p='[
  {"op":"replace","path":"/spec/template/spec/containers/0/env/1/valueFrom","value":{"secretKeyRef":{"name":"app-secrets","key":"API_TOKEN"}}}
]' >/dev/null

# Remove the sensitive key from the ConfigMap entirely.
kubectl patch configmap app-config -n default --type=json -p='[
  {"op":"remove","path":"/data/API_TOKEN"}
]' >/dev/null

kubectl rollout status deployment/app -n default --timeout=60s >/dev/null

# Observational step: base64 is encoding, not encryption — prove it by
# decoding the Secret's own stored value and recording it.
decoded=$(kubectl get secret app-secrets -n default -o jsonpath='{.data.API_TOKEN}' | base64 -d)
cat <<EOF > /workspace/answers.env
BASE64_IS_NOT_ENCRYPTION_PROOF=${decoded}
EOF
