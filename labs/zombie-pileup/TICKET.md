# INC-414: The box keeps getting slower

Nothing crashed and CPU/memory both look fine, but the process table is
quietly filling up with `<defunct>` entries and it's getting worse by
the minute. Find what's leaking and stop it for good.

Run `restart-workers` once you think it's fixed.

Get a shell: docker exec -it breaklab-zombie-pileup bash
