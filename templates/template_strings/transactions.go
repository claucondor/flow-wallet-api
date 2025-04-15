package template_strings

const AddAccountContractWithAdmin = `
import FungibleToken from 0x9a0766d93b6608b7
import FlowToken from 0x7e60df042a9c0868
import FlowIDTableStaking from 0x8624b52f9ddcd04a
import LockedTokens from 0x8d0e87b65159ae63
import FlowStakingCollection from 0x8d0e87b65159ae63
import NodeOperatorAccounts from 0x8d0e87b65159ae63

transaction(adminAccount: Address, name: String, code: String, publicKeys: [String]) {
	prepare(signer: auth(Storage, Contracts) &Account) {
		let acct = Account.create()

		for key in publicKeys {
			let publicKey = PublicKey(
				publicKey: key.decodeHex(),
				signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
			)
			
			acct.keys.add(
				publicKey: publicKey,
				hashAlgorithm: HashAlgorithm.SHA3_256,
				weight: 1000.0
			)
		}

		let admin = getAccount(adminAccount)

		let adminRef = admin.storage.borrow<auth(FungibleToken.Withdraw) &FlowToken.Vault>(from: /storage/flowTokenVault)
			?? panic("Could not borrow reference to admin token vault")

		let initialFundingAmount: UFix64 = 0.001

		let fundingProvider = adminRef.withdraw(amount: initialFundingAmount)

		let tokenReceiver = acct.capabilities.get<auth(FungibleToken.Deposit) &{FungibleToken.Receiver}>(/public/flowTokenReceiver)!

		tokenReceiver.borrow()!.deposit(from: <-fundingProvider)

		acct.contracts.add(name: name, code: code.decodeHex())
	}
}
`

const CreateAccount = `
transaction(publicKeys: [String]) {
	prepare(signer: auth(Storage, Keys) &Account) {
		let acct = Account.create()

		for key in publicKeys {
			let publicKey = PublicKey(
				publicKey: key.decodeHex(),
				signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
			)
			
			acct.keys.add(
				publicKey: publicKey,
				hashAlgorithm: HashAlgorithm.SHA3_256,
				weight: 1000.0
			)
		}
	}
}
`

const GenericFungibleTransfer = `
import FungibleToken from 0x{{.FungibleTokenAddress}}

transaction(amount: UFix64, to: Address) {
	let vault: &{FungibleToken.Vault}
	let receiver: &{FungibleToken.Receiver}

	prepare(signer: auth(Storage) &Account) {
		self.vault = signer.storage.borrow<&{FungibleToken.Vault}>(/storage/flowTokenVault)
			?? panic("Could not borrow provider vault")

		self.receiver = getAccount(to).capabilities.get<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)
			.borrow() ?? panic("Could not borrow receiver vault")
	}

	execute {
		self.receiver.deposit(from: <-self.vault.withdraw(amount: amount))
	}
}
`

const GenericFungibleSetup = `
import FungibleToken from 0x{{.FungibleTokenAddress}}
import {{.TokenContractName}} from 0x{{.TokenAddress}}

transaction {
	prepare(signer: auth(Storage, Keys) &Account) {
		// Configure storage for {{.TokenContractName}} if it doesn't already exist
		if signer.storage.borrow<&{{.TokenContractName}}.Vault>(/storage/{{.TokenStoragePath}}) == nil {
			// Create a new Vault and put it in storage
			signer.storage.save(
				<-{{.TokenContractName}}.createEmptyVault(),
				to: /storage/{{.TokenStoragePath}}
			)

			// Create a public capability to the stored Vault that exposes
			// the deposit function through the Receiver interface
			signer.capabilities.publish(
				signer.capabilities.storage.issue<&{FungibleToken.Receiver}>(/storage/{{.TokenStoragePath}}),
				at: /public/{{.TokenPublicReceiverPath}}
			)

			// Create a public capability to the stored Vault that exposes
			// the balance field through the Balance interface
			signer.capabilities.publish(
				signer.capabilities.storage.issue<&{FungibleToken.Balance}>(/storage/{{.TokenStoragePath}}),
				at: /public/{{.TokenPublicBalancePath}}
			)
		}
	}
}
`

const AddProposalKeyTransaction = `
transaction(adminKeyIndex: Int, numProposalKeys: UInt16) {
  prepare(account: auth(Storage, Keys) &Account) {
    let key = account.keys.get(keyIndex: adminKeyIndex)!
    var count: UInt16 = 0
    while count < numProposalKeys {
      account.keys.add(
            publicKey: key.publicKey,
            hashAlgorithm: key.hashAlgorithm,
            weight: 0.0
        )
        count = count + 1
    }
  }
}
`

const AddAccountKeysTransaction = `
transaction(publicKeys: [String]) {
  prepare(signer: auth(Storage, Keys) &Account) {
    for pbk in publicKeys {
      let key = PublicKey(
        publicKey: pbk.decodeHex(),
        signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
      )

      signer.keys.add(
        publicKey: key,
        hashAlgorithm: HashAlgorithm.SHA3_256,
        weight: 1000.0
      )
    }
  }
}
`

const GenericFungibleTransferMemo = `
import FungibleToken from 0x{{.FungibleTokenAddress}}
import {{.TokenContractName}} from 0x{{.TokenAddress}}

transaction(amount: UFix64, to: Address, memo: String) {
	let vault: &{FungibleToken.Vault}
	let receiver: &{FungibleToken.Receiver}

	prepare(signer: auth(Storage) &Account) {
		self.vault = signer.storage.borrow<&{FungibleToken.Vault}>(/storage/{{.TokenStoragePath}})
			?? panic("Could not borrow provider vault")

		self.receiver = getAccount(to).capabilities.get<&{FungibleToken.Receiver}>(/public/{{.TokenPublicReceiverPath}})
			.borrow() ?? panic("Could not borrow receiver vault")
	}

	execute {
		self.receiver.deposit(from: <-self.vault.withdraw(amount: amount))
	}
}
`

const GenericNFTTransfer = `
import NonFungibleToken from 0x631e88ae7f1d7c20
import {{.TokenContractAddress}} from {{.TokenContractAddress}}

transaction(id: UInt64, recipient: Address) {
    prepare(signer: auth(Storage) &Account) {
        
        // get the recipients public account object
        let recipient = getAccount(recipient)

        // borrow a reference to the signer's NFT collection
        let collectionRef = signer.storage.borrow<auth(NonFungibleToken.Provider, NonFungibleToken.CollectionPublic) &{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>(from: {{.TokenStoragePath}})
            ?? panic("Could not borrow a reference to the owner's collection")
        
        // borrow a public reference to the receivers collection
        let depositRef = recipient.capabilities.get<auth(NonFungibleToken.CollectionPublic) &{NonFungibleToken.CollectionPublic}>({{.TokenPublicCollectionPath}})!.borrow()
            ?? panic("Could not borrow a reference to the receiver's collection")
        
        // withdraw the NFT from the owner's collection
        let nft <- collectionRef.withdraw(withdrawID: id)
        
        // deposit the NFT in the recipient's collection
        depositRef.deposit(token: <-nft)
    }
}
`

const GenericNFTSetup = `
import NonFungibleToken from 0x631e88ae7f1d7c20
import {{.TokenContractAddress}} from {{.TokenContractAddress}}

transaction {
    prepare(signer: auth(Storage, Keys) &Account) {
        // If the account doesn't already have a collection
        if signer.storage.borrow<&{{.TokenContractAddress}}.Collection>(from: {{.TokenStoragePath}}) == nil {

            // Create a new empty collection
            let collection <- {{.TokenContractAddress}}.createEmptyCollection()
            
            // Save it to the account
            signer.storage.save(<-collection, to: {{.TokenStoragePath}})

            // Create a public capability for the collection
            signer.capabilities.publish(
                signer.capabilities.storage.issue<auth(NonFungibleToken.CollectionPublic) &{NonFungibleToken.CollectionPublic, {{.TokenContractAddress}}.{{.TokenPublicCollectionName}}}>({{.TokenStoragePath}}),
                at: {{.TokenPublicCollectionPath}}
            )
        }
    }
}
`

const DeployContract = `
transaction(name: String, code: String) {
    prepare(signer: auth(Storage, Contracts) &Account) {
        signer.contracts.add(name: name, code: code.decodeHex())
    }
}
`
