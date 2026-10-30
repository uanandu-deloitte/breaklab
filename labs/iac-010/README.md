# Renaming a Resource Isn't Free

Someone cleaned up `/infra/main.tf` and renamed the resource that manages
`/var/www/legal/notice.txt` from `legal_notice` to `privacy_notice` — same
file, same content, just a clearer name. They didn't touch Terraform's
state when they did it.

Find out what Terraform actually plans to do about that rename before
you let it do anything. The real file (and whatever's actually serving
it) must never go down or get recreated just because a label in a `.tf`
file changed.

Get a shell: `docker exec -it breaklab-iac-010 bash`
