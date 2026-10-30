# A Service's targetPort is independent of the container's real port

Nothing checks that a Service's `targetPort` actually matches where the
pod's container listens — Kubernetes creates the Endpoint entry
(`pod-ip:targetPort`) exactly as declared, whether or not anything is
listening there. A Deployment can be perfectly `Ready` and a Service
can have `Endpoints` populated, and traffic through that Service can
still go nowhere, because `targetPort` and the container's actual
listening port were never the same number.

## Why this matters

In a hand-templated Helm chart, it's extremely easy for
`service.targetPort` to drift from the container's real port: someone
changes the app's listening port, updates the Deployment template, and
never touches the Service template because nothing forces them to —
the two values are only connected if the chart author connects them.

## Where to look

- `helm get manifest web` — the real rendered YAML Helm actually
  applied, not the source templates.
- `kubectl get pods -l app=web-web -o jsonpath='{.items[0].spec.containers[0].ports}'`
  — the real container port, straight from the running Pod spec.

## The fix

Template `targetPort` from the same value the Deployment uses for
`containerPort`, so the two can never drift apart again:

```yaml
targetPort: {{ .Values.containerPort }}
```
