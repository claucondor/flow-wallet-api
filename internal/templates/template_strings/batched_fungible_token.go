package template_strings

import (
	"bytes"
	"text/template"
)

type BatchedFungibleOpsInfo struct {
	FungibleTokenContractAddress string
	Tokens                       []FungibleTokenInfo
}

type FungibleTokenInfo struct {
	ContractName       string
	Address            string
	VaultStoragePath   string
	ReceiverPublicPath string
	BalancePublicPath  string
}

func AddFungibleTokenVaultBatchTransaction(i BatchedFungibleOpsInfo) (string, error) {
	return executeTemplate("AddFungibleTokens", AddFungibleTokenVaultBatchTransactionTemplate, i)
}

func CreateAccountAndSetupTransaction(i BatchedFungibleOpsInfo) (string, error) {
	return executeTemplate("CreateAccount", CreateAccountAndSetupTransactionTemplate, i)
}

const CreateAccountAndSetupTransactionTemplate = `
import Crypto
import FungibleToken from {{ .FungibleTokenContractAddress }}
{{ range .Tokens }}
import {{ .ContractName }} from {{ .Address }}
{{ end }}

transaction(publicKeys: [Crypto.KeyListEntry]) {
	prepare(signer: auth(CreateAccount, Storage, Capabilities) &Account) {
		let account = Account(payer: signer)

		// add all the keys to the account
		for key in publicKeys {
			account.keys.add(publicKey: key.publicKey, hashAlgorithm: key.hashAlgorithm, weight: key.weight)
		}

		{{ range .Tokens }}
		// initializing vault for {{ .ContractName }}
		account.storage.save(<-{{ .ContractName }}.createEmptyVault(), to: {{ .VaultStoragePath }})
		
		let receiverCap = account.capabilities.storage.issue<&{{ .ContractName }}.Vault{FungibleToken.Receiver}>(
			{{ .VaultStoragePath }}
		)
		account.capabilities.publish(receiverCap, at: {{ .ReceiverPublicPath }})
		
		let balanceCap = account.capabilities.storage.issue<&{{ .ContractName }}.Vault{FungibleToken.Balance}>(
			{{ .VaultStoragePath }}
		)
		account.capabilities.publish(balanceCap, at: {{ .BalancePublicPath }})
		{{ end }}
	}
}
`

const AddFungibleTokenVaultBatchTransactionTemplate = `
import FungibleToken from {{ .FungibleTokenContractAddress }}
{{ range .Tokens }}
import {{ .ContractName }} from {{ .Address }}
{{ end }}

transaction() {
	prepare(account: auth(Storage, Capabilities) &Account) {
		{{ range .Tokens }}
		// initializing vault for {{ .ContractName }}
		if account.storage.borrow<&{{ .ContractName }}.Vault>(from: {{ .VaultStoragePath }}) == nil {
			account.storage.save(<-{{ .ContractName }}.createEmptyVault(), to: {{ .VaultStoragePath }})
			
			let receiverCap = account.capabilities.storage.issue<&{{ .ContractName }}.Vault{FungibleToken.Receiver}>(
				{{ .VaultStoragePath }}
			)
			account.capabilities.publish(receiverCap, at: {{ .ReceiverPublicPath }})
			
			let balanceCap = account.capabilities.storage.issue<&{{ .ContractName }}.Vault{FungibleToken.Balance}>(
				{{ .VaultStoragePath }}
			)
			account.capabilities.publish(balanceCap, at: {{ .BalancePublicPath }})
		}
		{{ end }}
	}
}
`

func executeTemplate(name string, temp string, i BatchedFungibleOpsInfo) (string, error) {
	template, err := template.
		New(name).
		Parse(temp)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = template.Execute(buf, i)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
