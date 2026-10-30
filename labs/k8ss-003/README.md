# Secrets Aren't Config

There's a `ConfigMap` named `app-config` in the `default` namespace with
two keys: `LOG_LEVEL` (fine, that's genuinely config) and `API_TOKEN` (not
fine — that's a sensitive value stored as plain config). A `Deployment`
named `app` consumes both via `configMapKeyRef` into its one container's
environment.

## Tasks

1. Read the current `API_TOKEN` value out of the ConfigMap.
2. Create a `Secret` named `app-secrets` with the same key, `API_TOKEN`,
   holding that same value.
3. Update the `app` Deployment so its `API_TOKEN` env var resolves via
   `secretKeyRef` (pointing at `app-secrets` / `API_TOKEN`) instead of
   `configMapKeyRef`. Leave `LOG_LEVEL` exactly as it is.
4. Remove the `API_TOKEN` key from `app-config` entirely — it shouldn't
   exist there anymore, in any form.
5. Confirm the running Pod still reads the correct value: `kubectl exec`
   into it and check the `API_TOKEN` environment variable.
6. As a final step, decode the Secret yourself:
   ```
   kubectl get secret app-secrets -o jsonpath='{.data.API_TOKEN}' | base64 -d
   ```
   Notice that this instantly reveals the plaintext value — a Secret's
   base64 encoding is not encryption, it's just how the API object
   represents arbitrary bytes as text. Record what you observed into
   `/workspace/answers.env`:
   ```
   BASE64_IS_NOT_ENCRYPTION_PROOF=<the exact decoded value you saw>
   ```

## Useful commands

```
kubectl get configmap app-config -o yaml
kubectl get deployment app -o yaml
kubectl create secret generic app-secrets --from-literal=API_TOKEN=<value>
kubectl patch deployment app --type=json -p='[{"op":"replace","path":"/spec/template/spec/containers/0/env/1/valueFrom","value":{"secretKeyRef":{"name":"app-secrets","key":"API_TOKEN"}}}]'
kubectl patch configmap app-config --type=json -p='[{"op":"remove","path":"/data/API_TOKEN"}]'
kubectl exec deploy/app -- printenv API_TOKEN
kubectl get secret app-secrets -o jsonpath='{.data.API_TOKEN}' | base64 -d
```

Get a shell: `docker exec -it breaklab-k8ss-003 bash`
