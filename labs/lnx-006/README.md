# Repair the Application Layout

A failed deployment copied acme-app's files into the wrong directories
under `/opt/acme-app`. The layout should be:

```
/opt/acme-app/
├── bin/     — the application executable
├── config/  — the configuration file
├── data/    — application data
└── logs/
```

Right now, three files are sitting in the wrong subdirectory. Nothing
has been deleted — they're just misplaced.

## Tasks

- Inspect every file under `/opt/acme-app` and work out, from its name
  and content, which directory it actually belongs in.
- Move each misplaced file to its correct directory.
- The executable must be runnable (`+x`) once it's in `bin/`.

## Constraints

- Don't delete any file.
- Don't replace a file with an empty placeholder — every file's content
  must be exactly what it was before you moved it.

Get a shell: `docker exec -it breaklab-lnx-006 bash`
