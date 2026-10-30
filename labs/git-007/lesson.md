# When a merge can't resolve itself

A merge conflict happens when the same region of the same file was
changed differently on both sides being merged — git has no way to know
which version (or what combination) you actually want, so it stops and
hands the decision to you instead of guessing.

## What a conflict looks like

Git leaves both versions directly in the file, wrapped in markers:

```
<<<<<<< HEAD
timeout: 45
=======
timeout: 60
>>>>>>> feature/raise-timeout
```

- Everything between `<<<<<<< HEAD` and `=======` is what's currently on
  the branch you're merging *into*.
- Everything between `=======` and `>>>>>>> <branch>` is what's coming in
  from the branch you're merging *from*.
- The file is left in a half-finished state on disk — it won't compile
  or parse correctly as-is. Resolving the conflict means editing it down
  to exactly the content you want, markers and all removed.

## Resolving and finishing the merge

1. Edit the file so it contains the correct final content — not one side
   or the other by default, but whatever the right answer actually is.
   Delete all three marker lines; nothing with `<<<<<<<`, `=======`, or
   `>>>>>>>` should remain.
2. `git add <file>` tells git you've resolved that file's conflict.
3. `git commit` (with no `-m` needed — git already prepared a merge
   commit message) finishes the merge, creating a commit with two
   parents just like a non-conflicting merge would.

## Why the "right" resolution matters

Picking a side at random isn't resolving a conflict, it's discarding
someone's work and hoping it didn't matter. A real resolution reflects
an actual decision about what the correct final state is — here, that
decision has already been made for you (see the task), so the exercise
is purely in applying it correctly.
