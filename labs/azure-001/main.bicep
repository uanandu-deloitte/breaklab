param storageAccountName string = 'acmecheckoutdata'
param location string = 'eastus'
// BUG: retentionDays is typed int, but its default value is the string
// literal '30' instead of the number 30.
param retentionDays int = '30'

resource storageAccount 'Microsoft.Storage/storageAccounts@2023-01-01' = {
  name: storageAccountName
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'
  tags: {
    retentionDays: string(retentionDays)
  }
  properties: {
    accessTier: 'Hot'
  }
}
