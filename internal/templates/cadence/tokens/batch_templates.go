package tokens

// CreateAccountAndSetupTemplate creates account and sets up multiple token vaults in one transaction
const CreateAccountAndSetupTemplate = `
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

// AddFungibleTokenVaultBatchTemplate adds multiple token vaults to an existing account
const AddFungibleTokenVaultBatchTemplate = `
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