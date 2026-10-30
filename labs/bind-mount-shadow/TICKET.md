# INC-430: Nothing is showing up in /data

The writer process is supposed to update /data/heartbeat.txt every
second. /data has been empty all morning. The Dockerfile clearly
chowns /data to appuser at build time, so permissions "should" be
fine — but something's still blocking the writes.

Get a shell: docker exec -it breaklab-bind-mount-shadow bash
