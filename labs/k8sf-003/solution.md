# Solution: Services and In-Cluster DNS

```bash
ip=$(kubectl get svc api -o jsonpath='{.spec.clusterIP}')

kubectl run probe --rm -it --image=busybox:1.36 --restart=Never -- \
  sh -c "wget -qO- api.default.svc.cluster.local; echo; wget -qO- $ip"
```

Both `wget` calls should print `ok`. Then record the findings:

```bash
cat <<EOF > /workspace/answers.env
SERVICE_CLUSTERIP=${ip}
DNS_LOOKUP_WORKED=yes
EOF
```

`SERVICE_CLUSTERIP` comes from a live `kubectl get`, never typed in from
memory — a Service recreated from scratch would get a different ClusterIP,
which is exactly why real applications use the DNS name instead of a
hard-coded IP.
