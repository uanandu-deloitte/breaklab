# The Cluster That Stopped Listening to Git

This cluster runs a real ArgoCD (just the controller/repo-server, no
UI) watching a real git repo for the `web` app's manifests. Someone
reorganized the repo's directory layout last week. Ever since, the
`web` Application has sat there refusing to reconcile — not failed
outright, just permanently out of sync with whatever's actually
committed.

## Tasks

1. Look at the real Application resource and its real status:

   ```
   kubectl get application web -n argocd
   kubectl get application web -n argocd -o yaml
   ```

2. Find the real repo content it's supposed to be tracking — it's
   served by a real git daemon in this container, under `/repos/app.git`.

3. Reconcile the Application's source so it actually points at where
   the manifests really live in the repo, and confirm it reaches
   `Synced` / `Healthy` for real — not just no errors, the actual `web`
   Deployment running in the `default` namespace.

Get a shell: `docker exec -it breaklab-k8sd-007 bash`
