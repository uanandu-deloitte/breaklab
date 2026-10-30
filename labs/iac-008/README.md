# State Files Only Know What You Told Them

This container has a small infra-as-code setup:

- `/infra/resources/` holds one declaration file per managed resource
  (`app.yaml`, `cache.yaml`), each with a `path:` and `content:`.
- `/infra/state.json` lists every path `infra-apply` is actually
  allowed to touch, under `"managed"`.
- `infra-apply` loops over every file in `/infra/resources/` and
  applies it — but only if its declared path is also listed in
  `state.json`. If it isn't, `infra-apply` prints a skip message and
  moves on, without touching anything.
- `infra-apply --import <path>` adds a path to `state.json`'s managed
  list.

`infra-apply` already ran once when this image was built, and it
didn't error.

## Tasks

1. Run `infra-apply` and read its output for both declared resources.
2. `ls /infra/resources/` and `cat /infra/state.json` — compare what's
   declared against what's actually tracked.
3. Find the one declared resource that's missing from the managed
   list.
4. `cat` its real file on disk — confirm it's never matched its
   declared content.
5. Run `infra-apply --import <that path>`, then run `infra-apply`
   again.
6. Confirm the resource now matches its declaration.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
UNMANAGED_PATH=
APPLY_SUCCEEDED_AFTER_IMPORT=
```

- `UNMANAGED_PATH` — the path of the declared resource that was never
  in `state.json`'s managed list.
- `APPLY_SUCCEEDED_AFTER_IMPORT` — `yes` or `no`: after importing it,
  did `infra-apply` actually write its declared content for real?

Get a shell: `docker exec -it breaklab-iac-008 bash`
