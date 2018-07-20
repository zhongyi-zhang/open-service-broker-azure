package mssqlfg

// nolint: lll
var dbmsFeARMTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"tags": {
			"type": "object"
		}
	},
	"resources": [
		{
			"type": "Microsoft.Sql/servers",
			"name": "{{ .serverName }}",
			"apiVersion": "2015-05-01-preview",
			"location": "{{.location}}",
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
		"fullyQualifiedDomainName": {
			"type": "string",
			"value": "[reference('{{ .serverName}}').fullyQualifiedDomainName]"
		}
	}
}
`)
