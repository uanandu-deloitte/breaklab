# Rolling updates, in six commands

A Deployment doesn't edit Pods in place — changing its Pod template (most
commonly its container image) creates a new ReplicaSet with the new
template, and the Deployment controller scales that new ReplicaSet up
while scaling the old one down, one Pod at a time, until only the new
ReplicaSet remains.

## Why this matters

`kubectl set image` (or editing a manifest's image field and re-applying
it) returns as soon as the Deployment's *spec* is updated — that happens
instantly, before a single Pod has actually changed. Treating that command
returning as "the rollout is done" is one of the most common Kubernetes
mistakes: the old Pods, with the old image, can keep serving traffic for
seconds to minutes afterward while the rollout works through its update
strategy (by default, `RollingUpdate`: a bounded number of Pods replaced at
a time, with each new Pod required to pass its readiness probe before the
next old one is torn down).

## Where to look

- `kubectl set image deployment/<name> <container>=<image:tag>` changes
  the Pod template's image for one named container inside the Deployment.
  The container name (left of `=`) must match the one in the Pod spec —
  not the Deployment's own name.
- `kubectl rollout status deployment/<name>` is the authoritative way to
  know whether a rollout finished: it blocks and streams progress, then
  exits 0 once every new Pod is up and Ready, or reports a problem (e.g. a
  new Pod stuck failing its readiness probe) if the rollout can't
  complete.
- `kubectl get pods -l <selector> -o jsonpath=...` lets you read each live
  Pod's actual container image directly — the only way to be certain
  every Pod, not just the Deployment's spec, reflects the change.
- `kubectl rollout history deployment/<name>` and `kubectl rollout undo
  deployment/<name>` exist for when a rollout succeeds but the new image
  turns out to be bad — not needed for this lab, but worth knowing they're
  there.

The habit this lab builds: a rollout isn't "done" when you press enter on
the command that starts it — it's done when `kubectl rollout status`, or
an independent read of every live Pod, says so.
