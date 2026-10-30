# ConfigMaps, Two Ways

The `app-config` ConfigMap has two keys: `LOG_LEVEL` and `app.conf`.
The `configmap-demo` Deployment is supposed to expose each of them a
different way — but right now its Pod's container doesn't see either
one.

## Tasks

1. `kubectl get configmap app-config -o yaml` — see both keys the
   ConfigMap actually holds.
2. `kubectl exec <pod> -- printenv APP_LOG_LEVEL` and
   `kubectl exec <pod> -- cat /etc/app/app.conf` — confirm neither
   currently exists inside the Pod.
3. Fix the `configmap-demo` Deployment (edit the manifest and
   re-apply, or `kubectl patch`) so that:
   - the ConfigMap's `LOG_LEVEL` key appears inside the container as
     an environment variable named `APP_LOG_LEVEL`
   - the ConfigMap's `app.conf` key appears as file content at
     `/etc/app/app.conf`
4. Re-run both `kubectl exec` commands from step 2 and confirm they
   now return the ConfigMap's real values.

Get a shell: `docker exec -it breaklab-k8sd-001 bash`
