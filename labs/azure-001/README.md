# main.bicep Won't Compile

`/infra/main.bicep` declares the storage account the checkout team
wants for their new export pipeline. Nobody can deploy it — it doesn't
even compile.

```
az bicep build --file /infra/main.bicep --stdout
```

Read the real compiler error (not a guess, the actual `az bicep build`
output), fix `main.bicep`, and get it compiling — without just changing
`retentionDays`'s declared type to make the error go away.

Get a shell: `docker exec -it breaklab-azure-001 bash`
