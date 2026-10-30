# Config Precedence: File vs Environment vs Flag

This container has a small infra-as-code setup:

- `/infra/resources.yaml` declares a `value:` for a managed file.
- `/etc/app/mode.conf` is the real file `infra-apply` writes.
- `/usr/local/bin/infra-apply` can read the value it applies from three
  places: a `--value=...` CLI flag, the `$INFRA_VALUE` environment
  variable, or the `value:` field in `resources.yaml` — in that order
  of priority. Only one of the three ever actually wins.

Right now, the config file and the environment disagree about what the
value should be.

## Tasks

1. `cat /infra/resources.yaml` — see what the file declares.
2. `env | grep INFRA_VALUE` — see what the environment declares.
3. Run `infra-apply` with no flag and read its own output carefully —
   it reports which source it actually used.
4. `cat /etc/app/mode.conf` — confirm the real resource matches
   whichever source actually won, not necessarily the file.
5. Run `infra-apply --value=maintenance-mode-forced` and check
   `/etc/app/mode.conf` again, to see the flag beat both of the other
   two sources.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
FILE_VALUE=
ENV_VALUE=
ACTUAL_SOURCE=
RESOLVED_VALUE=
```

- `FILE_VALUE` — the value declared in `resources.yaml`.
- `ENV_VALUE` — the value of `$INFRA_VALUE` in this container.
- `ACTUAL_SOURCE` — which one `infra-apply` actually used when run with
  no flag: `file` or `env`.
- `RESOLVED_VALUE` — the value `infra-apply` actually wrote, when run
  with no flag.

Get a shell: `docker exec -it breaklab-iac-006 bash`
