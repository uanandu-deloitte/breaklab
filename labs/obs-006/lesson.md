# A dashboard panel is just a name lookup

A panel doesn't know anything about your system beyond one string: the
metric name its config says to render. At render time, that name gets
looked up against whatever the metrics exporter is currently emitting —
an exact string match, nothing fuzzy. If that lookup succeeds, the panel
shows a value. If it fails, the panel shows "no data." There is no third
outcome.

## Where renames come from

Metric names drift for ordinary reasons: a unit gets clarified
(`disk_free_pct` becomes `disk_free_percent`), a naming convention gets
standardized across services, an exporter library gets upgraded and
changes its default naming. Every one of these is a reasonable,
deliberate change on the exporter side. The problem is never the rename
itself — it's that nothing forces every dashboard config that referenced
the old name to update in lockstep with it.

## Why this is quiet instead of loud

The render step treats a missing key the same way a well-behaved program
treats any other expected-to-sometimes-be-absent value: it substitutes a
placeholder ("no-data") and moves on. That's the right behavior for
genuinely missing data (a metric that hasn't reported yet). It's the
wrong signal for this situation, but the render step has no way to tell
the two apart — both look identical from where it's standing: a config
asking for a key, and that key not being in the data.

## How to diagnose it

- Read the panel's config to get the exact name it's asking for.
- List the exporter's actual keys (`jq keys` against its JSON output,
  or the equivalent for your metrics system) and compare them to that
  name character by character — renames are usually close, not wildly
  different, so a quick skim will fool you.
- Re-run the render step directly and look at its real output, rather
  than trusting the dashboard UI's cached or summarized view.

## The fix

Either point the panel at the new name directly, or introduce an
alias/compatibility mapping from the old name to the new one — useful
when other dashboards still reference the old name and can't all be
migrated at once.
