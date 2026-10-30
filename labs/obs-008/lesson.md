# `promtool test rules` is a real unit test for alerting rules

Prometheus alerting rules are just PromQL expressions evaluated on a
schedule. `promtool test rules` lets you unit-test one without running
a live Prometheus server or waiting for a real incident: you declare a
synthetic time series (`input_series`), tell it what should (or
shouldn't) be firing at a given point in simulated time
(`alert_rule_test`), and it evaluates the real rule file against that
series exactly the way a running Prometheus would.

## Why this matters

A silently-wrong alert rule is worse than no rule at all — it looks
like coverage exists, so nobody adds a second line of defense, and the
first real incident is the only time anyone finds out it never worked.
Unit-testing alerting rules in CI (exactly what `promtool test rules`
is for) catches this before it ships, the same way you'd unit-test any
other logic.

## Where to look

- `rules.yml` — the real rule: `expr`, `for`, `labels`, `annotations`.
- `rules_test.yml` — the real test: a synthetic `error_rate` series that
  breaches at `t=3m` and stays breached, with an assertion that
  `HighErrorRate` must be firing by `t=9m` (well past the rule's
  `for: 5m`).
- `promtool test rules rules_test.yml` — run it and read exactly what it
  says it expected versus what it got.

## The fix

Read the comparison operator in `expr` against what "high" actually
means for this metric, and point it the right direction.
