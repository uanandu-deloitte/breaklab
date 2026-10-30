# ';' runs the next command regardless; '&&' doesn't

A shell line like `a; b; c` runs `a`, then `b`, then `c` — in order, but
completely independently of whether each one succeeded. `;` is just a
command separator; it carries no information about exit codes at all. If
`a` fails, `b` and `c` still run exactly as if nothing had gone wrong.

`a && b && c` is different: `&&` only runs the next command if the one
before it exited `0`. If `a` fails, `b` never runs, and `c` never runs
either — the chain stops at the first failure.

## Why this matters for a pipeline

A driver script that chains build, test, and deploy stages with `;`
treats a failing test suite exactly the same as a passing one — it moves
on to deploy regardless. The fix isn't to change what the test stage
does; it's to change how the driver reacts to that stage's exit code.

Two ways to get this right:

- Chain the stage invocations with `&&` instead of `;`, so a non-zero
  exit from any stage skips everything after it.
- Add `set -e` near the top of the script, which makes the whole script
  exit immediately the first time any command in it exits non-zero —
  without needing to touch how the stages are separated at all.

Either one turns "ran every stage no matter what" into "stopped at the
first real failure," which is the entire point of having stages in the
first place.
