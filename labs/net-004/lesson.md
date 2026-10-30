# Looking below the application layer

`curl` tells you what *it* thinks happened: a status code, a body, maybe
a timing summary. That's a report from one of the two participants. A
packet capture is the ground truth of what actually went over the wire
between them — useful the moment you stop trusting (or need to prove)
what either side claims.

## Capturing with tcpdump

`tcpdump` reads raw packets directly from a network interface. On a
single host talking to itself, that interface is `lo` (loopback) — every
local request/response pair you make to `localhost` still produces real
packets on `lo`, even though they never reach a physical network.

```
tcpdump -i lo -w /workspace/capture.pcap port 8005
```

- `-i lo` — capture on the loopback interface
- `-w <file>` — write raw packets to a file instead of printing them
  (a `.pcap` file, readable by `tcpdump -r` or tools like Wireshark)
- `port 8005` — a capture filter: only packets involving port 8005 in
  either direction, so the file doesn't fill up with unrelated noise

Running it in the background (`&`), waiting a moment for it to actually
attach before generating traffic, and waiting another moment before
stopping it, matters: a capture that starts after the request began or
stops the instant the client gets its response can miss the last few
packets of either side.

## Reading a capture back

```
tcpdump -r /workspace/capture.pcap -nn
```

- `-r <file>` — read from a file instead of a live interface
- `-nn` — don't resolve hostnames or port names, show raw numbers (so
  port 8005 prints as `8005`, not some guessed service name)

Each line names a source `IP:port` and a destination `IP:port`. A full
exchange shows packets in *both* directions: your client's SYN and
request going to port 8005, and the server's response packets coming
back from port 8005 to whatever ephemeral source port your client used
for this connection.

## Source port vs. destination port

The server's well-known port (8005 here) is the *destination* port on
every packet your client sends, and the *source* port on every packet
the server sends back. The client itself uses a throwaway, high-numbered
"ephemeral" port for its own side of the connection — that's expected,
and different every time. Confusing the two is a common mistake when
reading output for the first time.

## Why this matters beyond this one capture

Once you can read a capture confidently, you stop having to take either
side's word for what happened during a network incident — you can see
whether a request ever left, whether a response ever came back, or
whether retransmissions mean packets were being dropped along the way.
