# INC-4417: Consumer process can't read the producer's export

run-producer.sh writes /var/lib/appdata/export.json successfully every
time, and both it and the reader are supposed to be in the same
appdata group. But the reader keeps getting permission denied on the
file producer just wrote.

Get a shell: docker exec -it breaklab-umask-mismatch bash
