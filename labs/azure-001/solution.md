# Solution: main.bicep Won't Compile

```bicep
param retentionDays int = 30
```

```bash
az bicep build --file /infra/main.bicep --stdout | jq '.parameters.retentionDays'
# {
#   "type": "int",
#   "defaultValue": 30
# }
```
