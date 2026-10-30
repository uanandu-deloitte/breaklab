# Solution: A Policy Engine That Isn't Enforcing Anything

```bash
kubectl get clusterpolicy disallow-privileged-containers -o yaml
# spec.validationFailureAction: Audit   <- reports only, never blocks

sed -i 's/validationFailureAction: Audit/validationFailureAction: Enforce/' /manifests/policy.yaml
kubectl apply -f /manifests/policy.yaml

kubectl delete pod bad-pod

kubectl apply -f /manifests/bad-pod.yaml
# Error from server: admission webhook "validate.kyverno.svc-fail" denied
# the request: ... Privileged containers are not allowed.
```
