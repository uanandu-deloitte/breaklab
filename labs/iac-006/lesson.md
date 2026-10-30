# Only one source of config ever wins

A lot of real tools let you set the same setting three different ways:
a config file, an environment variable, and a command-line flag. That's
convenient — until two of them disagree, and you discover the tool
never actually averages them or asks which one you meant. Exactly one
precedence order is baked into the tool, and it decides the outcome
every time, whether or not anyone remembers that order exists.

A common order — and the one this lab's `infra-apply` uses — is:

```
flag > environment variable > config file
```

A flag passed explicitly on this specific invocation wins over
everything. Failing that, an environment variable set anywhere in the
current process's environment wins over the file. The file is only
ever consulted as the last resort, when nothing more specific was
provided.

## Why "the file is correct" isn't the whole story

It's tempting to treat a declared config file as automatically
authoritative — it's checked in, it's reviewed, it's the thing everyone
agrees should be true. But a tool that also reads an environment
variable doesn't know any of that. If `$INFRA_VALUE` happens to be set
— baked into a container image, exported in a CI job, left over in
someone's shell profile — the tool will use it instead, silently,
every single run, no matter how correct the file is. The file isn't
wrong. It's just not the thing actually being honored.

This is exactly why "I checked the file, it's correct" is not the same
claim as "I checked what the tool actually resolved to." The first is
about intent; the second is about what really happened. Only running
the tool and reading its own report of which source it used (or
checking the real resource it wrote) answers the second question.

## Why this matters

Config precedence bugs are invisible by construction: nothing is
malformed, nothing errors, the file people are looking at is genuinely
correct. The actual cause — a stray environment variable nobody
remembered was set — lives somewhere else entirely. Knowing the
precedence order in advance, and knowing to check for it, turns a
"this makes no sense" incident into a two-command diagnosis.
