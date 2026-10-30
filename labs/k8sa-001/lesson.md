# A ResourceQuota is enforced by the API server, not by you

A `Namespace` on its own has no capacity limit — anyone with access to it
could, in principle, create Pods until the whole cluster runs out of
room. A `ResourceQuota` object closes that gap: it's attached to one
namespace and caps some total across every object of a kind in it — Pod
count, aggregate CPU/memory requests, number of Services, and more.

## What "hard" and "used" mean

```
kubectl get resourcequota pod-quota -n quota-demo
NAME         AGE   REQUEST     LIMIT
pod-quota    2m    pods: 2/2
```

`hard` is the ceiling you configured (`spec.hard`). `used` is a live
counter the API server maintains itself, incrementing it the instant an
object that counts against the quota is admitted, decrementing it the
instant one is deleted. Neither number is something you compute by hand
— `kubectl get resourcequota` always shows the real, current count.

## Why the rejection happens where it does

A `ResourceQuota` is enforced by an admission controller — code that runs
inside the API server itself, on every create request, before the object
is ever written to storage. That means:

- There's no asynchronous reconciliation involved. The moment a request
  would push a count over `hard`, the API server's response to that exact
  `kubectl apply`/`kubectl create` call is a rejection — not a warning
  later, not a Pod that gets created and then torn down.
- The error is specific. It names the quota that was violated and shows
  the exact numbers: what was requested, what's already used, what the
  limit is. That's real, structured information from the server, not a
  guess — which is why reading it (rather than assuming its wording) is
  the point of this lab.

## Why this matters operationally

Quotas are how a cluster operator keeps one team or namespace from
starving every other tenant on a shared cluster. A learner who has only
ever worked in an unconstrained sandbox namespace has never seen this
failure mode — and the failure mode is common: CI pipelines that spin up
scratch Pods, runaway Jobs, or just an application that scales more
replicas than a namespace was budgeted for. Recognizing "exceeded quota"
or "forbidden" in an error instantly tells you which layer blocked the
request — no application code, no scheduler, no node capacity problem,
just a namespace-level ceiling.
