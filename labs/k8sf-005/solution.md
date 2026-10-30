# Solution: Readiness Probes Gate Traffic, Not Liveness

```bash
kubectl get pods -l app=catalog-api
kubectl describe pod -l app=catalog-api
kubectl get endpoints catalog-api
kubectl get deployment catalog-api -o jsonpath='{.spec.template.spec.containers[0].readinessProbe}'
```

Record the findings, reading every value live rather than typing in a
guess:

```bash
ready_column=$(kubectl get pods -l app=catalog-api --no-headers | awk '{print $2}')
addrs=$(kubectl get endpoints catalog-api -o jsonpath='{.subsets[0].addresses[*].ip}' 2>/dev/null)
count=0
[ -n "$addrs" ] && count=$(echo "$addrs" | wc -w | tr -d ' ')
probe_path=$(kubectl get deployment catalog-api -o jsonpath='{.spec.template.spec.containers[0].readinessProbe.httpGet.path}')

cat <<EOF > /workspace/answers.env
READY_COLUMN=${ready_column}
ENDPOINT_COUNT=${count}
PROBE_PATH=${probe_path}
EOF
```

`ENDPOINT_COUNT` is 0 because the container only serves `/index.html`,
while the readiness probe checks `/health` — a 404 the probe treats as a
failure, forever, with no effect on the container's restart count at
all, since no liveness probe is configured.
