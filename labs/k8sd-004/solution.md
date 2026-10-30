# Solution: Tracing a Failure Across a Service Dependency

```bash
dep="backend"   # the name frontend's own logs report as unreachable

kubectl get endpoints "$dep"
kubectl get deployment "$dep"

endpoints=$(kubectl get endpoints "$dep" -o jsonpath='{.subsets[*].addresses[*].ip}' | wc -w)
desired=$(kubectl get deployment "$dep" -o jsonpath='{.spec.replicas}')
ready=$(kubectl get deployment "$dep" -o jsonpath='{.status.readyReplicas}')

if [ "$desired" = "0" ]; then mode=scaled-to-zero; else mode=crash-looping; fi

cat <<EOF > /workspace/answers.env
DEPENDENCY_NAME=$dep
READY_ENDPOINTS=$endpoints
DESIRED_REPLICAS=$desired
READY_REPLICAS=${ready:-0}
FAILURE_MODE=$mode
EOF
```

`kubectl logs deployment/frontend --tail=20` names `backend` as the
dependency it currently can't reach — that's where to look next,
not back into frontend's own (healthy) Pod. `kubectl get endpoints
backend` confirms the Service genuinely has zero ready addresses: the
ground truth for "can anything reach this right now," independent of
whatever the Deployment object claims.

`kubectl get deployment backend`'s `DESIRED` and `READY` columns are
what actually tells the two failure modes apart, since both produce
the identical empty Endpoints list:

- `DESIRED` is `0` → nobody is running on purpose. The fix is scaling
  back up, not debugging a crash that isn't happening.
- `DESIRED` is greater than `0` but `READY` is `0` → replicas exist
  and keep failing their readiness probe or restarting.
  `kubectl describe pod -l app=backend` would show the restart count
  and the Events explaining why, before reaching for a fix like
  `kubectl rollout undo`.

Nothing here changes the cluster's state — the task is purely
diagnostic: the right answer depends on reading `DESIRED` against
`READY` correctly, not on guessing from the symptom alone.
