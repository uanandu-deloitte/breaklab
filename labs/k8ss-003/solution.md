# Solution: Secrets Aren't Config

```bash
# Read the value that's currently sitting in the ConfigMap.
token=$(kubectl get configmap app-config -o jsonpath='{.data.API_TOKEN}')
echo "$token"        # tok_live_8f3ac2b9e1

# Create a Secret with the same key, holding the same value.
kubectl create secret generic app-secrets --from-literal="API_TOKEN=${token}"

# Repoint the Deployment's API_TOKEN env var at the Secret.
kubectl patch deployment app --type=json -p='[
  {"op":"replace","path":"/spec/template/spec/containers/0/env/1/valueFrom","value":{"secretKeyRef":{"name":"app-secrets","key":"API_TOKEN"}}}
]'

# Remove the sensitive key from the ConfigMap entirely.
kubectl patch configmap app-config --type=json -p='[
  {"op":"remove","path":"/data/API_TOKEN"}
]'

kubectl rollout status deployment/app --timeout=60s

# Confirm the live Pod still resolves the correct value.
kubectl exec deploy/app -- printenv API_TOKEN     # tok_live_8f3ac2b9e1

# Observational step: decode the Secret yourself.
kubectl get secret app-secrets -o jsonpath='{.data.API_TOKEN}' | base64 -d
# tok_live_8f3ac2b9e1  — plaintext, instantly, no key or passphrase required

cat <<'EOF' > /workspace/answers.env
BASE64_IS_NOT_ENCRYPTION_PROOF=tok_live_8f3ac2b9e1
EOF
```

The value itself never changed — what changed is which resource type
holds it and how the Deployment reaches it. The ConfigMap no longer has
any trace of `API_TOKEN`; the Secret does, under the same key, so the
Deployment's env var resolves identically from the Pod's point of view.

The last step is the point of the lab: decoding `app-secrets` took one
pipe to `base64 -d`, no credentials beyond whatever let you run
`kubectl get secret` at all. That's not a bug — it's exactly what
base64 is for (representing bytes as text), and it means a Secret's real
protection was never the encoding. It's whoever has RBAC `get`/`list` on
`secrets` in this namespace. Moving a value out of a ConfigMap and into a
Secret only accomplishes something if access to Secret objects is
actually restricted — which is precisely the kind of boundary a Role
like the one in the RBAC lab earlier in this stage can enforce.
