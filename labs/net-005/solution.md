# Solution: Broken Name Resolution

```bash
# Confirm the diagnosis: the name itself doesn't resolve.
getent hosts api.internal          # prints nothing, exit status 2

# Confirm the backend is actually fine on its own.
curl http://127.0.0.1:8080/        # {"service": "api-internal", "status": "ok"}

# Fix name resolution — add the missing entry to /etc/hosts.
echo "127.0.0.1 api.internal" >> /etc/hosts

# Confirm the fix.
getent hosts api.internal          # now prints 127.0.0.1  api.internal
curl http://api.internal:8080/     # {"service": "api-internal", "status": "ok"}
```

The backend was never the problem — it was reachable on `127.0.0.1:8080`
the whole time. The failure was entirely in name resolution: nothing told
this host what address `api.internal` points to, so the lookup failed
before any connection was attempted. Adding the `/etc/hosts` entry fixes
the resolver, not the request — curling `127.0.0.1` directly would have
"worked" without actually fixing anything a real frontend (which only
knows the hostname) could rely on.
