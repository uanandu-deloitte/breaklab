# INC-5102: Deploys hang waiting for this container to stop

`docker stop` on this service never exits cleanly on its own — it
always ends up force-killed after the grace period. Whatever is
supposed to handle shutdown in here isn't reacting to it.

Get a shell: docker exec -it breaklab-swallowed-signal bash
