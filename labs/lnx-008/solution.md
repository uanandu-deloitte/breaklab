# Solution: Find the Stale Backend

```bash
# Find every occurrence first — read them before touching anything.
grep -RlF "database-old.internal" /opt/app
# /opt/app/config/app.conf                  <- the live config (the bug)
# /opt/app/releases/v1/config/app.conf       <- archived release snapshot
# /opt/app/CHANGELOG.md                      <- historical record
# /opt/app/docs/migration-notes.md           <- historical record

# Edit only the live config.
sed -i 's/database-old\.internal/database.internal/' /opt/app/config/app.conf

# Self-check.
/opt/app/bin/check-config.sh   # OK: active config points at database.internal
```

Only `/opt/app/config/app.conf` is the file the running service reads —
everything else that matched is either an archived snapshot of a
previous release or documentation describing the migration
historically, and both are correct exactly as they are. `sed -i` is
scoped to that one file explicitly, never run with a glob or `-R` across
`/opt/app`, so nothing else in the tree is touched.
