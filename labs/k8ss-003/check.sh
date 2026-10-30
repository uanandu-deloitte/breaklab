#!/usr/bin/env bash
set -uo pipefail

fail() { echo "FAIL: $1"; exit 1; }

EXPECTED_TOKEN="tok_live_8f3ac2b9e1"
NS=default

# 1. The ConfigMap must no longer carry the sensitive key at all.
cm_token=$(kubectl get configmap app-config -n "$NS" -o jsonpath='{.data.API_TOKEN}' 2>/dev/null)
[ -z "$cm_token" ] || fail "ConfigMap app-config still has an API_TOKEN key — it must be removed entirely"

# 2. A Secret must now hold it, with the correct value.
secret_b64=$(kubectl get secret app-secrets -n "$NS" -o jsonpath='{.data.API_TOKEN}' 2>/dev/null)
[ -n "$secret_b64" ] || fail "Secret app-secrets with key API_TOKEN not found"
secret_val=$(echo "$secret_b64" | base64 -d 2>/dev/null)
[ "$secret_val" = "$EXPECTED_TOKEN" ] || fail "Secret app-secrets' API_TOKEN decodes to '$secret_val', expected '$EXPECTED_TOKEN'"

# 3. The Deployment must consume it via secretKeyRef, not configMapKeyRef.
secret_ref_name=$(kubectl get deployment app -n "$NS" -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="API_TOKEN")].valueFrom.secretKeyRef.name}' 2>/dev/null)
secret_ref_key=$(kubectl get deployment app -n "$NS" -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="API_TOKEN")].valueFrom.secretKeyRef.key}' 2>/dev/null)
[ "$secret_ref_name" = "app-secrets" ] && [ "$secret_ref_key" = "API_TOKEN" ] \
  || fail "Deployment app's API_TOKEN env var isn't wired to secretKeyRef app-secrets/API_TOKEN (got name='$secret_ref_name' key='$secret_ref_key')"

cm_ref_name=$(kubectl get deployment app -n "$NS" -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="API_TOKEN")].valueFrom.configMapKeyRef.name}' 2>/dev/null)
[ -z "$cm_ref_name" ] || fail "Deployment app's API_TOKEN env var still references a configMapKeyRef ($cm_ref_name)"

# 4. The live Pod must actually resolve the correct value at runtime.
pod=$(kubectl get pod -l app=app -n "$NS" -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
[ -n "$pod" ] || fail "no running Pod found for app=app"
live_val=$(kubectl exec "$pod" -n "$NS" -- printenv API_TOKEN 2>/dev/null)
[ "$live_val" = "$EXPECTED_TOKEN" ] || fail "live Pod's API_TOKEN env resolves to '$live_val', expected '$EXPECTED_TOKEN'"

# 5. The learner's own recorded proof that base64 isn't confidentiality.
ANSWERS=/workspace/answers.env
[ -f "$ANSWERS" ] || fail "$ANSWERS not found"
# shellcheck disable=SC1090
source "$ANSWERS"
[ -n "${BASE64_IS_NOT_ENCRYPTION_PROOF:-}" ] || fail "BASE64_IS_NOT_ENCRYPTION_PROOF is empty in $ANSWERS"
echo "$BASE64_IS_NOT_ENCRYPTION_PROOF" | grep -qF "$EXPECTED_TOKEN" \
  || fail "BASE64_IS_NOT_ENCRYPTION_PROOF doesn't contain the actual decoded token — looks like a placeholder, not something you really decoded"

echo "PASS: the sensitive value lives in a Secret, the Deployment reads it via secretKeyRef, and the ConfigMap no longer carries it"
