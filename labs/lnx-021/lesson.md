# Finding what's actually using the disk

`df` tells you a filesystem is full. It never tells you why. Answering
"why" means walking down from a filesystem, to a directory, to the
specific file or subdirectory actually responsible — and `du` and `find`
are the tools that do that walk.

## Comparing sibling directories with du

`du -sh /data/*/` prints one total size per top-level directory matched
by the glob, each already summed recursively (`-s`) and shown in a
human-readable unit (`-h`). Piped through `sort -rh`, the largest
directory sorts to the top:

```bash
du -sh /data/*/ | sort -rh | head -1
```

Without `-s`, `du` prints a line for every single subdirectory it walks
through — useful for drilling in further, but noisy when you just want
one total per top-level directory.

## Finding the single largest file

A directory's total size doesn't tell you whether it's one huge file or
a thousand small ones. `find` combined with its `-printf` action can
print every file's size next to its path, which `sort -rn` can then rank:

```bash
find /data -type f -printf '%s %p\n' | sort -rn | head -1
```

`%s` is the file's size in bytes, `%p` is its path — sorting numerically
descending (`-rn`) on that output puts the biggest file first, regardless
of which subdirectory it happens to live in.

## Why "compute it live" matters

A size recorded once goes stale the moment anything changes — a log
rotates, a cache evicts, a new file lands. `du` and `find` re-derive the
answer from whatever is actually on disk right now, which is the only
version of "where did the space go" worth trusting.
