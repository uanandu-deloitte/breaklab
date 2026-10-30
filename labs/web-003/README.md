# INC: 502 Bad Gateway

Customers hitting http://localhost:8080 get a 502 Bad Gateway. nginx
itself is up — this is nginx telling you it couldn't reach whatever
it's supposed to forward the request to.

## Tasks

- read nginx's error log to see exactly what address it tried and
  failed to reach
- confirm what port the real backend is actually listening on
- compare that against nginx's configured `proxy_pass` address
- fix the mismatch and reload nginx — without bypassing it

## Constraints

- the fix must go through nginx: the exposed port (8080) must serve the
  backend's content, not the backend's own port directly
- don't move the backend to a different port to "match" nginx — nginx's
  config is what's wrong here

Get a shell: `docker exec -it breaklab-web-003 bash`

Check it from outside the container: `curl -i http://localhost:8080/`
