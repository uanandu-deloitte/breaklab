# Solution: Ship Your Logs

```bash
ship-logs --source /var/log/app.log --destination /var/log/aggregator.log

echo "2026-10-20T00:01:00Z app new line shipped later" >> /var/log/app.log
ship-logs --source /var/log/app.log --destination /var/log/aggregator.log
```

The first run copies every existing line from `app.log` into
`aggregator.log`. `ship-logs` remembers how many lines it already
shipped from this source in a small state file under
`/var/lib/ship-logs/`, so the second run only appends the one new line
— not a duplicate of everything that shipped before.
