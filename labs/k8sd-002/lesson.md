# How Kubernetes derives a Pod's QoS class

Every Pod gets one of three QoS (Quality of Service) classes, and
Kubernetes computes it — nobody sets it directly — purely from the
requests and limits on its containers:

- **BestEffort**: no container sets any CPU or memory request or
  limit, for anything. Lowest priority — the first killed when a node
  runs low on memory.
- **Burstable**: at least one container sets a request or limit, but
  it doesn't qualify for Guaranteed (see below). Most real-world Pods
  land here.
- **Guaranteed**: every container sets both a CPU request and a CPU
  limit, and both a memory request and a memory limit — and for every
  container, request equals limit on both resources. This is the
  highest priority class; these Pods are the last killed under memory
  pressure and get a fixed, dedicated CPU share rather than a flexible
  one.

The requirement is strict: one container with `requests.cpu` slightly
different from `limits.cpu` anywhere in the Pod drops the whole Pod
out of Guaranteed and into Burstable, even if every other container
is set up correctly.

## Why this matters beyond a label

QoS class isn't cosmetic — it's what the kubelet actually uses to
decide eviction order under node memory pressure, and (for Guaranteed
Pods specifically) it can pin a container to dedicated CPU cores if
the node's `--cpu-manager-policy=static` is set. A Pod that's supposed
to behave predictably under load — a budget-sensitive service, a
database, anything that can't tolerate getting evicted first — needs
to actually earn Guaranteed, not just have requests/limits of some
kind.

## Why the limit still has to be realistic

Setting request == limit gets you the QoS class, but it says nothing
about whether the number itself is enough for the container to run. A
Guaranteed Pod with a memory limit too low for its own process still
gets OOM killed — it just gets OOM killed with a fancier QoS label.
Confirming the Pod is actually Running, not just that qosClass reads
Guaranteed, is part of the job, not an afterthought.
