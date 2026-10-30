# Write an Alert Rule

`/var/lib/metrics/current.json` is reporting live gauge values, but
`/etc/alerts/rules.yaml` doesn't define a rule yet — nothing is watching
these metrics at all.

Author a rule that alerts when `error_rate` exceeds `5`.

Create `/etc/alerts/rules.yaml` with exactly these three fields:

```yaml
metric: error_rate
op: ">"
threshold: 5
```

Then run `alert-evaluator` and confirm it actually fires: check
`/var/log/alerts.log` for a new line reporting the breach.

Get a shell: `docker exec -it breaklab-obs-003 bash`
