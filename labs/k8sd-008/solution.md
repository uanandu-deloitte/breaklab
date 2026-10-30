# Solution: Helm Rendered the Wrong Port

```bash
helm get manifest web | grep -A3 "kind: Service"
#   targetPort: 8080   <- the container never listens here

kubectl get pods -l app=web-web -o jsonpath='{.items[0].spec.containers[0].ports}'
# containerPort: 80    <- this is the real listening port
```

Fix `/chart/templates/service.yaml`:

```yaml
ports:
  - port: {{ .Values.service.port }}
    targetPort: {{ .Values.containerPort }}
```

```bash
helm upgrade web /chart
kubectl run check --image=busybox:1.36 --restart=Never --rm -i \
  --command -- wget -qO- http://web-web:8080
# Server address: ...
```
