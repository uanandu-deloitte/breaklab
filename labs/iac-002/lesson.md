# See the diff before you cause it

A real apply tool rewrites a managed resource the moment you run it —
which is fine when you're sure what it's about to do, and risky when
you're not. Most real infra-as-code tools solve this with a **plan** (or
"dry run") mode: run the tool with a flag, and instead of changing
anything, it tells you exactly what it *would* change, so you can read
that description before committing to it.

In this container, `infra-apply --plan` does exactly that:

- It reads `/infra/resources.yaml` (the declared, desired content) and
  the real file's current content.
- If they already match, it prints `no changes`.
- If they don't, it prints a line describing the exact before → after
  transition it would make — without writing anything.

Only `infra-apply` with no flag actually writes the new content.

## Why this matters

A plan step turns "run it and hope" into "read it, then decide." It lets
you catch a declaration that's wrong, or a change that's bigger (or
different) than you expected, before it ever touches the real resource —
and it's a genuinely safe, repeatable thing to run as many times as you
want, because by design it never writes.

The discipline this lab is building: when a tool offers a plan or dry-run
mode, run that first, read its output, and only then run the real apply.
