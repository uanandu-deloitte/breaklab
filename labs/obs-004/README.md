# Ship Your Logs

`/var/log/app.log` already has lines nobody's aggregator has seen.
`/var/log/aggregator.log` exists but is still empty. A shipper script at
`/usr/local/bin/ship-logs` can move the lines over, but it has never
been run.

Run it with its two required flags:

```
ship-logs --source /var/log/app.log --destination /var/log/aggregator.log
```

Confirm every line in `app.log` now appears in `aggregator.log`.

Then append one new line to `/var/log/app.log` and run `ship-logs` again
with the same two flags. Confirm only the new line lands in
`aggregator.log` — not a second copy of everything that already
shipped.

Get a shell: `docker exec -it breaklab-obs-004 bash`
