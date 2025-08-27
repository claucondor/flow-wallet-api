package admin

// AddProposalKeys adds multiple proposal keys to an admin account for parallel transaction execution
const AddProposalKeys = `
transaction(adminKeyIndex: Int, numProposalKeys: UInt16) {
  prepare(account: auth(Keys) &Account) {
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

// AddAccountContract adds a contract to an account with admin privileges
const AddAccountContract = `
transaction(name: String, code: String) {
	prepare(signer: auth(Contracts) &Account) {
		signer.contracts.add(name: name, code: code.decodeHex(), adminAccount: signer)
	}
}
`