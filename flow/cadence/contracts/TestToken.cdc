import FungibleToken from "./FungibleToken.cdc"
import "ViewResolver"

access(all) contract TestToken: FungibleToken {

    // Total supply of Test tokens in existence
    access(all) var totalSupply: UFix64

    // Event that is emitted when the contract is created
    access(all) event TokensInitialized(initialSupply: UFix64)

    // Event that is emitted when tokens are withdrawn from a Vault
    access(all) event TokensWithdrawn(amount: UFix64, from: Address?, uuid: UInt64, providerUUID: UInt64)

    // Event that is emitted when tokens are deposited to a Vault
    access(all) event TokensDeposited(amount: UFix64, to: Address?, uuid: UInt64, receiverUUID: UInt64)

    // Event that is emitted when new tokens are minted
    access(all) event TokensMinted(amount: UFix64)

    // Event that is emitted when tokens are destroyed
    access(all) event TokensBurned(amount: UFix64)

    // Event that is emitted when a new minter resource is created
    access(all) event MinterCreated(allowedAmount: UFix64)

    // Event that is emitted when a new burner resource is created
    access(all) event BurnerCreated()

    // Vault
    //
    // Each user stores an instance of only the Vault in their storage
    // The functions in the Vault and governed by the pre and post conditions
    // in FungibleToken when they are called.
    // The checks happen at runtime whenever a function is called.
    //
    // Resources can only be created in the context of the contract that they
    // are defined in, so there is no way for a malicious user to create Vaults
    // out of thin air. A special Minter resource needs to be defined to mint
    // new tokens.
    //
    access(all) resource Vault: FungibleToken.Provider, FungibleToken.Receiver, FungibleToken.Balance, ViewResolver.Resolver {

        // holds the balance of a users tokens
        access(all) var balance: UFix64

        // Event emitted when tokens are withdrawn from the vault
        access(all) event TokenResourceDestroyed(uuid: UInt64, balance: UFix64)

        // initialize the balance at resource creation time
        init(balance: UFix64) {
            self.balance = balance
        }

        // withdraw
        //
        // Function that takes an integer amount as an argument
        // and withdraws that amount from the Vault.
        // It creates a new temporary Vault that is used to hold
        // the money that is being transferred. It returns the newly
        // created Vault to the context that called so it can be deposited
        // elsewhere.
        //
        access(FungibleToken.Withdraw) fun withdraw(amount: UFix64): @{FungibleToken.Vault} {
            pre {
                self.balance >= amount:
                    "Amount withdrawn must be less than or equal than the balance of the Vault"
            }
            self.balance = self.balance - amount
            let vault <- create Vault(balance: amount)
            emit TokensWithdrawn(amount: amount, from: self.owner?.address, uuid: vault.uuid, providerUUID: self.uuid)
            return <-vault
        }

        // deposit
        //
        // Function that takes a Vault object as an argument and adds
        // its balance to the balance of the owners Vault.
        // It is allowed to destroy the sent Vault because the Vault
        // was a temporary holder of the tokens. The Vault's balance has
        // been consumed and therefore can be destroyed.
        access(all) fun deposit(from: @{FungibleToken.Vault}) {
            pre {
                from.balance > 0.0:
                    "FungibleToken.Receiver.deposit: Cannot deposit a Vault with zero balance."
                from.getType() != self.getType():
                    "FungibleToken.Receiver.deposit: Cannot deposit a Vault of the same type as the receiver."
            }
            let vault <- from as! @TestToken.Vault
            self.balance = self.balance + vault.balance
            emit TokensDeposited(amount: vault.balance, to: self.owner?.address, uuid: vault.uuid, receiverUUID: self.uuid)
            vault.balance = 0.0
            destroy vault
        }

        access(all) fun createEmptyVault(): @{FungibleToken.Vault} {
            return <-create Vault(balance: 0.0)
        }

        access(all) view fun getSupportedVaultTypes(): [Type] {
            return [Type<@TestToken.Vault>()]
        }

        access(all) view fun isAvailableToWithdraw(amount: UFix64): Bool {
            return self.balance >= amount
        }

        access(all) view fun isAvailableToDeposit(from: &{FungibleToken.Vault}): Bool {
            return from.isInstance(Type<@TestToken.Vault>())
        }

        access(all) view fun getViews(): [Type] {
            return [Type<FungibleToken.Balance>()]
        }

        access(all) view fun resolveView(_ view: Type): AnyStruct? {
            if view == Type<FungibleToken.Balance>() {
                return FungibleToken.Balance(balance: self.balance)
            }
            return nil
        }

        destroy() {
            TestToken.totalSupply = TestToken.totalSupply - self.balance
            emit TokenResourceDestroyed(uuid: self.uuid, balance: self.balance)
        }
    }

    // createEmptyVault
    //
    // Function that creates a new Vault with a balance of zero
    // and returns it to the calling context. A user must call this function
    // and store the returned Vault in their storage in order to allow their
    // account to be able to receive deposits of this token type.
    //
    access(all) fun createEmptyVault(): @{FungibleToken.Vault} {
        return <-create Vault(balance: 0.0)
    }

    access(all) resource Administrator {
        // createNewMinter
        //
        // Function that creates and returns a new minter resource
        //
        access(all) fun createNewMinter(allowedAmount: UFix64): @Minter {
            emit MinterCreated(allowedAmount: allowedAmount)
            return <-create Minter(allowedAmount: allowedAmount)
        }

        // createNewBurner
        //
        // Function that creates and returns a new burner resource
        //
        access(all) fun createNewBurner(): @Burner {
            emit BurnerCreated()
            return <-create Burner()
        }
    }

    // Minter
    //
    // Resource object that token admin accounts can hold to mint new tokens.
    //
    access(all) resource Minter {

        // the amount of tokens that the minter is allowed to mint
        access(all) var allowedAmount: UFix64

        // mintTokens
        //
        // Function that mints new tokens, adds them to the total supply,
        // and returns them to the calling context.
        //
        access(all) fun mintTokens(amount: UFix64): @TestToken.Vault {
            pre {
                amount > UFix64(0): "Amount minted must be greater than zero"
                amount <= self.allowedAmount: "Amount minted must be less than the allowed amount"
            }
            TestToken.totalSupply = TestToken.totalSupply + amount
            self.allowedAmount = self.allowedAmount - amount
            emit TokensMinted(amount: amount)
            return <-create Vault(balance: amount)
        }

        init(allowedAmount: UFix64) {
            self.allowedAmount = allowedAmount
        }
    }

    // Burner
    //
    // Resource object that token admin accounts can hold to burn tokens.
    //
    access(all) resource Burner {

        // burnTokens
        //
        // Function that destroys a Vault instance, effectively burning the tokens.
        //
        // Note: the burned tokens are automatically subtracted from the
        // total supply in the Vault destructor.
        //
        access(all) fun burnTokens(from: @{FungibleToken.Vault}) {
            let vault <- from as! @TestToken.Vault
            let amount = vault.balance
            destroy vault
            emit TokensBurned(amount: amount)
        }
    }

    init(adminAccount: auth(Storage, Keys, Capabilities) &Account) {
        self.totalSupply = 0.0

        // Create the Vault with the total supply of tokens and save it in storage
        //
        let vault <- create Vault(balance: self.totalSupply)
        adminAccount.storage.save(<-vault, to: /storage/testTokenVault)

        // Create a public capability to the stored Vault that only exposes
        // the `deposit` method through the `Receiver` interface
        //
        let receiverCap = adminAccount.capabilities.storage.issue<&TestToken.Vault{FungibleToken.Receiver}>(
            /storage/testTokenVault
        )
        adminAccount.capabilities.publish(receiverCap, at: /public/testTokenReceiver)

        // Create a public capability to the stored Vault that only exposes
        // the `balance` field through the `Balance` interface
        //
        let balanceCap = adminAccount.capabilities.storage.issue<&TestToken.Vault{FungibleToken.Balance}>(
            /storage/testTokenVault
        )
        adminAccount.capabilities.publish(balanceCap, at: /public/testTokenBalance)

        let admin <- create Administrator()
        adminAccount.storage.save(<-admin, to: /storage/testTokenAdmin)

        // Emit an event that shows that the contract was initialized
        emit TokensInitialized(initialSupply: self.totalSupply)
    }
}
