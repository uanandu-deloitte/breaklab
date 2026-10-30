# Build the Application Layout

Create the following directory structure:

```
/opt/acme-app/
├── bin/
├── config/
├── data/
└── logs/
```

Then:

- add a placeholder configuration file at `/opt/acme-app/config/app.conf`
- add a placeholder log file at `/opt/acme-app/logs/app.log`
- copy (don't move) the supplied executable at `/opt/fixtures/acme-app` into
  `/opt/acme-app/bin/acme-app`, keeping it executable

The original file at `/opt/fixtures/acme-app` must still exist afterward.

Get a shell: `docker exec -it breaklab-lnx-005 bash`
