# Solution: Normalize Deployment Data

```bash
tr ',' '\n' < /workspace/hosts.txt \
  | tr 'A-Z' 'a-z' \
  | sed 's/^[ \t]*//; s/[ \t]*$//' \
  | sed '/^$/d' \
  | sort -u \
  > /workspace/normalized.txt
```

`tr ',' '\n'` splits any comma-separated entries onto their own lines, so
every hostname ends up one per line regardless of how it was originally
delimited. `tr 'A-Z' 'a-z'` lowercases everything. The `sed` strips
leading/trailing whitespace (spaces or tabs) from each line, and the
second `sed` drops any now-empty lines left behind by a trailing comma.
`sort -u` sorts the result alphabetically and removes duplicates in one
step.
