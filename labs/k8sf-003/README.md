# Services and In-Cluster DNS

An `api` Deployment (2 replicas) and its Service are already running and
healthy. Your job is to reach the Service from inside the cluster two
different ways, and record proof of both.

## Find the Service's ClusterIP

```
kubectl get svc api -o jsonpath='{.spec.clusterIP}'
```

## Launch a throwaway debug Pod

```
kubectl run probe --rm -it --image=busybox:1.36 --restart=Never -- sh
```

From inside that shell, reach the `api` Service **both** ways:

```
wget -qO- api.default.svc.cluster.local      # by DNS name (short form "api" also works)
wget -qO- <the ClusterIP you found above>    # by ClusterIP
```

Both should print `ok`. Exit the shell (`exit`) when you're done — `--rm`
cleans the Pod up automatically.

## Record your findings

Write these to `/workspace/answers.env`, one `KEY=value` per line:

```
SERVICE_CLUSTERIP=
DNS_LOOKUP_WORKED=
```

- `SERVICE_CLUSTERIP`: the exact ClusterIP you read from `kubectl get svc
  api`, not a guess.
- `DNS_LOOKUP_WORKED`: `yes` if the DNS-name request above returned `ok`,
  `no` otherwise.

Get a shell: `docker exec -it breaklab-k8sf-003 bash`
