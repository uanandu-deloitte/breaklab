# The documentation is already on the host

You won't always have internet access, a cheat sheet, or a command you've
memorized the exact flags for. Every standard Linux command carries its
own documentation with it, reachable with nothing more than the command
itself.

## Two ways to ask a command about itself

- `<command> --help` prints a short usage summary and a list of options,
  straight to your terminal — fast, and almost always available.
- `man <command>` opens the full manual page: longer, more thorough, and
  searchable with `/` once you're inside it (press `q` to leave).

Both describe the real options the installed version of the command
supports — more reliable than recalling a flag from memory or from a
different system where the same command behaved slightly differently.

## Why this matters

Restricted or unfamiliar hosts are exactly where guessed syntax is most
likely to be wrong — a different distribution, a different version, or a
minimal install can all change which flags exist. Checking first costs a
few seconds; assuming costs you a failed command, or worse, a command
that silently does something other than what you intended.

## Reading an options list

`--help` output lists each option with a one-line description, usually
grouped as short form (`-a`) and long form (`--all`) together. Scan for
the words that match what you're trying to do — "hidden", "parents",
"interactive" — rather than reading every line top to bottom.
