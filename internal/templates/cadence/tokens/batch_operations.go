package tokens

import (
	"bytes"
	"fmt"
	"text/template"
)

// FungibleTokenInfo holds information about a fungible token for batch operations
type FungibleTokenInfo struct {
	ContractName       string
	Address            string
	VaultStoragePath   string
	ReceiverPublicPath string
	BalancePublicPath  string
}

// BatchedFungibleOpsInfo holds information for batch fungible token operations
type BatchedFungibleOpsInfo struct {
	FungibleTokenContractAddress string
	Tokens                       []FungibleTokenInfo
}

// ExecuteTemplate processes a template with the given data
func ExecuteTemplate(name string, templ string, info BatchedFungibleOpsInfo) (string, error) {
	t, err := template.New(name).Parse(templ)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, info)
	if err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.String(), nil
}

// AddFungibleTokenVaultBatchTransaction generates a batch transaction to add fungible token vaults
func AddFungibleTokenVaultBatchTransaction(info BatchedFungibleOpsInfo) (string, error) {
	return ExecuteTemplate("AddFungibleTokenVaultBatchTransaction", AddFungibleTokenVaultBatchTemplate, info)
}

// CreateAccountAndSetupTransaction generates a transaction to create account and setup token vaults
func CreateAccountAndSetupTransaction(info BatchedFungibleOpsInfo) (string, error) {
	return ExecuteTemplate("CreateAccountAndSetupTransaction", CreateAccountAndSetupTemplate, info)
}