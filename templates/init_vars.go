package templates

import "github.com/onflow/flow-go-sdk"

var KnownAddresses = templateVariables{
	"FungibleToken.cdc": knownAddresses{
		flow.Emulator: "0xee82856bf20e2aa6",
		flow.Testnet:  "0x9a0766d93b6608b7",
		flow.Mainnet:  "0xf233dcee88fe0abe",
	},
	"NonFungibleToken.cdc": knownAddresses{
		flow.Emulator: "0xf8d6e0586b0a20c7",
		flow.Testnet:  "0x631e88ae7f1d7c20",
		flow.Mainnet:  "0x1d7e57aa55817448",
	},
	"FlowToken.cdc": knownAddresses{
		flow.Emulator: "0x0ae53cb6e3f42a79",
		flow.Testnet:  "0x7e60df042a9c0868",
		flow.Mainnet:  "0x1654653399040a61",
	},
	"FUSD.cdc": knownAddresses{
		flow.Emulator: "0xf8d6e0586b0a20c7",
		flow.Testnet:  "0xe223d8a629e49c68",
		flow.Mainnet:  "0x3c5959b568896393",
	},
}

func init() {
	knownAddressesReplacers = makeReplacers(KnownAddresses)
}
