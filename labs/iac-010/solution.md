# Solution: Renaming a Resource Isn't Free

```bash
cd /infra
terraform plan
# shows: local_file.legal_notice will be destroyed
#        local_file.privacy_notice will be created
# -- the rename, read naively, looks like a destroy+recreate.

terraform state mv local_file.legal_notice local_file.privacy_notice
# Move succeeded.

terraform plan -detailed-exitcode
# No changes. Your infrastructure matches the configuration. (exit 0)
```

`/var/www/legal/notice.txt` keeps the exact same mtime throughout —
`terraform state mv` only ever edits the state file, it never touches
the real resource.
