package template_strings

import (
	"github.com/flow-hydraulics/flow-wallet-api/internal/templates/cadence/accounts"
	"github.com/flow-hydraulics/flow-wallet-api/internal/templates/cadence/admin"
	"github.com/flow-hydraulics/flow-wallet-api/internal/templates/cadence/scripts"
	"github.com/flow-hydraulics/flow-wallet-api/internal/templates/cadence/tokens"
)

// All transactions now organized and imported from their respective modules
// This maintains backward compatibility while providing clean organization

// Account Management Transactions
const CreateAccount = accounts.CreateAccount
const AddAccountKeysTransaction = accounts.AddAccountKeys

// Token Management Transactions  
const GenericFungibleTransfer = tokens.GenericFungibleTransfer
const GenericFungibleSetup = tokens.GenericFungibleSetup

// Scripts for Reading Data
const GenericFungibleBalance = scripts.GenericFungibleBalance

// Admin Operations
const AddProposalKeyTransaction = admin.AddProposalKeys
const AddAccountContractWithAdmin = admin.AddAccountContract

// Batch Operations Templates
const CreateAccountAndSetupTransactionTemplate = tokens.CreateAccountAndSetupTemplate
const AddFungibleTokenVaultBatchTransactionTemplate = tokens.AddFungibleTokenVaultBatchTemplate

// Re-export batch operation types and functions
type FungibleTokenInfo = tokens.FungibleTokenInfo
type BatchedFungibleOpsInfo = tokens.BatchedFungibleOpsInfo
var AddFungibleTokenVaultBatchTransaction = tokens.AddFungibleTokenVaultBatchTransaction
var CreateAccountAndSetupTransaction = tokens.CreateAccountAndSetupTransaction