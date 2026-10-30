# Why Did the Script Fail?

`/opt/app/start.sh` is supposed to start the app using a configuration
file named by the `APP_CONFIG_PATH` environment variable. Operations
reports that it "always works" — it always prints a success message and
exits 0 — but the app behind it has been failing to actually start for
days.

## Tasks

Repair `/opt/app/start.sh` so that:

- Run with `APP_CONFIG_PATH` **unset**, it exits **non-zero** and prints a
  clear, safe error message (no leaking anything that shouldn't be
  exposed) — not a silent, misleading success.
- Run with `APP_CONFIG_PATH` pointing at a real, valid config file (one is
  already seeded at `/opt/app/config.yaml`), it genuinely succeeds:
  prints `Service started successfully` and exits 0.

## Verifying it yourself

```bash
env -u APP_CONFIG_PATH bash /opt/app/start.sh; echo "exit: $?"
APP_CONFIG_PATH=/opt/app/config.yaml bash /opt/app/start.sh; echo "exit: $?"
```

Get a shell: `docker exec -it breaklab-lnx-019 bash`
