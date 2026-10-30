# Solution: Inspect Configuration

```bash
diff /workspace/config.old /workspace/config.new
```

```
1c1
< database_host: db-old.internal
---
> database_host: db-new.internal
2c2
< log_level: debug
---
> log_level: info
3c3
< timeout: 30
---
> timeout: 60
```

```bash
cat <<EOF > /workspace/answers.env
DATABASE_HOST_OLD=db-old.internal
DATABASE_HOST_NEW=db-new.internal
LOG_LEVEL_OLD=debug
LOG_LEVEL_NEW=info
TIMEOUT_OLD=30
TIMEOUT_NEW=60
EOF
```

`diff`'s output names the three changed keys directly — `max_connections`
and `retry_count` don't appear because they're identical in both files.
