# Solution: HTTP 500 Investigation

```bash
cat <<'EOF' > /workspace/top-500-ip.sh
#!/usr/bin/env bash
awk '$5 == 500 {print $1}' /workspace/access.log | sort | uniq -c | sort -rn | head -1 | awk '{print $1, $2}'
EOF
chmod +x /workspace/top-500-ip.sh
```

`awk` pulls out just the IP column (`$1`) for lines whose status column
(`$5`) is 500. `sort | uniq -c` tallies how many times each IP appears;
`sort -rn` puts the highest count first; `head -1` takes just the winner;
the final `awk` reformats `uniq -c`'s `"  4 198.51.100.9"` into the
requested `"4 198.51.100.9"`.
