import Crypto

transaction(publicKeys: [Crypto.KeyListEntry], contracts: {String: String}) {
	prepare(signer: auth(Keys) &Account) {
		panic("Account initialized with custom script")
	}
}
