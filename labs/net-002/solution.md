# Solution: Who Owns the Port?

```bash
# Confirm the diagnosis: who is actually listening on port 80?
ss -ltnp | grep :80
# LISTEN ... *:80 ... users:(("python3",pid=...,fd=...))
ps aux | grep python3

# Stop only the decoy process.
pkill -f "http.server 80"
sleep 1

# Now the port is free — start the real web server.
nginx -t
nginx

# Confirm.
curl http://localhost:8080/
```

`ss -ltnp` shows the port is held by a `python3` process, not nginx —
that's the decoy. `pkill -f "http.server 80"` targets only that specific
command line, leaving everything else untouched. Only once the port is
actually free does `nginx` succeed in binding it.
