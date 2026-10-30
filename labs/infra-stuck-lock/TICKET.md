# INC-430: This morning's config change never landed

Ops pushed an update to /infra/site.yaml over an hour ago, but
http://localhost:8082 still shows the old message. Every attempt to
run infra-apply since then has failed outright. Get it applying again.

Get a shell: docker exec -it breaklab-infra-stuck-lock bash
