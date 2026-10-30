# Solution: Nothing Changes When Nothing Changed

```bash
# 1. Nothing edited yet — this should be a no-op.
infra-apply
# -> applied: /etc/app-release.conf (no changes)

# 2. Edit the declared value.
sed -i 's/^value:.*/value: release-2026.10.2/' /infra/resources.yaml

# 3. Re-run — this time it's a real change.
infra-apply
# -> applied: /etc/app-release.conf (changed)

cat <<EOF > /workspace/answers.env
FIRST_RUN_RESULT=no-change
SECOND_RUN_RESULT=changed
FINAL_VALUE=release-2026.10.2
EOF
```

The first run is a no-op because `/etc/app-release.conf` already held
exactly what `/infra/resources.yaml` declared — nothing to do, nothing
touched. Only after the declared value itself changes does the second
run have real work to do, and `infra-apply` reports that difference
plainly instead of silently rewriting a file that was already correct.
