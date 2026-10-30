# INC-614: New jobs stop completing

Jobs that were completing fine this morning have stopped. `free -h`
and CPU usage both look completely normal — nothing is starved for
memory or compute. New jobs just never finish.

Get a shell: docker exec -it breaklab-pids-ceiling bash
