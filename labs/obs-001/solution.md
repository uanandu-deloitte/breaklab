# Solution: Logs vs Metrics

```bash
jq 'select(.level=="ERROR")' /var/log/app.jsonl | jq -s 'length'
# -> 3

jq -r 'select(.request_id=="req-3005") | .message' /var/log/app.jsonl
# -> upstream inventory service timed out

jq -r '.disk_used_pct' /var/lib/metrics/current.json
# -> 78

cat <<EOF > /workspace/answers.env
ERROR_COUNT=3
REQUEST_ID_MESSAGE="upstream inventory service timed out"
DISK_USED_PCT=78
WHY_ANSWER=LOG
HOWMUCH_ANSWER=METRIC
EOF
```

A log line is a discrete, timestamped event — it can tell you *why*
something happened (the exact message tied to one request id). A
metric is a single current numeric value — it can tell you *how much*
is happening right now, but carries no detail about any one event.
