# INC-452: Config change didn't take effect

Someone updated app-config's message from v1 to v2 this morning, but
config-reader is still acting on the old value. Get it actually picking
up v2.

Get a shell: docker exec -it breaklab-stale-configmap bash
