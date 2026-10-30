# Two ways to get a ConfigMap into a container

A ConfigMap just stores key/value pairs — it does nothing on its own.
A Pod has to explicitly ask for each key, and there are two independent
ways to ask:

**As an environment variable.** A container's `env[]` entry can set
`valueFrom.configMapKeyRef` instead of a literal `value`, naming the
ConfigMap and the key to pull from it. The environment variable's own
name doesn't have to match the ConfigMap key's name — `APP_LOG_LEVEL`
can be sourced from a key called `LOG_LEVEL`. This is the right choice
for small values an application reads once at startup, like a log
level or a feature flag.

**As a mounted file.** A Pod's `volumes[]` can declare a `configMap`
volume, and a container's `volumeMounts[]` mounts it at some path. Each
key in the ConfigMap becomes a file named after that key, containing
its value — unless the volume's `items[]` restricts it to specific
keys (and lets you rename the file with `path`). This is the right
choice for anything that looks like a config *file* — multi-line
content an application expects to `open()` and read, like `app.conf`.

## Why both exist on purpose

These aren't competing options where one is "better" — they match
different consumption patterns. A twelve-factor app that reads
`os.Getenv("LOG_LEVEL")` wants the environment variable. A process
that expects a real file path on disk, like `-c /etc/app/app.conf`,
wants the volume mount. The same ConfigMap can back both at once,
exactly as this lab's `app-config` does — nothing stops one key from
going out as an env var while a different key in the same object goes
out as a file.

## Why this needs confirming live, not just reading the manifest

A manifest *declaring* `configMapKeyRef` or a `configMap` volume
proves intent, not effect. The key might not exist in the ConfigMap,
the mount path might be wrong, or the volume might never have been
attached to the container at all. The only way to know an injection
actually took is to look from inside the running Pod — `printenv` for
the environment variable, `cat` for the file — and compare what you
see against the ConfigMap's current live content.
