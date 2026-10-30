# A dashboard panel is just a stored query, nothing more

A Grafana panel doesn't know anything about the metric it's supposed to
show — it stores a literal query string (`targets[].expr` for a
Prometheus panel) and re-runs it against the datasource every time the
panel renders. If the metric that query names was ever renamed on the
exporter side, the panel doesn't break loudly — it just queries a name
that no longer exists and renders empty.

## Why this matters

Metric renames happen constantly during real refactors (cleaning up
naming conventions, migrating instrumentation libraries, merging two
near-duplicate metrics). Nothing about a rename on the exporter side
automatically updates every dashboard that queried the old name — each
one silently goes dark, one panel at a time, and the only signal is an
empty graph nobody's necessarily looking at that day.

## Where to look

- `/etc/grafana/provisioning/dashboards/json/checkout.json` — the real
  provisioned dashboard; read `panels[].targets[].expr`.
- `curl http://localhost:9090/api/v1/label/__name__/values` — lists
  every real metric name Prometheus currently has, straight from the
  source of truth.
- `curl --data-urlencode 'query=<expr>' http://localhost:9090/api/v1/query`
  — run the panel's exact query directly against Prometheus to see what
  it returns, independent of Grafana.

## The fix

Update the panel's `expr` in the dashboard JSON to the metric's current
real name, then confirm the same query returns real data straight from
Prometheus.
