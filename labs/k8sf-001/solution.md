# Solution: Your First Pod

Orientation commands (any order):

```bash
kubectl get pods -o wide
kubectl describe pod -l app=hello
kubectl logs -l app=hello
kubectl get deployments
kubectl get svc hello
```

Make the real change:

```bash
kubectl scale deployment/hello --replicas=4
kubectl rollout status deployment/hello
```

Record the findings, reading every value live rather than typing in a
guess:

```bash
image=$(kubectl get deployment hello -o jsonpath='{.spec.template.spec.containers[0].image}')
cat <<EOF > /workspace/answers.env
INITIAL_REPLICA_COUNT=2
CONTAINER_IMAGE=${image}
FINAL_REPLICA_COUNT=4
EOF
```

`INITIAL_REPLICA_COUNT` is 2 because that's what the Deployment's manifest
originally declared, before you scaled it — once you've scaled up, that
number is no longer visible live, so it has to come from what you observed
*before* making the change.
