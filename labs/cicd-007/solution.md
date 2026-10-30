# Solution: Promoting Artifacts Between Environments

```bash
ls -la /srv/artifacts/
# app-2.1.0.tar  app-2.1.5.tar  app-2.2.0.tar  current -> app-2.1.0.tar

readlink /srv/artifacts/current
# app-2.1.0.tar

ls -la --time-style=full-iso /srv/artifacts/*.tar
# app-2.2.0.tar has the most recent mtime; app-2.1.0.tar has the oldest.

cat <<EOF > /workspace/answers.env
PROMOTED_ARTIFACT=app-2.1.0.tar
WOULD_GUESSING_BY_NAME_BE_RIGHT=no
EOF
```

`current` is a symlink, and `readlink` shows it genuinely points at
`app-2.1.0.tar` — the *oldest* build by both version number and mtime,
not the newest. `app-2.2.0.tar` has the highest version number and was
built most recently, which is exactly why guessing by filename or
recency would have gotten this wrong: both point at a build that was
never promoted.

This is a realistic shape for a deliberate rollback or a pinned
release: production is intentionally running an older, known-good
build while newer ones exist elsewhere in the same directory. The
symlink is the one piece of state whose entire job is recording what's
actually live — not the filenames sitting next to it.
