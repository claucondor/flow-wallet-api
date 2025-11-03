package tokens

import (
	"context"
	"strings"

	"github.com/flow-hydraulics/flow-wallet-api/internal/accounts"
	"github.com/flow-hydraulics/flow-wallet-api/internal/chain_events"
	"github.com/flow-hydraulics/flow-wallet-api/internal/flow_helpers"
	"github.com/flow-hydraulics/flow-wallet-api/internal/templates"
	"github.com/onflow/flow-go-sdk"
	log "github.com/sirupsen/logrus"
)

type ChainEventHandler struct {
	AccountService  accounts.Service
	ChainListener   chain_events.Listener
	TemplateService templates.Service
	TokenService    Service
}

func (h *ChainEventHandler) Handle(ctx context.Context, event flow.Event) {
	isDeposit := strings.Contains(event.Type, "Deposit")
	if isDeposit {
		h.handleDeposit(ctx, event)
	}
}

func (h *ChainEventHandler) handleDeposit(ctx context.Context, event flow.Event) {
	// We don't have to care about tokens that are not in the database
	// as we could not even listen to events for them
	token, err := h.TemplateService.TokenFromEvent(event)
	if err != nil {
		log.
			WithFields(log.Fields{"error": err}).
			Warn("Failed to extract token from event")
		return
	}

	// Acceder a los campos del evento usando SearchFieldByName
	amountOrNftID := event.Value.SearchFieldByName("amount")
	if amountOrNftID == nil {
		// Para NFTs, intentar con id
		amountOrNftID = event.Value.SearchFieldByName("id")
	}
	accountAddress := event.Value.SearchFieldByName("to")

	// DEBUG: Log ALL deposit events
	log.
		WithFields(log.Fields{
			"event":   event.Type,
			"to":      accountAddress,
			"amount":  amountOrNftID,
			"txId":    event.TransactionID,
		}).
		Info("Deposit event detected")

	if amountOrNftID == nil || accountAddress == nil {
		log.WithField("event", event.Type).Warn("Could not find required fields in event")
		return
	}

	// Get the target account from database
	account, err := h.AccountService.Details(flow_helpers.HexString(accountAddress.String()))
	if err != nil {
		log.
			WithFields(log.Fields{
				"error":   err,
				"address": accountAddress,
				"event":   event.Type,
			}).
			Warn("Could not get account for deposit event")
		return
	}

	if err = h.TokenService.RegisterDeposit(ctx, token, event.TransactionID, account, amountOrNftID.String()); err != nil {
		log.
			WithFields(log.Fields{"error": err}).
			Warn("Error while registering a deposit")
		return
	}

	log.
		WithFields(log.Fields{
			"token":         token.Name,
			"account":       accountAddress,
			"amountOrNftID": amountOrNftID,
		}).
		Debug("New deposit")
}
