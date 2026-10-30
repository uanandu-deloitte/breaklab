# An Application's source is just a repo URL, a revision, and a path

ArgoCD doesn't understand your repo's structure — `spec.source.path`
is a literal string it hands to the repo server on every reconcile. If
a repo reorganization moves manifests to a new directory and nobody
updates the Application, ArgoCD doesn't guess or fall back — it reports
a `ComparisonError` and leaves the cluster exactly where it last
successfully synced, forever, until the path is fixed.

## Why this matters

This is the single most common real GitOps failure mode: the
Application resource and the repo it watches drift out of agreement
after a refactor, and the cluster quietly stops tracking git at all —
not because anything crashed, but because the one string connecting
them went stale.

## Where to look

- `kubectl get application web -n argocd -o yaml` — read
  `status.conditions` for the real `ComparisonError` message, and
  `spec.source.path` for what it's currently pointed at.
- The repo itself, served at `/repos/app.git` by this container's real
  git daemon — clone it or just check what directories the latest
  commit on `main` actually has.

## The fix

```bash
kubectl patch application web -n argocd --type merge \
  -p '{"spec":{"source":{"path":"manifests/web"}}}'
```
