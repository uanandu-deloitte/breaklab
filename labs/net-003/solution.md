# Solution: Connectivity Ladder

```bash
# Endpoint 1 — healthy
curl -v http://localhost:8001/
# HTTP/1.1 200 OK

# Endpoint 2 — DNS failure
curl -v http://nonexistent.internal:8002/
# curl: (6) Could not resolve host: nonexistent.internal

# Endpoint 3 — connection refused
curl -v http://localhost:8003/
# curl: (7) Failed to connect to localhost port 8003: Connection refused
nc -zv localhost 8003   # also reports "Connection refused"
ss -ltn | grep :8003    # nothing — no process has this port

# Endpoint 4 — timeout
curl --max-time 3 -v http://localhost:8004/
# curl: (28) Operation timed out after 3000 milliseconds
ss -ltn | grep :8004    # something IS listening...
nc -zv localhost 8004   # ...and the TCP connection itself succeeds

cat <<EOF > /workspace/answers.env
ENDPOINT_1=healthy
ENDPOINT_2=dns-failure
ENDPOINT_3=refused
ENDPOINT_4=timeout
EOF
```

Each endpoint fails at a different stage: 8002 never gets past name
resolution; 8003 gets an immediate, explicit refusal because nothing
holds that port; 8004 gets a successful TCP connection but no
application-layer response, so nothing happens until the client's own
timeout fires.
