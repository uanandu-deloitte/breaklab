# Pods, Deployments, and Services, in six commands

A **Pod** is the smallest thing Kubernetes schedules and runs — one or more
containers sharing a network namespace. A **Deployment** manages a set of
identical Pods on your behalf: it owns a desired replica count, and
continuously reconciles the real Pods toward that number. A **Service** is
a stable network identity in front of a changing set of Pods, matched by
label selector.

## Why this matters

Almost every Kubernetes lab after this one starts from the same question:
"what does the cluster actually say is happening right now?" Guessing, or
trusting what a manifest file claims, isn't enough — manifests describe
*desired* state; `kubectl get`/`describe`/`logs` show *live* state, and the
two can disagree for all sorts of reasons you'll spend the rest of this
stage diagnosing.

## Where to look

- `kubectl get pods` shows STATUS (Pending/Running/Succeeded/Failed) and a
  READY column (`n/m` containers passing their readiness probe) —
  these are two different questions.
- `kubectl describe pod <name>` prints a Conditions section and an Events
  log at the bottom — the Events log is usually the fastest way to see
  *why* something is or isn't happening.
- `kubectl logs <pod>` shows a container's stdout/stderr — the same output
  you'd see if you ran it in a foreground terminal.
- `kubectl get deployment <name>` shows DESIRED, CURRENT, UP-TO-DATE, and
  AVAILABLE as separate columns. A Deployment can have the desired number
  of Pods that exist (CURRENT) without all of them being ready yet
  (AVAILABLE) — especially right after a scale or rollout.
- `kubectl scale deployment/<name> --replicas=<n>` edits the Deployment's
  desired replica count directly. It's a live change to cluster state, not
  a change to any file — the Deployment's controller loop does the rest.
- `kubectl get svc <name>` shows a ClusterIP (a stable virtual IP inside
  the cluster) and the port(s) it listens on, independent of which Pods are
  currently backing it.

None of this is about fixing anything yet — it's building the habit of
checking live state with `kubectl` before trusting any file or any
assumption about what "should" be true.
