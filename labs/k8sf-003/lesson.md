# Services and DNS, in six commands

A Kubernetes Service gets two independent ways to be reached from inside
the cluster: a stable **ClusterIP** (a virtual IP, routed by the cluster's
networking layer to whichever Pods the Service currently selects), and a
**DNS name** — `<service>.<namespace>.svc.cluster.local`, automatically
registered by the cluster's internal DNS the moment the Service exists.
Inside the same namespace, the short form `<service>` resolves too.

## Why this matters

Hard-coding a ClusterIP anywhere outside a quick manual test is a mistake —
ClusterIPs are stable for a Service's lifetime, but a Service recreated
from scratch gets a new one. DNS names are what real applications use for
exactly this reason: `api.default.svc.cluster.local` means "whatever
Service is currently called `api` in namespace `default`," independent of
any specific IP.

## Where to look

- `kubectl get svc <name> -o jsonpath='{.spec.clusterIP}'` reads the live
  ClusterIP directly, instead of assuming or remembering it.
- `kubectl run <name> --rm -it --image=<image> --restart=Never -- <cmd>`
  launches a one-off Pod, attaches your terminal to it, runs `<cmd>`
  interactively, and deletes the Pod the moment you exit — the standard
  way to poke at a cluster from the inside without leaving anything behind.
- A Service's DNS name follows a fixed pattern:
  `<service>.<namespace>.svc.cluster.local` from any namespace, or just
  `<service>` from within the same namespace.
- `wget -qO-` (quiet, output to stdout) is a convenient way to prove
  reachability end to end from a minimal image like `busybox`, without
  needing `curl` installed.

The habit this lab builds: reaching something by its DNS name and reaching
it by its current IP are two different claims — proving one doesn't prove
the other, and a thorough check proves both independently.
