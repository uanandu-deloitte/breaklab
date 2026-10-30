# The Dashboard Panel That Queries Nothing

The "Checkout overview" dashboard is already provisioned in this box's
real Grafana (backed by a real Prometheus actually scraping a live
metrics endpoint) — open it up and its one panel, "HTTP requests",
renders completely empty. No error, no "no data" banner that points at
why, it's just blank.

Everything's actually running:

```
curl http://localhost:3000/api/health      # Grafana
curl http://localhost:9090/api/v1/targets  # Prometheus, scraping for real
```

Find the panel's real query (via the dashboard JSON under
`/etc/grafana/provisioning/dashboards/json/`, or Grafana's own API), and
figure out why it returns nothing even though Prometheus is actively
scraping real data. Fix the dashboard's query, not the metric name on
the export side.

Get a shell: `docker exec -it breaklab-obs-009 bash`
