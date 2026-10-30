# Taints repel; tolerations are the one exception you grant

A taint is applied to a *node*, not a Pod, and it works in the opposite
direction from most scheduling concepts: instead of a Pod declaring
where it wants to go (`nodeSelector`, affinity), a taint lets a node
declare who it's willing to accept. By default, a tainted node accepts
nothing — every Pod is repelled from it unless that Pod explicitly
tolerates the exact taint in question.

## The three fields that must match

```
kubectl taint nodes <node> dedicated=special:NoSchedule
```

breaks down into three parts a toleration must match *all* of:

- **key**: `dedicated`
- **value**: `special`
- **effect**: `NoSchedule`

A toleration with the right key but wrong value (or no value, under
`Exists` instead of `Equal`) doesn't count. Neither does matching
key/value with the wrong effect. All three have to line up exactly:

```yaml
tolerations:
  - key: dedicated
    operator: Equal
    value: special
    effect: NoSchedule
```

## Why this surfaces as Pending with no logs

Taints are enforced by the scheduler itself, before a Pod is ever placed
on a node — so an untolerated Pod never gets a `nodeName` assigned, never
has a container start, and never produces logs. The only place the
reason is recorded is the Pod's own Events, visible through
`kubectl describe pod`, with a line like:

```
Warning  FailedScheduling  ...  0/1 nodes are available: 1 node(s) had untolerated taint {dedicated: special}.
```

This is the same category of problem as an impossible resource request
— the Pod simply doesn't fit anywhere the scheduler is allowed to place
it — but the *reason* is a policy decision (this node is reserved) rather
than a capacity shortfall.

## Why tolerations can't be patched in

Like a Pod's resource requests, `tolerations` is effectively fixed once
a Pod is created — there's no live `kubectl patch` that adds a missing
toleration to a Pod that already exists. The only way to change it is
the same pattern used for an impossible resource request: delete the
Pod, fix the manifest, and let the scheduler evaluate it fresh.

## Why this matters operationally

Taints are how cluster operators carve out dedicated capacity — GPU
nodes reserved for ML workloads, a pool reserved for one tenant, nodes
being drained for maintenance. A Pod landing on the wrong node because
nothing stopped it is a correctness problem; a Pod that can't land
*anywhere* because its author forgot the matching toleration is the
mirror-image mistake, and it's one every operator eventually has to
debug from the Pod side, not the node side.
