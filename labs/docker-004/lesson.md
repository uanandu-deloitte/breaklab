# Images bake in defaults; containers can override them

A Dockerfile's `ENV` instruction writes a value into the image itself — it
becomes part of the image's own config, fixed at build time. Every
container started from that image inherits that value as its starting
environment, whether or not anyone ever overrides it later.

## Where to look at runtime

- `env` and `printenv` show the current process's environment as the shell
  sees it.
- `/proc/1/environ` holds PID 1's own environment straight from the
  kernel — no shell involved. It is NUL-separated rather than
  newline-separated, so reading it directly with `cat` looks like one
  unbroken line; translating each NUL byte into a newline makes it
  readable one variable per line.
- These normally agree with each other — the container starts with one
  environment that every process in it inherits — but `/proc/1/environ` is
  the ground truth when something downstream might have re-exported or
  changed a value in its own shell session.

## Build time vs. run time

A value baked in with `ENV` in the Dockerfile is fixed the moment the image
is built. It changes only if:

- the image is rebuilt with a different Dockerfile, or
- something overrides it at container start (for example `docker run -e
  KEY=value`), or
- a process inside the container explicitly re-exports it in its own
  shell.

None of those is "the Dockerfile's value silently changing on its own" —
the baked-in value sits inert until something outside the image actively
overrides it.

## Which values belong in an image at all

Some values genuinely describe the build itself — a version string tied to
what was actually compiled — and are reasonable to bake in, because
"different image, different value" is exactly the behavior you want.

Other values describe the *deployment*, not the build: how verbose to log,
which backend to talk to, which feature flags are on. Baking those into
the image means the only way to change them is to rebuild and re-publish a
new image, which defeats the point of having one image you deploy to
multiple environments. Those values belong in a runtime override
(`docker run -e`, a Kubernetes ConfigMap, and so on), with the Dockerfile's
`ENV` only supplying a sane default.
