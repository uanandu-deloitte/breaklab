# Writing a health-check script worth trusting

A health-check script is only useful if its exit code can be trusted
without reading its output — anything that calls it (a supervisor, a
load balancer, a human running it at 3am) needs `$?` alone to mean
something.

## Check what you actually depend on, specifically

"Is the app healthy" usually decomposes into several independent
things: is the right process running, is it listening where it's
supposed to, does it actually answer requests correctly, and is there
room left to keep operating. Checking only the easiest one (e.g. just
the HTTP response) misses failure modes the others catch — a process can
still accept TCP connections while unable to serve a correct response,
or still be running while its disk fills up underneath it.

- A process check (`pgrep -f <pattern>`) confirms something matching
  is alive, independent of whether it's network-reachable yet.
- A port check (`ss -ltn`, or bash's own `/dev/tcp/<host>/<port>`
  redirection) confirms a listener exists, independent of whether its
  response is correct.
- An HTTP check (`curl`) confirms the full round trip actually works.
- A disk check (`du -sh` against the specific path the app writes to)
  catches the failure mode that's invisible to all three checks above
  until it's too late.

## Exit codes and concise output

Exit `0` only when every check passed. Exit non-zero — and print exactly
one line naming what failed — the moment any single check fails; there's
no reason to keep checking once the answer is already "unhealthy."
Resist the urge to print every intermediate step: a script whose healthy
output is one line and whose unhealthy output is one clear line is far
more useful than one that dumps raw command output for someone to
interpret.
