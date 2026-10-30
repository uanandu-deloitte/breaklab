# Structured vs unstructured: the same facts, two levels of risk

A log's *content* can be identical whether it's written as plain text
or as structured JSON — what changes is how safely you can get a fact
back out of it.

## Unstructured text: fragile extraction

`/var/log/legacy.log` is a classic plain-text log: `<date> <time>
<LEVEL> <message>`. To pull a fact out of it you lean on the message's
exact wording and the columns' exact position — `grep` for a phrase,
`awk` to grab the first two space-separated fields as a timestamp. It
works, right up until someone rewords a message, adds a field, or
reorders the columns — then the same command silently returns nothing,
or worse, the wrong thing, with no error to tell you it broke.

## Structured JSON: reliable extraction

`/var/log/app.jsonl` describes the exact same events, but as one JSON
object per line with named fields: `timestamp`, `level`, `message`,
`order_id`, `attempt`. A single `jq 'select(.order_id=="ord-5521" and
.level=="ERROR")'` filter keys off those real field names instead of
guessing at text position or wording. Rename a message, add a field,
reorder the keys — the filter still works, because it never depended on
any of that in the first place.

## Why this matters

Both representations can answer the same question today. The
difference shows up later, when the message wording changes for an
unrelated reason and every grep-based dashboard or alert quietly goes
blind — while every jq-based one keeps working, because it was never
reading the sentence, just the fields.
