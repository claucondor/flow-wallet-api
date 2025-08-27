package tokens

// GenericFungibleTransfer transfers fungible tokens between accounts
const GenericFungibleTransfer = `
import FungibleToken from "./FungibleToken.cdc"
import TOKEN_DECLARATION_NAME from TOKEN_ADDRESS

transaction(amount: UFix64, recipient: Address) {
  let sentVault: @FungibleToken.Vault

  prepare(signer: auth(Storage, FungibleToken.Withdraw) &Account) {
    let vaultRef = signer.storage
      .borrow<&TOKEN_DECLARATION_NAME.Vault>(from: TOKEN_VAULT)
      ?? panic("failed to borrow reference to sender vault")

    self.sentVault <- vaultRef.withdraw(amount: amount)
  }

  execute {
    let receiverRef = getAccount(recipient)
      .capabilities
      .borrow<&TOKEN_DECLARATION_NAME.Vault & FungibleToken.Receiver>(TOKEN_RECEIVER)
      ?? panic("failed to borrow reference to recipient vault")

    receiverRef.deposit(from: <-self.sentVault)
  }
}
`

// GenericFungibleSetup sets up a fungible token vault for an account
const GenericFungibleSetup = `
import FungibleToken from "./FungibleToken.cdc"
import TOKEN_DECLARATION_NAME from TOKEN_ADDRESS

transaction {
  prepare(signer: auth(Storage, Capabilities) &Account) {

    let existingVault = signer.storage.borrow<&TOKEN_DECLARATION_NAME.Vault>(from: TOKEN_VAULT)

    if (existingVault != nil) {
        panic("vault exists")
    }

    var vault: @TOKEN_DECLARATION_NAME.Vault? = nil
    if let f = TOKEN_DECLARATION_NAME.createEmptyVault as? fun(): @TOKEN_DECLARATION_NAME.Vault {
        vault <- f()
    } else if let f = TOKEN_DECLARATION_NAME.createEmptyVault as? fun(allowUnrestrictedFlow: Bool): @TOKEN_DECLARATION_NAME.Vault {
        vault <- f(allowUnrestrictedFlow: false)
    } else {
        panic("Could not determine the correct function signature for createEmptyVault")
    }

    signer.storage.save(<-vault!, to: TOKEN_VAULT)

    let cap = signer.capabilities.storage.issue<&TOKEN_DECLARATION_NAME.Vault & FungibleToken.Receiver>(
      TOKEN_VAULT
    )
    signer.capabilities.publish(cap, at: TOKEN_RECEIVER)

    let balanceCap = signer.capabilities.storage.issue<&TOKEN_DECLARATION_NAME.Vault & FungibleToken.Balance>(
      TOKEN_VAULT
    )
    signer.capabilities.publish(balanceCap, at: TOKEN_BALANCE)
  }
}
`