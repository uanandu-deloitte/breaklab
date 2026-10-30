# Solution: Resolve the Conflict

```bash
cd /workspace
git checkout main
git merge feature/raise-timeout
# CONFLICT (content): Merge conflict in config.yaml

cat <<EOF > config.yaml
service: acme-api
port: 8080
timeout: 60
EOF

git add config.yaml
git commit --no-edit
```

`git merge feature/raise-timeout` stops mid-merge because both branches
changed the `timeout:` line differently since their common ancestor —
git has no way to pick a side on its own. Overwriting `config.yaml` with
the agreed final content (keeping the longer timeout, discarding main's
shorter one) removes both the conflict and all of its marker lines in
one step. `git add` marks the file resolved, and `git commit` (git
already prepared a merge commit message, so `--no-edit` just accepts it)
finishes the merge with a real two-parent commit.
