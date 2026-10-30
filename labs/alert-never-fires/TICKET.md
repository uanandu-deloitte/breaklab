# INC-970: Disk nearly full, nobody got paged

/var/lib/metrics/current.json has shown disk_used_pct at 97 for a while
now, but nobody has been alerted. Get the alerting pipeline actually
catching this.

Get a shell: docker exec -it breaklab-alert-never-fires bash
