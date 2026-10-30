# Find It Yourself

You have shell access to a restricted host. Don't rely on memorized
syntax — use `--help` or `man` to discover the options you actually need,
then use them.

## Tasks

1. Find the `ls` option that lists hidden files (including dotfiles).
2. Find the `mkdir` option that creates a nested directory tree in a
   single command, even when the parent directories don't exist yet.
3. Find the `cp` option that makes `cp` ask for confirmation before it
   overwrites an existing file.

Record what you found in `/workspace/answers.env`, one `KEY=value` per
line:

```
HIDDEN_FILES_OPTION=
NESTED_MKDIR_OPTION=
INTERACTIVE_COPY_OPTION=
```

Then actually use each option:

- create a hidden file at `/workspace/.hidden-check`
- create the nested directory tree `/workspace/deep/nested/tree` with a
  single `mkdir` command
- copy `/opt/fixtures/readme.txt` to `/workspace/fixture-copy.txt` using
  the interactive copy option you discovered

Get a shell: `docker exec -it breaklab-lnx-002 bash`
