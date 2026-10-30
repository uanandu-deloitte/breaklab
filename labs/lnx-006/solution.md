# Solution: Repair the Application Layout

```bash
mv /opt/acme-app/config/acme-app /opt/acme-app/bin/acme-app
chmod +x /opt/acme-app/bin/acme-app

mv /opt/acme-app/data/app.conf /opt/acme-app/config/app.conf

mv /opt/acme-app/bin/data.db /opt/acme-app/data/data.db
```

Each `mv` relocates one file to the directory its content actually
belongs in — the shebang-and-echo script to `bin/`, the `key: value`
config to `config/`, the application data file to `data/` — without
touching its content. `chmod +x` restores the executable bit on the
script now that it's in `bin/`, since a plain `mv` doesn't guarantee
that bit survived being dropped in `config/` as a non-executable copy.
