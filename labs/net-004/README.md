# Packet Inspection Introduction

A small HTTP backend is running locally on port 8005. Rather than trust
`curl`'s own report of what happened, capture the actual packets.

## Tasks

- capture traffic on the loopback interface while making exactly one
  `curl` request to the backend, writing the capture to
  `/workspace/capture.pcap`
- read the capture back and confirm you can see both the request and the
  response
- identify the destination port of the request

Record your findings in `/workspace/answers.env`:

```
DEST_PORT=
CAPTURE_HAS_REQUEST=
CAPTURE_HAS_RESPONSE=
```

(`CAPTURE_HAS_REQUEST` / `CAPTURE_HAS_RESPONSE` are each `yes` or `no`.)

A reasonable sequence:

```bash
tcpdump -i lo -w /workspace/capture.pcap port 8005 &
sleep 1
curl http://localhost:8005/
sleep 1
kill %1
```

Get a shell: `docker exec -it breaklab-net-004 bash`
