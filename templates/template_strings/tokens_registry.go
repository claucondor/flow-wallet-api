package template_strings

// Registros de tokens fungibles comunes
var (
	FungibleTokenInfoRegistry = map[string]FungibleTokenInfo{
		"FlowToken": {
			ContractName:       "FlowToken",
			Address:            "0x7e60df042a9c0868",
			VaultStoragePath:   "/storage/flowTokenVault",
			ReceiverPublicPath: "/public/flowTokenReceiver",
			BalancePublicPath:  "/public/flowTokenBalance",
		},
		"FungibleToken": {
			ContractName:       "FungibleToken",
			Address:            "0x9a0766d93b6608b7",
			VaultStoragePath:   "",
			ReceiverPublicPath: "",
			BalancePublicPath:  "",
		},
		"FUSD": {
			ContractName:       "FUSD",
			Address:            "0xe223d8a629e49c68",
			VaultStoragePath:   "/storage/fusdVault",
			ReceiverPublicPath: "/public/fusdReceiver",
			BalancePublicPath:  "/public/fusdBalance",
		},
		"MyToken": {
			ContractName:       "MyToken",
			Address:            "",
			VaultStoragePath:   "/storage/myTokenVault",
			ReceiverPublicPath: "/public/myTokenReceiver",
			BalancePublicPath:  "/public/myTokenBalance",
		},
	}

	// Registros de tokens no fungibles comunes
	NonFungibleTokenInfoRegistry = map[string]string{
		"NonFungibleToken": "0x631e88ae7f1d7c20",
		"ExampleNFT":       "",
	}
)
