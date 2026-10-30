# Solution: Build the Application Layout

```bash
mkdir -p /opt/acme-app/{bin,config,data,logs}
touch /opt/acme-app/config/app.conf
touch /opt/acme-app/logs/app.log
cp /opt/fixtures/acme-app /opt/acme-app/bin/acme-app
chmod +x /opt/acme-app/bin/acme-app
```

`mkdir -p` with brace expansion builds all four directories in one call.
`cp` (not `mv`) leaves `/opt/fixtures/acme-app` untouched, and the explicit
`chmod +x` guarantees the copy stays executable regardless of what `cp`
preserved on its own.
