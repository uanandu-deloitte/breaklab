# Solution: The Cluster That Stopped Listening to Git

```bash
kubectl get application web -n argocd -o yaml
# status.conditions: ComparisonError -- "manifests/app: app path does not exist"

kubectl patch application web -n argocd --type merge \
  -p '{"spec":{"source":{"path":"manifests/web"}}}'

kubectl get application web -n argocd
# SYNC STATUS   HEALTH STATUS
# Synced        Healthy

kubectl get deployment web -n default
# web   1/1   1   1
```
