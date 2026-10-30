# Solution: Disk Full Incident

```bash
# Narrow down what's actually consuming space.
du -sh /var/log/* 2>/dev/null | sort -rh
du -sh /var/log/app/* 2>/dev/null | sort -rh

# /var/log/app/debug is disposable, generated debug output — remove it.
rm -rf /var/log/app/debug

# Confirm the protected fixtures are untouched.
cat /srv/data/customer.db
cat /opt/app/bin/app
cat /var/log/audit/retained.log

# Confirm the application can write again.
touch /var/log/app/new.log
```

`du -sh` on successive levels narrows down which specific directory is
responsible before deleting anything. Only `/var/log/app/debug` —
clearly generated, disposable debug output — gets removed; the three
fixtures the ticket calls out as protected (`customer.db`, the
application binary, and the retained audit log) are never touched, and
the final `touch` proves the application's actual failure mode (it
couldn't write new logs) is now resolved, not just that disk usage went
down.
