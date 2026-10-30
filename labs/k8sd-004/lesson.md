# A healthy Pod proves nothing about what it depends on

Kubernetes tracks a Pod's own health obsessively — liveness probes,
readiness probes, restart counts, phase. None of that machinery
watches the *other* services a Pod calls over the network. A frontend
that makes an HTTP call to a backend Service can be perfectly
`Running` and `1/1 Ready` while every single one of those calls fails,
because readiness only describes whether this Pod is ready to receive
traffic — not whether the things it calls are.

## Endpoints are the ground truth, not the Deployment

A Service doesn't send traffic anywhere by magic — it forwards to
whatever addresses appear in its matching `Endpoints` object, and
those addresses are populated only from Pods that are both selected by
the Service's label selector *and* currently passing their readiness
probe. `kubectl get endpoints <name>` showing an empty address list
means exactly one thing: nothing is currently ready to receive
traffic through that Service. It does not, by itself, tell you *why*.

## The same symptom, two different causes

An empty Endpoints list can come from either direction, and they look
identical from a caller's point of view:

- **Scaled to zero.** The Deployment's `spec.replicas` is 0 —
  deliberately, by a capacity change, a bad automation, or a human
  running `kubectl scale` by mistake. There are no Pods to be ready,
  because there are no Pods at all.
- **Crash-looping.** The Deployment's `spec.replicas` is still what it
  always was, and Pods exist, but every one of them keeps failing its
  readiness probe (or keeps restarting) before it can serve traffic.

`kubectl get deployment <name>` prints both `DESIRED` and `READY` side
by side specifically so you can tell these apart: `DESIRED` at 0 means
scaled down; `DESIRED` greater than 0 with `READY` at 0 means replicas
exist but aren't healthy. The fix is different in each case — scaling
back up versus finding and fixing why the replicas never turn ready —
so confusing one for the other sends you to the wrong remediation.

## Where to start looking

The caller is usually the one making noise (a broken feature, a
failed request), but the caller's own Pod is rarely where the actual
fault lives. A frontend that logs its own dependency checks on a loop
will name the specific dependency it considers unreachable — follow
that name to the dependency's own Service and Deployment, not back
into the caller's already-healthy Pod.
