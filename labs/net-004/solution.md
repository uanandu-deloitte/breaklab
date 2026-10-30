# Solution: Packet Inspection Introduction

```bash
tcpdump -i lo -w /workspace/capture.pcap port 8005 &
sleep 1
curl http://localhost:8005/
sleep 1
kill %1
sleep 1

# Read it back:
tcpdump -r /workspace/capture.pcap -nn
# 127.0.0.1.54321 > 127.0.0.1.8005: Flags [S], ...        <- request side
# 127.0.0.1.8005 > 127.0.0.1.54321: Flags [S.], ...       <- response side
# 127.0.0.1.54321 > 127.0.0.1.8005: Flags [P.], ... GET...
# 127.0.0.1.8005 > 127.0.0.1.54321: Flags [P.], ... HTTP/1.0 200...

cat <<EOF > /workspace/answers.env
DEST_PORT=8005
CAPTURE_HAS_REQUEST=yes
CAPTURE_HAS_RESPONSE=yes
EOF
```

The client's own port (`54321` above) is whatever ephemeral port the
kernel happened to assign — it varies every run. Port 8005 is constant:
it's the *destination* on every packet the client sends, and the
*source* on every packet the server sends back. Seeing both directions
in the capture is what proves the full exchange happened, not just that
a request was sent.
