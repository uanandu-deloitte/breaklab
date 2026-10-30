# Solution: Write an Alert Rule

```bash
cat <<EOF > /etc/alerts/rules.yaml
metric: error_rate
op: ">"
threshold: 5
EOF
alert-evaluator
cat /var/log/alerts.log
```

`alert-evaluator` reads the three fields back out of `rules.yaml`, looks
up `error_rate` in `/var/lib/metrics/current.json` (currently `7`), and
since `7 > 5` is true, appends a firing alert line to
`/var/log/alerts.log`.
