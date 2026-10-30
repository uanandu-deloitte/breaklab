# Who Owns the Port?

`curl http://localhost:8080/` is returning nothing useful — the web
server this host is supposed to run doesn't seem to be serving its real
page.

## Tasks

- identify exactly what process is currently listening on port 80 inside
  the container
- stop *only* that process — don't touch anything else
- start the real web server (nginx) so it can bind port 80
- confirm the real page is now served

Get a shell: `docker exec -it breaklab-net-002 bash`

Check from the host: `curl http://localhost:8080/`
