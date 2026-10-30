# Alert rules are just a boolean check

An alerting tool doesn't understand your system — it only evaluates a
tiny expression: a **metric** name, a **comparison operator**, and a
**threshold**, checked against whatever value a metrics file currently
reports. If no rule exists for a metric, nothing is ever compared, and
nothing ever fires, no matter how badly that metric is breaching.

## Why this matters

Plenty of real alerting silence isn't a broken pipeline — it's simply
that nobody ever wrote the rule. An evaluator running perfectly against
zero rules produces exactly the same visible symptom (silence) as one
that's genuinely broken, so before you debug a pipeline, check whether a
rule exists at all.

## Where to look

- `/var/lib/metrics/current.json` holds the live gauge values an
  evaluator reads from — a plain JSON object, one key per metric.
- `/etc/alerts/rules.yaml` is where a rule lives: three fields,
  `metric`, `op`, and `threshold`.
- `alert-evaluator` reads the rule file, looks up that metric's current
  value, applies the operator against the threshold, and appends a line
  to `/var/log/alerts.log` only when the comparison is true.

## Writing one

A rule that should fire once `error_rate` climbs past `5` looks like:

```yaml
metric: error_rate
op: ">"
threshold: 5
```

The operator is one of `>`, `>=`, `<`, or `<=` — pick the direction that
actually matches what "bad" looks like for that metric.
