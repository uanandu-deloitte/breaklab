# Solution: Identify Repo State

```bash
cd /workspace
git status
git branch
git log --oneline

cat <<EOF > answers.env
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
HEAD_COMMIT_SHORT_SHA=$(git rev-parse --short HEAD)
OTHER_BRANCH_NAME=$(git branch --format='%(refname:short)' | grep -v -x "$(git rev-parse --abbrev-ref HEAD)")
MODIFIED_FILE_COUNT=$(git status --porcelain | grep -c '^ M ')
EOF
```

`git rev-parse --abbrev-ref HEAD` names the checked-out branch directly
instead of parsing it out of `git status` output. `git rev-parse --short
HEAD` gives the exact short SHA `git log --oneline`'s top line shows.
`git branch --format` lists every local branch's name with no `*`
decoration, so filtering out the current one leaves the other branch
that exists but isn't checked out. `git status --porcelain` prints a
stable two-character status code per line — a leading space followed
by `M` marks a tracked file with unstaged modifications, which is what
gets counted; untracked files show `??` instead and are correctly left
out of that count.
