# A mount replaces what's underneath it, it doesn't merge with it

A Dockerfile's `RUN` and `COPY` instructions write into the image's own
filesystem layers. Those layers are real, and a container started from
that image with nothing else attached really does see whatever the
Dockerfile wrote there.

The moment a volume — named or bind-mounted — is attached at some path,
that changes completely. Linux's mount mechanism doesn't blend a volume's
contents with whatever already existed at that path in the image; it
replaces that path's view entirely, for as long as the mount stays
attached. The image's own content at that path still exists, down in the
image's layers — but nothing reachable from inside the running container
can see it anymore, because the volume is now what answers every read and
write at that path.

## Why this surprises people

"The Dockerfile clearly created this file" is a true statement about the
image. It is not evidence about what a running container's filesystem
actually contains, because a volume mounted at that path overrides it
completely. The two most common first reactions are also the two wrong
ones:

- Rebuilding the image does nothing — the new image's layers still get
  shadowed by the same mount on the very next start.
- Inspecting the image's own filesystem (for example exporting its
  layers) proves what *would* be there with no mount attached, which is
  not the same question as what's actually there right now, with the
  mount attached.

## Proving two paths are genuinely different mounts

`stat -c '%d' <path>` prints the numeric device id of the filesystem a
path actually lives on. Two paths on the same filesystem, with no mount
boundary between them, report the same device id; a path that's a
distinct mount point reports a different device id than its parent. This
is a factual check, not an inference — the same thing a full mount table
inspection would confirm, in one command.
