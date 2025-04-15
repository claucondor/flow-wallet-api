import FungibleToken from 0x9a0766d93b6608b7
import FlowToken from 0x0ae53cb6e3f42a79

transaction(amount: UFix64, recipient: Address) {
  let sentVault: @FungibleToken.Vault
  prepare(signer: auth(Keys) &Account) {
    let vaultRef = signer.storage.borrow<auth(FungibleToken.Withdraw) &FungibleToken.Vault>(from: /storage/flowTokenVault)
      ?? panic("Could not borrow reference to the owner's Vault!")

    self.sentVault <- vaultRef.withdraw(amount: amount)
  }

  execute {
    let receiverRef = getAccount(recipient)
      .capabilities.get<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)
      .borrow() ?? panic("Could not borrow receiver reference to the recipient's Vault")

    receiverRef.deposit(from: <-self.sentVault)
  }
}
