# Restarts, Backoff, and Which Pod You're Actually Looking At

A `flaky-worker` Deployment is running. Its single Pod crashed a couple
of times right after it was created, then stabilized — the live
container you'd see right now is not the first attempt.

## Look around first

- `kubectl get pods -l app=flaky-worker` — what does RESTARTS say?
- `kubectl describe pod -l app=flaky-worker` — what do the Events say
  about each restart, including the backoff delay before each one?
- `kubectl logs -l app=flaky-worker` — the current instance's own
  output (may say nothing about the earlier crash at all).
- `kubectl logs -l app=flaky-worker --previous` — the instance right
  before this one, including whatever it printed right before it exited.

## Record your findings

Write these to `/workspace/answers.env`, one `KEY=value` per line:

```
RESTART_COUNT=
PREVIOUS_EXIT_REASON=
```

- `RESTART_COUNT`: the Pod's live restart count, read from `kubectl get
  pod <name> -o
  jsonpath='{.status.containerStatuses[0].restartCount}'` — not guessed
  from the RESTARTS column alone.
- `PREVIOUS_EXIT_REASON`: the previous instance's termination reason,
  read from `kubectl get pod <name> -o
  jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}'`.

Get a shell: `docker exec -it breaklab-k8sf-004 bash`
