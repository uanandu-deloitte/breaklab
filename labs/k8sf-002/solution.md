# Solution: Deployments and Rollouts

```bash
kubectl set image deployment/web web=nginx:1.27-alpine
kubectl rollout status deployment/web
```

Confirm every live Pod actually picked up the new tag:

```bash
kubectl get pods -l app=web -o jsonpath='{range .items[*]}{.metadata.name}{" -> "}{.spec.containers[0].image}{"\n"}{end}'
```

Both lines should end in `nginx:1.27-alpine`. `kubectl rollout status`
having returned 0 already guarantees this, but reading the Pods directly
is the independent proof — the same thing `check.sh` does, rather than
trusting the Deployment's spec alone.
