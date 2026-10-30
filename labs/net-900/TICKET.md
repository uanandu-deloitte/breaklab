# INC-900: The public API is down

Customers report that the public API (http://localhost:8080) is
returning errors since a change window earlier today. Restore service.

Architecture: `client -> nginx -> api -> database`, each hop a
separate process in this container.

## Constraints

- do not replace the api process's code or the database's code
- do not weaken the api service account's permissions beyond what it
  actually needs

Get a shell: `docker exec -it breaklab-net-900 bash`
