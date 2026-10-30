# Solution: State Files Only Know What You Told Them

```bash
infra-apply
# applied: /etc/app/app.conf (no changes)
# skipping /etc/cache/cache.conf: not in state, run with --import first

cat /infra/state.json
# {"managed": ["/etc/app/app.conf"]}

cat /etc/cache/cache.conf
# cache config v0-manual      <- never matched resources/cache.yaml

infra-apply --import /etc/cache/cache.conf
# imported: /etc/cache/cache.conf

infra-apply
# applied: /etc/app/app.conf (no changes)
# applied: /etc/cache/cache.conf

cat /etc/cache/cache.conf
# cache config v1

cat <<EOF > /workspace/answers.env
UNMANAGED_PATH=/etc/cache/cache.conf
APPLY_SUCCEEDED_AFTER_IMPORT=yes
EOF
```

`resources/cache.yaml` was declared correctly from the start, and
`infra-apply` ran without ever erroring — but `/infra/state.json`'s
managed list only ever contained `/etc/app/app.conf`. Every run quietly
skipped `cache.conf`'s declaration instead of failing, because that's
the tool's deliberate behavior for anything outside its tracked state.
Running `infra-apply --import /etc/cache/cache.conf` added the path to
the managed list; only after that did a plain `infra-apply` finally
write the declared content for the first time.
