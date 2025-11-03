package scripts

// GenericFungibleBalance checks the balance of a fungible token for an account
const GenericFungibleBalance = `
import FungibleToken from 0x9a0766d93b6608b7
import TOKEN_DECLARATION_NAME from TOKEN_ADDRESS

access(all)
view fun main(account: Address): UFix64 {

    let vaultRef = getAccount(account)
        .capabilities
        .get<&{FungibleToken.Balance}>(TOKEN_BALANCE)
        .borrow()
        ?? panic("failed to borrow reference to vault")

    return vaultRef.balance
}
`