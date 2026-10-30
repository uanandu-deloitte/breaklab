# The config file is the truth; the real resource is just its shadow

An infra-as-code apply tool works backwards from how you might edit a
file by hand: you don't change the thing you want changed — you change a
*declaration* of what it should look like, then hand that declaration to
a tool that makes the real resource match it.

In this container, `/infra/resources.yaml` declares what
`/etc/banner.txt` should contain. `/usr/local/bin/infra-apply` reads the
declaration and overwrites the real file to match it. It already ran
once, at image build time — that's why `/etc/banner.txt` already has
content before you've touched anything.

## Which one is the source of truth?

`/infra/resources.yaml` is. `/etc/banner.txt` is not something you
maintain directly — it's a side effect, a cache of whatever
`resources.yaml` last said. If you edit `/etc/banner.txt` by hand, two
things are true at once:

- The file looks changed, right now.
- The *next* time anyone runs `infra-apply`, your manual edit is gone —
  overwritten back to whatever `resources.yaml` declares, because the
  tool's whole job is to make reality match the declaration, not to
  preserve whatever reality happened to be.

## Why this matters

Treating the declaration as the single source of truth (and the real
resource as disposable, regenerable output) is the core idea behind every
"infrastructure as code" tool, however small or large. The lesson that
starts here: never "fix" a managed resource by editing it directly —
fix the declaration, then re-apply.
