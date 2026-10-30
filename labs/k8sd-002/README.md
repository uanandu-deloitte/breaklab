# Requests, Limits, and QoS

The `qos-demo` Pod currently has no CPU or memory requests or limits
at all, which makes it BestEffort — the first Pod evicted under node
memory pressure, with no CPU share guaranteed either.

## Tasks

1. `kubectl get pod qos-demo -o jsonpath='{.status.qosClass}'` —
   confirm it currently reports `BestEffort`.
2. `kubectl get pod qos-demo -o jsonpath='{.spec.containers[0].resources}'`
   — confirm there's nothing there to derive a class from.
3. Fix the Pod so it ends up **Guaranteed**: every container's CPU
   and memory requests must exactly equal its limits. (A Pod's
   resources are immutable — delete and recreate it with
   `kubectl delete pod qos-demo` then `kubectl apply`, rather than
   trying to patch the running one.)
4. Pick a limit generous enough that the container doesn't get OOM
   killed or throttled into failure — confirm with
   `kubectl get pod qos-demo` that it's Running, not CrashLoopBackOff.
5. Confirm: `kubectl get pod qos-demo -o jsonpath='{.status.qosClass}'`
   now reports `Guaranteed`.

Get a shell: `docker exec -it breaklab-k8sd-002 bash`
