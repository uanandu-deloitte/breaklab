# INC-6033: docker logs is empty for a service that's supposedly running

Nobody can see anything in `docker logs` for this container, but the
on-call engineer who checked earlier swears the process is alive and
doing work.

Get a shell: docker exec -it breaklab-silent-stdout bash
