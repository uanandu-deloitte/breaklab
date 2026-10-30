# Restarts, backoff, and which instance you're looking at

When a container inside a Pod exits, Kubernetes restarts it according to
the Pod's `restartPolicy` (`Always` by default). Each restart replaces the
*container*, not the Pod — the Pod keeps its name and identity, but the
process inside it may already be on its second, third, or tenth attempt
by the time you look.

## Why this matters

`kubectl get pods` isn't just STATUS and READY — it also has a RESTARTS
column, and that number is your first clue that what you're seeing live
is not the Pod's original, first-ever execution. Checking `kubectl logs
<pod>` only shows the *current* instance's output; if it just restarted a
moment ago, that can be empty or unrelated to whatever actually went
wrong. The actual failure lives in the *previous* instance's own output
and exit status, not the current one's.

## Where to look

- `kubectl get pods` shows a RESTARTS count directly; `kubectl get pod
  <name> -o jsonpath='{.status.containerStatuses[0].restartCount}'` reads
  the same number as a precise integer.
- `kubectl describe pod <name>` prints an Events section listing each
  start and crash in order, including the backoff delay Kubernetes waited
  before each retry (crash backoff grows: roughly 10s, 20s, 40s, up to a
  cap).
- `kubectl logs <pod> --previous` shows the stdout/stderr of the instance
  immediately before the current one — this is where an actual crash
  message is, when the live instance's own logs show nothing useful.
- `kubectl get pod <name> -o
  jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}'`
  (and `...lastState.terminated.exitCode`) report exactly why and how the
  previous instance ended, straight from the API rather than inferred
  from logs alone.

The habit this lab builds: before trusting anything about "what
happened," check whether you're looking at the first attempt or the Nth —
and if it's the Nth, go find the (N-1)th's own output before concluding
anything.
