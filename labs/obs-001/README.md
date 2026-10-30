# Logs vs Metrics

You have shell access to a host carrying two small observability files:

- `/var/log/app.jsonl` — a structured JSON-lines application log, one
  JSON object per line.
- `/var/lib/metrics/current.json` — a snapshot of current metric
  gauges.

## Tasks

1. Count how many lines in `/var/log/app.jsonl` have `"level": "ERROR"`.
2. Find the log line whose `request_id` is `req-3005` and record its
   `message`.
3. Read the current value of the `disk_used_pct` metric.
4. Decide, for each of these two questions, which file actually answers
   it:
   - "Why did this particular request fail?" — answered by the file
     that records discrete, timestamped **events** with detail: a LOG.
   - "How full is the disk right now?" — answered by the file that
     holds a current **numeric value**, no event detail: a METRIC.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line (quote any value containing spaces):

```
ERROR_COUNT=
REQUEST_ID_MESSAGE=
DISK_USED_PCT=
WHY_ANSWER=
HOWMUCH_ANSWER=
```

`WHY_ANSWER` and `HOWMUCH_ANSWER` should each be exactly `LOG` or
`METRIC`.

Get a shell: `docker exec -it breaklab-obs-001 bash`
