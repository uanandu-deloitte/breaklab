# Connectivity Ladder

This host exposes four local endpoints. Each one fails (or succeeds) in a
genuinely different way. Diagnose them **in this order**:

1. `curl http://localhost:8001/`
2. `curl http://nonexistent.internal:8002/`
3. `curl http://localhost:8003/`
4. `curl --max-time 3 http://localhost:8004/`

Use `ping`, `ss`, `nc -zv`, `curl -v`, and route info to work out exactly
*why* each one behaves the way it does — don't just note the symptom.

Record your classification in `/workspace/answers.env`, one `KEY=value`
per line, using exactly one of `healthy`, `dns-failure`, `refused`, or
`timeout` for each:

```
ENDPOINT_1=
ENDPOINT_2=
ENDPOINT_3=
ENDPOINT_4=
```

Get a shell: `docker exec -it breaklab-net-003 bash`
