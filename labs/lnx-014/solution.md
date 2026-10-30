# Solution: Permission Playground

```bash
# Verify each prediction as tester, not as root.
su -s /bin/bash -c "test -r /workspace/playground/owner-only.txt && echo yes || echo no" tester
su -s /bin/bash -c "test -r /workspace/playground/group-readable.txt && echo yes || echo no" tester
su -s /bin/bash -c "test -r /workspace/playground/world-readable.txt && echo yes || echo no" tester

cat <<EOF > /workspace/answers.env
OWNER_ONLY_READABLE_BY_TESTER=no
GROUP_READABLE_READABLE_BY_TESTER=no
WORLD_READABLE_READABLE_BY_TESTER=yes
EOF
```

`owner-only.txt` is mode 600, owned by `root` — only root's own access
matters, so `tester` gets nothing. `group-readable.txt` is mode 640,
group `opsgroup` — `tester` isn't a member of `opsgroup`, so the group
triplet (which does allow read) never applies to them, and they fall
through to "other," which allows nothing. `world-readable.txt` is mode
644 — `tester` is neither the owner nor in the group, so the "other"
triplet applies, and it grants read.
