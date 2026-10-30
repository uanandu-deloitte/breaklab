# LimitRange: defaults that get injected, not just documented

A `ResourceQuota` caps totals across a namespace. A `LimitRange` is the
companion object that handles the per-object case: what happens when a
single Pod shows up with no resource requests or limits of its own at
all?

Without a LimitRange, the answer is "nothing" — a container with no
`resources` block runs with unbounded limits and a request of zero,
which the scheduler effectively treats as "doesn't need anything" even
though it will happily consume whatever's available on its node. That's
exactly the gap a LimitRange closes.

## What `defaultRequest` and `default` mean

```yaml
spec:
  limits:
    - type: Container
      defaultRequest:
        cpu: 100m
        memory: 64Mi
      default:
        cpu: 200m
        memory: 128Mi
```

- `defaultRequest` is what gets filled into `resources.requests` on any
  container in the namespace that doesn't specify its own.
- `default` is what gets filled into `resources.limits` the same way.

Both are injected by an admission controller at the moment a Pod is
created — the same point in the request pipeline where a ResourceQuota
gets enforced, just doing the opposite job: instead of rejecting a
request, it's filling in blanks the Pod didn't set.

## Why this is admission-time, not a controller loop

There's no background process scanning the namespace, finding Pods with
missing resources, and patching them. If there were, you'd expect
`pod-before` — created before the LimitRange existed — to eventually get
patched too. It never will. The LimitRange only ever looks at objects as
they're being created; anything that already existed keeps whatever it
had, forever, unless something deletes and recreates it.

That's the real test of whether defaults were actually injected versus
just read from a manifest: a Pod's own YAML with no `resources` block,
applied *after* the LimitRange exists, coming back with real values when
you read its live spec. If you'd typed those numbers into the Pod's own
manifest, this lab would prove nothing about the LimitRange at all —
which is exactly why the task insists the new Pod's manifest stays empty
of a `resources` block.

## Why this matters operationally

Unbounded containers are how one noisy neighbor starves everything else
on a node — no request means the scheduler doesn't reserve anything for
it, and no limit means nothing stops it from consuming as much CPU or
memory as the node has free. A LimitRange is a safety net: teams that
forget to size their containers (which is most teams, most of the time)
still get *something* reasonable rather than *nothing at all*.
