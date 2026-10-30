# The Volume That Was Already There

This image's Dockerfile wrote a welcome file at `/data/welcome.txt` during
the build, with this exact content:

```
Welcome to the shared data volume.
```

But this lab's environment mounts a fresh volume at `/data` — so whatever
you find there right now did not come from the Dockerfile. Figure out why,
and restore the file.

## Tasks

1. Check what's actually in `/data` right now, and compare it with what
   the Dockerfile wrote.
2. Recreate `/data/welcome.txt` with exactly the content shown above — the
   volume, not the image, is what actually persists for this container,
   so writing it there is the only way to get the file back for real.
3. Prove `/data` is a genuinely distinct mount point from `/`, by
   comparing `stat -c '%d' /data` against `stat -c '%d' /`.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
WELCOME_FILE_RESTORED=
MOUNT_IS_DISTINCT=
```

- `WELCOME_FILE_RESTORED`: `yes` once `/data/welcome.txt` exists again
  with the exact original content.
- `MOUNT_IS_DISTINCT`: `yes` or `no`, based on whether the two `stat -c
  '%d'` device ids from task 3 actually match or differ.

Get a shell: `docker exec -it breaklab-docker-005 bash`
