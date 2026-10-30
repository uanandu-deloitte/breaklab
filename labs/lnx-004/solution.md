# Solution: Find the Production Log

```bash
grep -RlE "TX-8842" /var/log /opt /tmp 2>/dev/null
# -> /var/log/app/archive/2024-01/host-03/trace.log

echo "/var/log/app/archive/2024-01/host-03/trace.log" > /workspace/answer.txt
```

`grep -R` recurses through every file under the given directories;
`-l` prints just the matching file's path instead of the matching
line. The decoy logs under `/var/log/app`, `/opt/acme-app/logs`, and
`/tmp/old-logs` all have plausible names but none contain
`TX-8842` — only the buried, unassumingly-named
`.../archive/2024-01/host-03/trace.log` does.
