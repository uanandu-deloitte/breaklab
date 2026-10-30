# Two ways a merge can go

`git merge <branch>` always does the same thing conceptually — bring
`<branch>`'s commits into the branch you're on — but it produces one of
two very different results depending on history's shape.

## Fast-forward: no new commit

If the branch you're on hasn't moved since the other branch's starting
point, merging it in is just a matter of sliding your branch's pointer
forward to match. No new commit is created — the history afterward looks
exactly like it would have if you'd made every commit directly on the
branch you're on. `git log` shows a straight line.

## True merge: a new commit with two parents

If *both* branches have moved forward independently since they diverged,
there's no single pointer move that captures both histories — git has to
create a new commit that has two parents: the tip of each branch being
combined. This is a real, permanent record in history that a merge
happened at this point, combining two independent lines of work.

## Reading the graph before merging

```
git log --graph --oneline --all
```

`--graph` draws the commit DAG with ASCII art, `--oneline` keeps each
commit to one line, and `--all` shows every branch at once instead of
just the one you're on. Before merging anything, this is how you confirm
whether you're looking at a situation that can fast-forward or one that
can't — the shape of the graph (a single line vs. two diverging lines
that need to reconverge) tells you which one you're about to get.

## Forcing a true merge commit on purpose

Even when a fast-forward *would* be possible, `git merge --no-ff <branch>`
forces git to create a merge commit anyway — useful when a team wants
every feature branch's integration point recorded explicitly in history,
rather than disappearing into a straight line.
