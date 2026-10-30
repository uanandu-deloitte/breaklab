# What Even Is a Container

You have shell access inside a running container. Before anything else,
establish which parts of what you're looking at are actually isolated to
this container, and which parts are just the host showing through.

## Tasks

1. Count how many processes `ps aux` shows inside this container. It should
   be a small handful — nowhere near the hundreds a real host would have.
2. Check `uname -r`. This is **not** a kernel that shipped with this
   image — containers don't bundle a kernel at all. Figure out whose
   kernel this actually is.
3. Check `hostname` (and `/etc/hostname`). Figure out where this value
   came from — it was never written by this lab's Dockerfile.

Record your findings in `/workspace/answers.env`, one `KEY=value` per line:

```
PROCESS_COUNT_SEEN=
KERNEL_SHARED=
HOSTNAME_SOURCE=
```

- `PROCESS_COUNT_SEEN`: the number of lines `ps aux --no-headers` prints
  right now, in this container.
- `KERNEL_SHARED`: `yes` or `no` — is the kernel version you see shared
  with (i.e. actually belonging to) the host, rather than bundled in this
  image?
- `HOSTNAME_SOURCE`: a short phrase naming what actually set this
  container's hostname (for example: `container-runtime` or
  `container-id`) — not `dockerfile`, since nothing in this image's
  Dockerfile sets a hostname.

Get a shell: `docker exec -it breaklab-docker-001 bash`
