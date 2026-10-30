# Labels and selectors, in six commands

A Kubernetes Service never targets a Pod by name. It matches
`spec.selector` against each Pod's `metadata.labels`, and the result —
every currently-Ready Pod whose labels satisfy that selector — becomes
the Service's Endpoints list. Nothing else about a Pod (its name, which
Deployment created it, how it was scheduled) matters to a Service at
all.

## Why this matters

Because matching is pure label comparison, a mismatch produces zero
Endpoints with no error logged anywhere — not on the Service, not on the
Pod, not in cluster events. The Pod can be perfectly `Running` and
`Ready`; the Service's own status shows nothing wrong either. The only
observable symptom is an empty Endpoints list, which is why checking
Endpoints has to be a reflex, not a last resort.

## Where to look

- `kubectl get pods --show-labels` prints every Pod's actual, current
  labels — the ground truth a selector either matches or doesn't.
- `kubectl get svc <name> -o jsonpath='{.spec.selector}'` reads exactly
  what a Service is trying to match, live, rather than trusting a
  manifest file that may no longer reflect the cluster.
- `kubectl describe svc <name>` shows the selector and the Endpoints it
  currently produces together, in one place.
- `kubectl get endpoints <name>` is the direct, observable symptom: an
  empty ADDRESSES column means no Pod's labels currently satisfy the
  selector, independent of whether that Pod is otherwise healthy.
- `kubectl get pods -l <key>=<value>` run with the Service's own selector
  value is a quick, direct proof: if it returns nothing, that's exactly
  what the Service itself is seeing too.

The habit this lab builds: a Service's health and a Pod's health are two
separate claims. Checking Endpoints — not just the Pod, and not just the
Service's own status — is the only way to see what's actually routable
right now.
