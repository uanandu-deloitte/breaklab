# Dashboards Query by Name — Renames Break Silently

A dashboard panel is rendering "no data." The exporter it's supposed to
read from is up and reporting real numbers. Find out exactly which name
the panel is asking for, which name the exporter is actually using, and
why those two things are no longer the same.

## Tasks

1. Read `/etc/dashboards/overview.yaml` — this is the panel's config,
   and it names exactly one metric to look up.
2. Read the keys in `/var/lib/metrics/exported.json` with `jq 'keys'` —
   this is everything the exporter is actually emitting right now.
3. Run `render-dashboard` and read `/var/lib/dashboards/overview.rendered`
   to confirm the panel really is getting no data right now, live.
4. Find the exporter key that's the renamed version of the metric the
   panel is asking for — compare the two names exactly, not just by eye.

Record your findings in `/workspace/answers.env`, one `KEY=value` per line:

```
PANEL_METRIC=
EXPORTED_METRIC=
RENDER_RESULT=
```

- `PANEL_METRIC`: the exact metric name `/etc/dashboards/overview.yaml`
  is configured to look up.
- `EXPORTED_METRIC`: the exact key, actually present in
  `/var/lib/metrics/exported.json`, that is the renamed successor of
  `PANEL_METRIC` — the same underlying metric, new name.
- `RENDER_RESULT`: whatever `render-dashboard` actually writes to
  `/var/lib/dashboards/overview.rendered` right now (it should show the
  lookup is failing, not succeeding).

Get a shell: `docker exec -it breaklab-obs-006 bash`
