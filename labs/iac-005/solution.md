# Solution: Lockfiles and Crash Safety

```bash
infra-apply
# refusing to run: lock held by PID 99999 (started 2026-10-21T03:14:00Z)

lock_pid=$(sed -n 's/^PID=//p' /var/run/infra/apply.lock)

ps -p "$lock_pid"
# (no output, non-zero exit — no such process)

stat /var/run/infra/apply.lock
# shows a modification time well before "now"

rm /var/run/infra/apply.lock

infra-apply
# applied: /etc/release-notes.txt

cat <<EOF > /workspace/answers.env
LOCK_PID=$lock_pid
LOCK_PROCESS_ALIVE=no
APPLY_SUCCEEDED=yes
EOF

cat /etc/release-notes.txt   # now prints "Build server release notes v3"
```

The lock recorded PID 99999 and an old timestamp. `ps -p 99999` finds no
such process — whatever created the lock is long gone, and it never ran
the cleanup step that would have removed its own lock on exit. That
makes this a stale lock, not a live one: safe to remove. Once it's
gone, `infra-apply` can create its own lock, apply the change, and
(this time) clean up after itself when it exits — leaving
`/etc/release-notes.txt` matching `resources.yaml`'s declared content
and no lock file behind.
