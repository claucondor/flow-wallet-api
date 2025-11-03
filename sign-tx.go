package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run sign-tx.go <encodedTx> <privateKey>")
		os.Exit(1)
	}

	encodedTxHex := os.Args[1]
	privateKeyHex := os.Args[2]

	// Decode the transaction
	txBytes, err := hex.DecodeString(encodedTxHex)
	if err != nil {
		fmt.Printf("Error decoding transaction hex: %v\n", err)
		os.Exit(1)
	}

	tx, err := flow.DecodeTransaction(txBytes)
	if err != nil {
		fmt.Printf("Error decoding transaction: %v\n", err)
		os.Exit(1)
	}

	// Decode the private key
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		fmt.Printf("Error decoding private key: %v\n", err)
		os.Exit(1)
	}

	privateKey, err := crypto.DecodePrivateKey(crypto.ECDSA_P256, privateKeyBytes)
	if err != nil {
		fmt.Printf("Error loading private key: %v\n", err)
		os.Exit(1)
	}

	// Create signer
	signer, err := crypto.NewInMemorySigner(privateKey, crypto.SHA3_256)
	if err != nil {
		fmt.Printf("Error creating signer: %v\n", err)
		os.Exit(1)
	}

	// Sign the payload (user signs as authorizer)
	// The proposer address and key index should be the first authorizer
	if len(tx.Authorizers) == 0 {
		fmt.Println("No authorizers in transaction")
		os.Exit(1)
	}

	authorizerAddress := tx.Authorizers[0]
	keyIndex := tx.ProposalKey.KeyIndex

	err = tx.SignPayload(authorizerAddress, keyIndex, signer)
	if err != nil {
		fmt.Printf("Error signing transaction: %v\n", err)
		os.Exit(1)
	}

	// Encode the signed transaction
	signedTxBytes := tx.Encode()
	signedTxHex := hex.EncodeToString(signedTxBytes)

	fmt.Println(signedTxHex)
}
