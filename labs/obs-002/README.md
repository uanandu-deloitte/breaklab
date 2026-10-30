# Structured vs Unstructured Logs

Two files on this host describe the exact same sequence of events for
two checkout attempts — one logged as plain text, one as structured
JSON:

- `/var/log/legacy.log` — plain-text lines: `<date> <time> <LEVEL>
  <message>`.
- `/var/log/app.jsonl` — the same events, one JSON object per line
  (`timestamp`, `level`, `message`, `order_id`, and `attempt` where it
  applies).

Order `ord-5521`'s checkout retried its database connection a few times
before failing for good. Order `ord-5522` succeeded immediately.

## Tasks

For order `ord-5521`, answer the same two questions **twice** — once
using only `/var/log/legacy.log` (grep/awk/sed — whatever text tool
gets you there), and once using only `/var/log/app.jsonl` with a
single `jq` filter:

1. How many database-connection attempts did it retry before failing?
2. What was the exact timestamp of the final failure line
   (`checkout failed for order ord-5521: ...`)?

Record both sets of answers in `/workspace/answers.env`, one
`KEY=value` per line, quoting any value containing spaces:

```
ATTEMPTS_FROM_LEGACY=
FINAL_FAILURE_TIME_FROM_LEGACY=
METHOD_LEGACY=
ATTEMPTS_FROM_JSON=
FINAL_FAILURE_TIME_FROM_JSON=
METHOD_JSON=
```

`METHOD_LEGACY` and `METHOD_JSON` should each just name the tool you
used to get that answer (e.g. `grep`, `awk`, `jq`).

Get a shell: `docker exec -it breaklab-obs-002 bash`
