# Solution: The Dashboard Panel That Queries Nothing

```bash
curl -s http://localhost:9090/api/v1/label/__name__/values | jq .
# ["http_requests_count_total", "up", ...]
```

The panel's `expr` is `http_requests_total` — that metric was renamed
to `http_requests_count_total` and nothing updated this dashboard.

Edit `/etc/grafana/provisioning/dashboards/json/checkout.json`'s panel
query:

```json
"targets": [
  {
    "expr": "http_requests_count_total",
    "refId": "A",
    "datasource": { "type": "prometheus", "uid": "prometheus" }
  }
]
```

Grafana's file-based dashboard provisioning picks up the change
automatically (it polls the provisioning directory), or restart Grafana
to force an immediate reload. Querying the same `expr` directly against
Prometheus now returns real data.
