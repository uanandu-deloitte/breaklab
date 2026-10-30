# Solution: Read the Access Log

```bash
cat <<EOF > /workspace/answers.env
TOTAL_REQUESTS=$(wc -l < /workspace/access.log)
FIRST_TIMESTAMP=$(head -1 /workspace/access.log | cut -d' ' -f2)
LAST_TIMESTAMP=$(tail -1 /workspace/access.log | cut -d' ' -f2)
COUNT_2XX=$(cut -d' ' -f5 /workspace/access.log | grep -c '^2')
COUNT_4XX=$(cut -d' ' -f5 /workspace/access.log | grep -c '^4')
COUNT_5XX=$(cut -d' ' -f5 /workspace/access.log | grep -c '^5')
MOST_FREQUENT_PATH=$(cut -d' ' -f4 /workspace/access.log | sort | uniq -c | sort -rn | head -1 | awk '{print $2}')
EOF
```

`wc -l` counts every line. `head -1`/`tail -1` grab the first and last
lines, and `cut -d' ' -f2` pulls out just the timestamp column. `cut -d'
' -f5` isolates the status column; `grep -c '^2'` (and `'^4'`, `'^5'`)
counts how many of those start with the given digit. `cut -d' ' -f4 |
sort | uniq -c | sort -rn | head -1` ranks paths by frequency and takes
the winner, and the final `awk '{print $2}'` strips `uniq -c`'s leading
count to leave just the path.
