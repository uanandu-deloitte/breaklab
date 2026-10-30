# Build Time vs Run Time

This image's Dockerfile bakes three `ENV` values into it at build time.
Find out what they are, and decide which one should really be supplied at
deploy time instead of hardcoded into the image forever.

## Tasks

1. Read this container's actual running environment — with `env` or
   `printenv`, and directly from `/proc/1/environ` (remember: it's
   NUL-separated, not newline-separated).
2. Identify the three baked-in values: `APP_ENV`, `LOG_LEVEL`, and
   `APP_VERSION`.
3. Decide which *one* of those three looks like it should normally be
   overridden per-deployment rather than baked into the image once and
   forgotten.

Record your findings in `/workspace/answers.env`, one `KEY=value` per line:

```
APP_ENV=
LOG_LEVEL=
APP_VERSION=
OVERRIDE_CANDIDATE=
```

`OVERRIDE_CANDIDATE` is the name of the one baked-in key (just the key
name, e.g. `LOG_LEVEL`) that makes the least sense frozen into the image
at build time.

Collect every value from the running container itself — don't type in
what you expect to see.

Get a shell: `docker exec -it breaklab-docker-004 bash`
