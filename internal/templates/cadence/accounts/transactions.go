package accounts

// CreateAccount creates a new Flow account with the provided public keys
const CreateAccount = `
transaction(publicKeys: [String]) {
	prepare(signer: auth(CreateAccount) &Account) {
		let acct = Account(payer: signer)

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

// AddAccountKeys adds multiple keys to an existing account
const AddAccountKeys = `
transaction(publicKeys: [String]) {
  prepare(signer: auth(Keys) &Account) {
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