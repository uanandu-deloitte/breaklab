# Solution: The Blob Nobody Outside The Team Can Read

```bash
az storage container show --name public-assets --query properties.publicAccess -o tsv
# off

az storage container set-permission --name public-assets --public-access blob

curl -i http://127.0.0.1:10000/devstoreaccount1/public-assets/notice.txt
# HTTP/1.1 200 OK
# ...
# Checkout is back online. Thanks for your patience.
```
