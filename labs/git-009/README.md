# Push Rejected

You've been working in `/workspace`, a local clone of `/srv/origin.git`.
You made a commit of your own. When you try to push it, git refuses.

## Task

1. Figure out why `git push` is being rejected.
2. Bring your local `main` and origin's `main` together without losing
   either commit -- not yours, and not the one a teammate already pushed
   while you were working.
3. Push successfully. Origin's history should end up containing both
   commits.

## Constraints

- Don't force-push. A force-push would make this "succeed" by throwing
  away the teammate's commit -- that's the mistake this lab exists to
  catch, not the fix.

Get a shell: `docker exec -it breaklab-git-009 bash`
