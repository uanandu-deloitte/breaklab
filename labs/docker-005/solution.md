# Solution: The Volume That Was Already There

```bash
cat > /data/welcome.txt <<'EOF'
Welcome to the shared data volume.
EOF

data_dev=$(stat -c '%d' /data)
root_dev=$(stat -c '%d' /)

cat <<EOF2 > /workspace/answers.env
WELCOME_FILE_RESTORED=yes
MOUNT_IS_DISTINCT=$([ "$data_dev" != "$root_dev" ] && echo yes || echo no)
EOF2
```

The Dockerfile's `RUN echo ... > /data/welcome.txt` wrote into the image's
own layer, but the volume this lab's environment declares at `/data`
mounts a separate, empty filesystem over that same path the moment the
container starts — so the image's file is unreachable from inside this
container, not merged with the volume's (empty) content. Recreating the
file writes directly onto the volume, which is the only copy that will
actually persist across a reset or restart.

`stat -c '%d'` reports the device id of the filesystem backing a path. A
real mount boundary at `/data` reports a different device id than `/`
itself, confirming this isn't just an empty subdirectory on the same
filesystem but an actual separate mount.
