# docker logs only ever shows one thing

`docker logs` looks like a general-purpose window into "whatever this
container is doing," but it is much narrower than that. It shows exactly
one thing: whatever PID 1 (and any process that inherited PID 1's own
stdout/stderr file descriptors) writes to **stdout and stderr**. That's
the entire contract. Docker's logging driver attaches to those two
streams the moment the container starts, and nothing else.

## What this means in practice

- A process can be completely healthy, actively producing output every
  second, and `docker logs` can show absolutely nothing — if that output
  is going somewhere other than stdout/stderr.
- A file written to disk inside the container is real. You can `cat` it,
  `tail -f` it, watch it grow. None of that makes it visible to
  `docker logs`. Docker's logging driver has no idea the file exists; it
  was never attached to it in the first place.
- This isn't a bug or a misconfiguration to "fix" in Docker — it's the
  whole mechanism. If an application (or whatever launched it) chooses to
  log to a file instead of stdout, `docker logs` was never going to show
  it, no matter how long you wait or how many times you restart.

## Why this surprises people

"The app is clearly alive and logging something" is true and also not
evidence that `docker logs` will show any of it. The two facts are
decided independently: whether the process is healthy, and where its
output is actually being written. Confirming the first tells you nothing
about the second — you have to go look at the container's filesystem
directly to find out where output is really landing.

## How to check, instead of assuming

- `docker logs <container>` run from the **host** shows the stdout/stderr
  capture — this is the only thing that command can ever show.
- A shell inside the container (`docker exec -it <container> bash`) lets
  you inspect the filesystem directly and find any file a process might
  be writing to instead of, or in addition to, stdout.

Both are real signals. Only one of them is what `docker logs` will ever
report.
