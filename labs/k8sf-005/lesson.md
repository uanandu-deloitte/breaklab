# Readiness vs. liveness, in six commands

Kubernetes gates Service traffic on **readiness**, not on whether a
container process is merely alive. A Pod can be `Running` (its container
started and hasn't exited) while being `NotReady` (failing its readiness
probe) indefinitely — and the Service in front of it will never route
traffic there. Liveness is a completely separate mechanism: it only
controls whether Kubernetes restarts the container, and a Pod with no
liveness probe configured at all will simply stay Running-but-NotReady
forever if its readiness probe never passes.

## Why this matters

"The container didn't crash" and "the container is receiving traffic"
are two different claims, and conflating them is a common source of
confusing incidents: `kubectl get pods` can show `Running` with a
healthy-looking restart count of `0`, while the Service in front of it
silently has nowhere to send a single request. Nothing about this state
produces an error anywhere by default — it's a deliberate, quiet safety
mechanism, not a glitch.

## Where to look

- `kubectl get pods -o wide` shows STATUS and READY as two separate
  columns — `Running` and `0/1` is a completely normal, valid, and
  deliberate combination.
- `kubectl describe pod <name>` prints the readiness probe's own failure
  message in its Events section — connection refused, a non-2xx status
  code, or a timeout, naming exactly what's wrong.
- `kubectl get endpoints <service>` shows the Service's actual computed
  backend list — a Pod that's Running-but-NotReady simply isn't in it,
  with no accompanying error logged anywhere else.
- `kubectl get deployment <name> -o
  jsonpath='{.spec.template.spec.containers[0].readinessProbe}'` reads the
  probe's live configuration (path, port, timing) exactly as deployed.
- A Pod's `readinessProbe` is immutable once the Pod is created — fixing
  one always means changing the Deployment and letting it replace the
  Pod with a corrected one, never editing the live Pod directly.

The habit this lab builds: when traffic never reaches a Pod that "looks
fine," check readiness and Endpoints before suspecting the network, the
Service, or anything else.
