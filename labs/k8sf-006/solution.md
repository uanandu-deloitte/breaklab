# Solution: Labels, Selectors, and How Services Find Pods

```bash
kubectl get pods --show-labels
kubectl get svc reports -o jsonpath='{.spec.selector}'
kubectl describe svc reports
kubectl get endpoints reports
```

Record the findings, reading every value live rather than typing in a
guess:

```bash
selector=$(kubectl get svc reports -o jsonpath='{.spec.selector.app}')
pod_label=$(kubectl get pods -l tier=backend -o jsonpath='{.items[0].metadata.labels.app}')
addrs=$(kubectl get endpoints reports -o jsonpath='{.subsets[0].addresses[*].ip}' 2>/dev/null)
count=0
[ -n "$addrs" ] && count=$(echo "$addrs" | wc -w | tr -d ' ')

cat <<EOF > /workspace/answers.env
SERVICE_SELECTOR_VALUE=${selector}
POD_LABEL_VALUE=${pod_label}
ENDPOINT_COUNT=${count}
EOF
```

`SERVICE_SELECTOR_VALUE` (`reports-old`) and `POD_LABEL_VALUE`
(`reports`) differ — that mismatch alone is why `ENDPOINT_COUNT` is 0,
even though the Pod itself is completely healthy and Ready.
