# Preview Before You Apply

This container has a small infra-as-code setup, same shape as before:

- `/infra/resources.yaml` declares the desired content of a file.
- `/infra/managed-file.txt` is the real file.
- `/usr/local/bin/infra-apply` reconciles the real file to match the
  declaration.

This time, `resources.yaml` was already edited to a new value *after*
the last real apply — the real file still holds the old content.
`infra-apply` also supports a `--plan` flag: run with `--plan` and it
only prints what it *would* change, without touching anything; run with
no flag and it actually applies the change.

## Tasks

1. `cat /infra/resources.yaml` and `cat /infra/managed-file.txt` —
   confirm they currently differ.
2. Run `infra-apply --plan`. Read its output.
3. `cat /infra/managed-file.txt` again — confirm it's still exactly
   what it was before you ran `--plan`. A plan never writes.
4. Run `infra-apply` with no flag.
5. `cat /infra/managed-file.txt` again — confirm it now matches
   `resources.yaml`'s declared content.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
PLANNED_CHANGE=
```

- `PLANNED_CHANGE` — exactly what `infra-apply --plan` reported it would
  change, before you ran the real apply.

`PLANNED_CHANGE` contains spaces — quote it (`KEY="value with spaces"`)
so sourcing the file doesn't word-split it.

Get a shell: `docker exec -it breaklab-iac-002 bash`
