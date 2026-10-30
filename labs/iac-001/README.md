# Declarative vs Imperative

This container has a tiny infra-as-code setup:

- `/infra/resources.yaml` declares the desired content of a file.
- `/etc/banner.txt` is the real file `infra-apply` writes to match it.
- `/usr/local/bin/infra-apply` reads `resources.yaml` and overwrites the
  real file to match the declared `content:`.

`infra-apply` already ran once when the image was built, so
`/etc/banner.txt` currently matches `resources.yaml`.

## Tasks

1. `cat /infra/resources.yaml` and `cat /etc/banner.txt` — confirm they
   currently match.
2. Edit the `content:` line in `/infra/resources.yaml` to a new value of
   your choice.
3. Run `infra-apply`.
4. `cat /etc/banner.txt` again and confirm it now shows your new value —
   not because you edited the file yourself, but because the tool
   reconciled it to match the declaration.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
DESIRED_VALUE_BEFORE=
DESIRED_VALUE_AFTER=
SOURCE_OF_TRUTH_FILE=
```

- `DESIRED_VALUE_BEFORE` — what `content:` said before you edited it.
- `DESIRED_VALUE_AFTER` — what `content:` says now (and what
  `/etc/banner.txt` should now contain).
- `SOURCE_OF_TRUTH_FILE` — the path of the file that's actually in
  charge here (not the file you'd naively expect if you'd never seen an
  apply tool before).

`DESIRED_VALUE_BEFORE` and `DESIRED_VALUE_AFTER` contain spaces — quote
them (`KEY="value with spaces"`) so sourcing the file doesn't word-split
it.

Get a shell: `docker exec -it breaklab-iac-001 bash`
