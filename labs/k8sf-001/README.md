# Your First Pod

A small `hello` Deployment (2 replicas) and its Service are already running
and healthy in this cluster. Nothing here is broken — the point of this lab
is to get comfortable reading live cluster state before you start fixing it
in later labs.

## Look around first

Work through these using plain `kubectl`, in whatever order you like:

- `kubectl get pods -o wide` — are the `hello` Pods Running? How many are
  Ready?
- `kubectl describe pod <one of the hello pods>` — what does the Events
  section say about how it started?
- `kubectl logs <one of the hello pods>` — what did the container print?
- `kubectl get deployments` — what does it show for `hello`'s desired vs.
  ready replica count?
- `kubectl get svc hello` — what ClusterIP and port does the Service have?

## Now make a real change

Scale `hello` from 2 replicas to **4**:

```
kubectl scale deployment/hello --replicas=4
```

Confirm all 4 actually become Ready before moving on — `kubectl rollout
status deployment/hello` or `kubectl get pods -l app=hello -w` both show
this.

## Record your findings

Write these to `/workspace/answers.env`, one `KEY=value` per line:

```
INITIAL_REPLICA_COUNT=
CONTAINER_IMAGE=
FINAL_REPLICA_COUNT=
```

- `INITIAL_REPLICA_COUNT`: how many replicas `hello` had before you scaled
  it.
- `CONTAINER_IMAGE`: the exact image (with tag) the `hello` containers run —
  read it from a live `kubectl get`, don't guess from this file.
- `FINAL_REPLICA_COUNT`: the replica count you scaled it to.

Get a shell: `docker exec -it breaklab-k8sf-001 bash`
