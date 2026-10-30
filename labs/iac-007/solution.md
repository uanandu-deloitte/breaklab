# Solution: Why You Pin Module Versions

```bash
declared=$(sed -n 's/^module: *//p' /infra/resources.yaml)   # render-v2
pinned=$(cat /infra/MODULE_VERSION)                            # render-v1

cat /etc/app/release.conf
# RELEASE(v1): build-42        <- rendered by the stale pin, not render-v2

echo "$declared" > /infra/MODULE_VERSION   # fix the pin

infra-apply
# applied: /etc/app/release.conf (rendered by render-v2)

cat /etc/app/release.conf
# Release build-42 [rendered by v2]

cat <<EOF > /workspace/answers.env
DECLARED_MODULE=$declared
PINNED_MODULE=$pinned
MODULES_AGREED_BEFORE_FIX=no
RESOLVED_VERSION_AFTER_FIX=$declared
EOF
```

`resources.yaml` declares `module: render-v2`, but `infra-apply` never
reads that field — it only reads `/infra/MODULE_VERSION`, which was
still pinned to `render-v1`. The apply ran cleanly and exited 0 both
before and after the fix; a clean exit never told you which module ran,
only that whichever one did run didn't crash. Comparing the deployed
file's format against each module's source was what actually revealed
`render-v1` had produced it. Fixing the pin — not `resources.yaml`,
which was already correct — and re-running `infra-apply` is what made
the resolved module finally match the declared one.
