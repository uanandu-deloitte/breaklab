# Nothing Listens Where You Think

A tiny status-page server runs inside this container, with host port
`8080` mapped to container port `80`. From the host, `curl
http://localhost:8080/` returns nothing. The process hasn't crashed, and
nothing else is obviously wrong.

## Tasks

1. Run `curl` from the host against the mapped port, then get a shell
   (`docker exec -it breaklab-docker-006 bash`) and run `curl localhost`
   from inside the container. Compare the two results.
2. Use `ss -tlnp` inside the container to see exactly which address the
   server process is bound to, not just which port.
3. Fix the server so the mapped host port actually reaches it, without
   changing which port it listens on.

Get a shell: `docker exec -it breaklab-docker-006 bash`
