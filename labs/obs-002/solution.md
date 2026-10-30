# Solution: Structured vs Unstructured Logs

```bash
# From the unstructured log (fragile: depends on exact wording and column position)
grep -c "ERROR could not connect to db" /var/log/legacy.log
# -> 3
grep "checkout failed for order ord-5521" /var/log/legacy.log | awk '{print $1, $2}'
# -> 2026-01-15 10:22:07

# From the structured log (reliable: a single jq filter keyed on real fields)
jq -r 'select(.order_id=="ord-5521" and .level=="ERROR" and (.message|startswith("could not connect to db")))|.timestamp' /var/log/app.jsonl | wc -l
# -> 3
jq -r 'select(.order_id=="ord-5521" and (.message|startswith("checkout failed")))|.timestamp' /var/log/app.jsonl
# -> 2026-01-15 10:22:07

cat <<EOF > /workspace/answers.env
ATTEMPTS_FROM_LEGACY=3
FINAL_FAILURE_TIME_FROM_LEGACY="2026-01-15 10:22:07"
METHOD_LEGACY=grep+awk
ATTEMPTS_FROM_JSON=3
FINAL_FAILURE_TIME_FROM_JSON="2026-01-15 10:22:07"
METHOD_JSON=jq
EOF
```

Both answers come out the same, because both files describe the same
events. The text version works, but only because the message wording
and column order happened to cooperate — rename the "ERROR" level or
reorder the fields and the grep/awk commands quietly break. The jq
filter keys off actual field names (`order_id`, `level`, `message`), so
it survives changes that would silently break a text-pattern search.
