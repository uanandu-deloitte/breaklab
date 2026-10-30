# Normalize Deployment Data

`/workspace/hosts.txt` is a deployment inventory someone pasted together
by hand: mixed case, duplicate hosts, and inconsistent delimiters — some
entries are comma-separated, some are one per line, and several have
stray leading/trailing whitespace or a tab.

Produce a clean, canonical host list at `/workspace/normalized.txt`:

- every hostname lowercase
- one hostname per line
- no duplicates
- sorted alphabetically
- no leading/trailing whitespace on any line

## Constraints

- Build the result with a pipeline (`tr`, `cut`/`sed`, `sort`, `uniq`) —
  no manually editing or retyping the host list.
- Don't modify `hosts.txt`.

Get a shell: `docker exec -it breaklab-lnx-011 bash`
