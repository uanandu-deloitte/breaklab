# Labels, Selectors, and How Services Find Pods

A `reports` Deployment (1 replica) is running, and its Pod is `Running`
and `Ready`. A `reports` Service also exists — but it has no backends to
send traffic to at all.

## Look around first

- `kubectl get pods --show-labels` — what labels does the `reports` Pod
  actually carry?
- `kubectl get svc reports -o jsonpath='{.spec.selector}'` — what is the
  Service actually trying to match?
- `kubectl describe svc reports` — selector and Endpoints, together.
- `kubectl get endpoints reports` — how many addresses does it list?
- `kubectl get pods -l app=reports-old` — try filtering Pods by the
  Service's own selector value directly. What comes back?

## Record your findings

Write these to `/workspace/answers.env`, one `KEY=value` per line:

```
SERVICE_SELECTOR_VALUE=
POD_LABEL_VALUE=
ENDPOINT_COUNT=
```

- `SERVICE_SELECTOR_VALUE`: the live value of the Service's
  `spec.selector.app`, read from `kubectl get svc reports -o
  jsonpath='{.spec.selector.app}'`.
- `POD_LABEL_VALUE`: the live value of the Pod's actual `app` label —
  not what the Service expects, what the Pod really has.
- `ENDPOINT_COUNT`: how many addresses `kubectl get endpoints reports`
  actually lists right now.

Get a shell: `docker exec -it breaklab-k8sf-006 bash`
