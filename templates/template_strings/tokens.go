package template_strings

// MyToken es un template de token fungible personalizado con características básicas
const MyToken = `
import FungibleToken from 0x9a0766d93b6608b7

access(all) contract MyToken {
    
    // Eventos
    access(all) event TokensInitialized(initialSupply: UFix64)
    access(all) event TokensMinted(amount: UFix64)
    access(all) event TokensBurned(amount: UFix64)
    access(all) event TokensWithdrawn(amount: UFix64, from: Address?)
    access(all) event TokensDeposited(amount: UFix64, to: Address?)
    access(all) event MinterCreated(allowedAmount: UFix64)
    access(all) event BurnerCreated()
    
    // Propiedades públicas del token
    access(all) let symbol: String
    access(all) let name: String
    access(all) var totalSupply: UFix64

    // Rutas para almacenar recursos
    access(all) let VaultStoragePath: StoragePath
    access(all) let ReceiverPublicPath: PublicPath
    access(all) let BalancePublicPath: PublicPath
    access(all) let AdminStoragePath: StoragePath
    access(all) let MinterStoragePath: StoragePath

    // El recurso Vault implementa la interfaz Vault de FungibleToken
    access(all) resource Vault: FungibleToken.Vault {
        // El saldo del Vault
        access(all) var balance: UFix64

        // Inicializa un nuevo Vault con saldo inicial
        init(balance: UFix64) {
            self.balance = balance
        }

        // Retira tokens del Vault
        access(all) fun withdraw(amount: UFix64): @FungibleToken.Vault {
            self.balance = self.balance - amount
            emit TokensWithdrawn(amount: amount, from: self.owner?.address)
            return <-create Vault(balance: amount)
        }

        // Deposita tokens en el Vault
        access(all) fun deposit(from: @FungibleToken.Vault) {
            let vault <- from as! @MyToken.Vault
            let amount = vault.balance
            self.balance = self.balance + amount
            emit TokensDeposited(amount: amount, to: self.owner?.address)
            vault.balance = 0.0
            destroy vault
        }
    }

    // Crea un Vault vacío
    access(all) fun createEmptyVault(): @FungibleToken.Vault {
        return <-create Vault(balance: 0.0)
    }

    // Recurso que permite al admin gestionar el token
    access(all) resource Administrator {
        // Crea un nuevo Minter
        access(all) fun createNewMinter(allowedAmount: UFix64): @Minter {
            emit MinterCreated(allowedAmount: allowedAmount)
            return <-create Minter(allowedAmount: allowedAmount)
        }

        // Crea un nuevo Burner
        access(all) fun createNewBurner(): @Burner {
            emit BurnerCreated()
            return <-create Burner()
        }
    }

    // Recurso para acuñar nuevos tokens
    access(all) resource Minter {
        // Cantidad de tokens que este minter puede acuñar
        access(all) var allowedAmount: UFix64

        // Función para acuñar nuevos tokens
        access(all) fun mintTokens(amount: UFix64): @MyToken.Vault {
            pre {
                amount > 0.0: "La cantidad a acuñar debe ser mayor que cero"
                amount <= self.allowedAmount: "La cantidad debe ser menor o igual que la cantidad permitida"
            }
            self.allowedAmount = self.allowedAmount - amount
            MyToken.totalSupply = MyToken.totalSupply + amount
            emit TokensMinted(amount: amount)
            return <-create Vault(balance: amount)
        }

        init(allowedAmount: UFix64) {
            self.allowedAmount = allowedAmount
        }
    }

    // Recurso para quemar tokens
    access(all) resource Burner {
        // Función para quemar tokens
        access(all) fun burnTokens(from: @FungibleToken.Vault) {
            let vault <- from as! @MyToken.Vault
            let amount = vault.balance
            MyToken.totalSupply = MyToken.totalSupply - amount
            emit TokensBurned(amount: amount)
            destroy vault
        }
    }

    init() {
        self.symbol = "MYT"
        self.name = "My Token"
        self.totalSupply = 0.0

        self.VaultStoragePath = /storage/myTokenVault
        self.ReceiverPublicPath = /public/myTokenReceiver
        self.BalancePublicPath = /public/myTokenBalance
        self.AdminStoragePath = /storage/myTokenAdmin
        self.MinterStoragePath = /storage/myTokenMinter

        // Crear el admin y guardarlo
        let admin <- create Administrator()
        
        // Crear un minter y acuñar el suministro inicial
        let minter <- admin.createNewMinter(allowedAmount: 1000000.0)
        let initialTokens <- minter.mintTokens(amount: 1000000.0)
        
        // Guardar el Vault con los tokens iniciales
        self.account.storage.save(<-initialTokens, to: self.VaultStoragePath)
        
        // Exponer las capacidades públicas
        self.account.capabilities.publish(
            self.account.capabilities.storage.issue<&{FungibleToken.Receiver}>(self.VaultStoragePath),
            at: self.ReceiverPublicPath
        )
        
        self.account.capabilities.publish(
            self.account.capabilities.storage.issue<&{FungibleToken.Balance}>(self.VaultStoragePath),
            at: self.BalancePublicPath
        )
        
        // Guardar el admin y el minter
        self.account.storage.save(<-admin, to: self.AdminStoragePath)
        self.account.storage.save(<-minter, to: self.MinterStoragePath)

        emit TokensInitialized(initialSupply: self.totalSupply)
    }
}
`
