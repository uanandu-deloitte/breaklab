# Solution: Restarts, Backoff, and Which Pod You're Actually Looking At

```bash
pod=$(kubectl get pods -l app=flaky-worker -o jsonpath='{.items[0].metadata.name}')

kubectl get pod "$pod"
kubectl describe pod "$pod"
kubectl logs "$pod" --previous
```

Record the findings, reading every value live rather than typing in a
guess:

```bash
restarts=$(kubectl get pod "$pod" -o jsonpath='{.status.containerStatuses[0].restartCount}')
reason=$(kubectl get pod "$pod" -o jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}')

cat <<EOF > /workspace/answers.env
RESTART_COUNT=${restarts}
PREVIOUS_EXIT_REASON=${reason}
EOF
```

`RESTART_COUNT` and `PREVIOUS_EXIT_REASON` both come from a live
`kubectl get`, not from counting lines in `kubectl describe`'s Events
section — Events can scroll off after enough restarts, but the API's own
`restartCount` and `lastState.terminated.reason` fields are exact and
always current.
