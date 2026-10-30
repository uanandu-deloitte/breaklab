# Solution: HighErrorRate Has Never Once Fired

```yaml
groups:
  - name: checkout
    rules:
      - alert: HighErrorRate
        expr: error_rate > 5
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Checkout error rate is high"
```

```bash
promtool test rules /infra/rules_test.yml
#   SUCCESS
```

The rule was comparing `error_rate < 5` — true when things are *fine*,
never true during an actual breach. Flipping the operator to `> 5`
makes it fire exactly when `error_rate` crosses the real threshold.
