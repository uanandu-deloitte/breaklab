# INC-903: Last night's backup is stale

The nightly backup for `opsuser` is supposed to pick up the day's data
automatically. This morning's backup doesn't match what's actually in
/var/app/data. Get it current, and make sure tonight's run will work too.

Get a shell: docker exec -it breaklab-cron-blindspot bash
