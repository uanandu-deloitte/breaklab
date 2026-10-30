# Logs vs metrics: two different questions

Observability tooling gives you two fundamentally different kinds of
data, and knowing which one to reach for saves time when something's
wrong.

## A log is a discrete, timestamped event

Every line in `/var/log/app.jsonl` is one thing that happened, at one
instant: a request arrived, a query ran slow, a charge was declined.
Logs answer **why** — they carry enough detail (a message, a request id,
a level) to explain one specific occurrence.

## A metric is a current numeric value

`/var/lib/metrics/current.json` holds no events at all — just the latest
reading of a handful of gauges: `cpu_pct`, `disk_used_pct`,
`error_rate`. A metric answers **how much** — it tells you the system's
current state, but nothing about any one request or decision that led
there.

## Reading JSON with jq

- `jq 'select(.level=="ERROR")' /var/log/app.jsonl` filters a stream of
  JSON-lines objects down to just the ones matching a condition — jq
  reads consecutive JSON values from a file even when each one sits on
  its own line, so no special "jsonl mode" is needed.
- `jq -r '.some_key' file.json` pulls a single key's value straight out
  of a JSON object, with `-r` printing it unquoted.
- Piping either into `wc -l` turns a list of matching lines into a
  count.

## Why this matters

If you only have metrics, you can tell something is wrong (error rate
is up) but not which request failed or why. If you only have logs, you
can explain any one failure in detail but have to scan a lot of lines
to estimate how bad things are right now. Real systems carry both, and
knowing which one answers your current question keeps you from digging
in the wrong place.
