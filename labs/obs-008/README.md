# HighErrorRate Has Never Once Fired

`/infra/rules.yml` declares a real Prometheus alerting rule,
`HighErrorRate`, meant to catch checkout's error rate breaching 5%.
Checkout has had real incidents where error rate sat at 8% for ten
minutes straight. `HighErrorRate` never fired for any of them.

`/infra/rules_test.yml` is a real `promtool` unit test for this rule —
synthetic input series, no live Prometheus server needed. Run it:

```
promtool test rules /infra/rules_test.yml
```

Read exactly what it says is wrong, fix the real rule in `rules.yml`,
and get the test to pass.

Get a shell: `docker exec -it breaklab-obs-008 bash`
