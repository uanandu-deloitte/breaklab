# Deployments and Rollouts

A `web` Deployment (2 replicas) is already running `nginx:1.25-alpine` and
healthy. Your job is to roll it forward to a new image tag and prove the
rollout actually finished, not just that you ran the command.

## Look around first

- `kubectl get deployment web -o jsonpath='{.spec.template.spec.containers[0].image}'`
  — confirm the current image tag.
- `kubectl get pods -l app=web -o wide` — two Pods, both Running.

## Roll out the new tag

Update the `web` container to **`nginx:1.27-alpine`**:

```
kubectl set image deployment/web web=nginx:1.27-alpine
```

Then watch the rollout to completion:

```
kubectl rollout status deployment/web
```

Don't move on until that command reports the rollout finished — a
Deployment's spec updates immediately, but the real Pods take a moment to
be replaced one at a time.

## Confirm it for yourself

Check that every live Pod for `web` is actually running the new tag, not
just that the Deployment's spec says so:

```
kubectl get pods -l app=web -o jsonpath='{range .items[*]}{.metadata.name}{" -> "}{.spec.containers[0].image}{"\n"}{end}'
```

Get a shell: `docker exec -it breaklab-k8sf-002 bash`
