# Solution: Dashboards Query by Name — Renames Break Silently

```bash
panel_metric=$(sed -n 's/^metric: *//p' /etc/dashboards/overview.yaml)
exported_metric=disk_free_percent

render-dashboard
rendered=$(cat /var/lib/dashboards/overview.rendered)

cat <<EOF > /workspace/answers.env
PANEL_METRIC=$panel_metric
EXPORTED_METRIC=$exported_metric
RENDER_RESULT=$rendered
EOF
```

`/etc/dashboards/overview.yaml` asks for `disk_free_pct`. `jq 'keys'`
against `/var/lib/metrics/exported.json` shows the exporter is actually
emitting `disk_free_percent` — the same underlying disk-space metric,
renamed from a `_pct` suffix to a spelled-out `_percent` one. The panel
config was never updated to match. `render-dashboard` does an exact
string lookup of `disk_free_pct` in the exporter's JSON, finds nothing,
and writes `no-data` — not an error, just silence, because a missing key
and a not-yet-reported metric look identical to the render step.
