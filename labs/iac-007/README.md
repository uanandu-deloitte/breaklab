# Why You Pin Module Versions

This container has a small infra-as-code setup:

- `/infra/resources.yaml` declares a resource, including a `module:`
  field naming which render module it expects to be used.
- `/infra/MODULE_VERSION` is a separate pin file. `infra-apply` reads
  *this* file — not the `module:` field in `resources.yaml` — to
  decide which render module to actually run.
- `/infra/modules/` holds the render module implementations
  themselves (`render-v1.sh`, `render-v2.sh`), each producing a
  slightly different output format.
- `/etc/app/release.conf` is the real, deployed output.

`infra-apply` already ran once when this image was built. It exited
cleanly. That doesn't tell you which module it actually used.

## Tasks

1. `cat /infra/resources.yaml` — note the `module:` field it declares.
2. `cat /infra/MODULE_VERSION` — note what the pin file actually says.
3. Compare the two. Do they agree?
4. `cat /etc/app/release.conf` and compare its format against
   `/infra/modules/render-v1.sh` and `render-v2.sh` to figure out which
   module actually rendered it.
5. Fix `/infra/MODULE_VERSION` so it matches what `resources.yaml`
   declares, then run `infra-apply` again.
6. Confirm `/etc/app/release.conf` is now rendered by the module that
   was actually declared.

Record your findings in `/workspace/answers.env`, one `KEY=value` per
line:

```
DECLARED_MODULE=
PINNED_MODULE=
MODULES_AGREED_BEFORE_FIX=
RESOLVED_VERSION_AFTER_FIX=
```

- `DECLARED_MODULE` — the module named in `resources.yaml`'s `module:`
  field.
- `PINNED_MODULE` — what `/infra/MODULE_VERSION` originally said,
  before you touched it.
- `MODULES_AGREED_BEFORE_FIX` — `yes` or `no`: did the declared module
  and the original pin actually match?
- `RESOLVED_VERSION_AFTER_FIX` — the module that actually renders
  `/etc/app/release.conf` now, after your fix.

Get a shell: `docker exec -it breaklab-iac-007 bash`
