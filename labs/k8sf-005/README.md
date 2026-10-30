# Readiness Probes Gate Traffic, Not Liveness

A `catalog-api` Deployment (1 replica) and its Service are running. The
Pod has been `Running` the whole time and has never restarted — but its
Service has nowhere to send traffic.

## Look around first

- `kubectl get pods -o wide` — is the Pod Running? What does the READY
  column say?
- `kubectl describe pod -l app=catalog-api` — what does the Events
  section say about the readiness probe specifically?
- `kubectl get endpoints catalog-api` — how many addresses does the
  Service actually have to route to?
- `kubectl get deployment catalog-api -o
  jsonpath='{.spec.template.spec.containers[0].readinessProbe}'` — what
  path and port is the probe actually configured to check?
- `kubectl exec deploy/catalog-api -- wget -qO- localhost:8080/` — what
  does the container actually serve, and on which port?

## Record your findings

Write these to `/workspace/answers.env`, one `KEY=value` per line:

```
READY_COLUMN=
ENDPOINT_COUNT=
PROBE_PATH=
```

- `READY_COLUMN`: the exact READY column value from `kubectl get pods
  -l app=catalog-api` (e.g. `0/1`).
- `ENDPOINT_COUNT`: how many addresses `kubectl get endpoints
  catalog-api` actually lists right now.
- `PROBE_PATH`: the readiness probe's configured path, read live from
  the Deployment, not guessed.

Get a shell: `docker exec -it breaklab-k8sf-005 bash`
