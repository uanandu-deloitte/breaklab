# Tracing a Failure Across a Service Dependency

The `frontend` Pod is `Running` and `1/1 Ready`. A feature that
depends on `backend` is broken anyway, because `backend` currently has
zero ready endpoints. Your job is not to fix anything — it's to
correctly diagnose *which* of two possible failures this is before
you would reach for a fix.

## Tasks

1. `kubectl get pods` — confirm frontend itself looks healthy. It
   does; that's the trap.
2. `kubectl logs deployment/frontend --tail=20` — see which
   dependency frontend itself reports as unreachable.
3. `kubectl get endpoints <dependency>` — confirm it really has zero
   ready addresses.
4. `kubectl get deployment <dependency>` — compare `DESIRED` against
   `READY`. This is what tells apart the two possible causes:
   - `DESIRED` is `0` → the dependency was scaled to zero.
   - `DESIRED` is greater than `0` but `READY` is `0` → replicas exist
     but are crash-looping.
5. If `DESIRED` is greater than `0`, `kubectl describe pod -l
   app=<dependency>` to see the restart count and events confirming
   the crash loop.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
DEPENDENCY_NAME=
READY_ENDPOINTS=
DESIRED_REPLICAS=
READY_REPLICAS=
FAILURE_MODE=
```

- `DEPENDENCY_NAME`: the name of the dependency frontend's own logs
  report as unreachable.
- `READY_ENDPOINTS`: the number of ready addresses `kubectl get
  endpoints <dependency>` currently shows.
- `DESIRED_REPLICAS`: the dependency Deployment's `DESIRED` count.
- `READY_REPLICAS`: the dependency Deployment's `READY` count.
- `FAILURE_MODE`: either `scaled-to-zero` or `crash-looping` — whichever
  this incident actually is, based on `DESIRED_REPLICAS`.

Collect every value from the running cluster itself — don't hard-code
anything from the lesson.

Get a shell: `docker exec -it breaklab-k8sd-004 bash`
