# The Blob Nobody Outside The Team Can Read

`public-assets/notice.txt` is supposed to be a public, anonymously
readable status notice — the kind of blob a status page or CDN would
fetch directly with no credentials at all. Right now it's real (you
can see it fine with the real connection string), but a plain anonymous
HTTP GET against it fails.

This is a real Azurite instance (Azure's official local Storage
emulator) and the real `az storage` CLI — the same commands and the
same behavior as real Azure Blob Storage, just running on this box with
no cloud account.

```
az storage container show --name public-assets
curl -i http://127.0.0.1:10000/devstoreaccount1/public-assets/notice.txt
```

Figure out why the anonymous request fails and fix the real container
configuration so it doesn't, without handing out the account key or a
SAS token to get there.

Get a shell: `docker exec -it breaklab-azure-002 bash`
