package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/flow-hydraulics/flow-wallet-api/errors"
	"github.com/flow-hydraulics/flow-wallet-api/templates"
	"github.com/flow-hydraulics/flow-wallet-api/tokens"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (s *Tokens) SetupFunc(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	tokenName := vars["tokenName"]

	// Decide whether to serve sync or async, default async
	sync := r.FormValue(SyncQueryParameter) != ""
	job, transaction, err := s.service.Setup(r.Context(), sync, tokenName, address)

	if err != nil {
		handleError(rw, r, err)
		return
	}

	var res interface{}
	if sync {
		res = transaction.ToJSONResponse()
	} else {
		res = job.ToJSONResponse()
	}

	handleJsonResponse(rw, http.StatusCreated, res)
}

func (s *Tokens) MakeAccountTokensFunc(tType templates.TokenType) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		a := vars["address"]

		res, err := s.service.AccountTokens(a, tType)

		if err != nil {
			handleError(rw, r, err)
			return
		}

		handleJsonResponse(rw, http.StatusOK, res)
	}
}

func (s *Tokens) DetailsFunc(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	tokenName := vars["tokenName"]

	res, err := s.service.Details(r.Context(), tokenName, address)

	if err != nil {
		handleError(rw, r, err)
		return
	}

	handleJsonResponse(rw, http.StatusOK, res)
}

func (s *Tokens) CreateWithdrawalFunc(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	tokenName := vars["tokenName"]

	var withdrawal tokens.WithdrawalRequest

	if r.Body == nil || r.Body == http.NoBody {
		err := &errors.RequestError{StatusCode: http.StatusBadRequest, Err: fmt.Errorf("empty body")}
		handleError(rw, r, err)
		return
	}

	// Try to decode the request body.
	if err := json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
		err = &errors.RequestError{StatusCode: http.StatusBadRequest, Err: fmt.Errorf("invalid body")}
		handleError(rw, r, err)
		return
	}

	withdrawal.TokenName = tokenName

	// Decide whether to serve sync or async, default async
	sync := r.FormValue(SyncQueryParameter) != ""
	job, transaction, err := s.service.CreateWithdrawal(r.Context(), sync, address, withdrawal)

	if err != nil {
		handleError(rw, r, err)
		return
	}

	var res interface{}
	if sync {
		res = transaction.ToJSONResponse()
	} else {
		res = job.ToJSONResponse()
	}

	handleJsonResponse(rw, http.StatusCreated, res)
}

func (s *Tokens) ListWithdrawalsFunc(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	tokenName := vars["tokenName"]

	res, err := s.service.ListWithdrawals(address, tokenName)

	if err != nil {
		handleError(rw, r, err)
		return
	}

	handleJsonResponse(rw, http.StatusOK, res)
}

func (s *Tokens) GetWithdrawalFunc(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	tokenName := vars["tokenName"]
	txId := vars["transactionId"]

	res, err := s.service.GetWithdrawal(address, tokenName, txId)

	if err != nil {
		handleError(rw, r, err)
		return
	}

	handleJsonResponse(rw, http.StatusOK, res)
}

func (s *Tokens) ListDepositsFunc(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	tokenName := vars["tokenName"]

	res, err := s.service.ListDeposits(address, tokenName)

	if err != nil {
		handleError(rw, r, err)
		return
	}

	handleJsonResponse(rw, http.StatusOK, res)
}

func (s *Tokens) GetDepositFunc(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	tokenName := vars["tokenName"]
	transactionId := vars["transactionId"]

	res, err := s.service.GetDeposit(address, tokenName, transactionId)

	if err != nil {
		handleError(rw, r, err)
		return
	}

	handleJsonResponse(rw, http.StatusOK, res)
}

// TokenBalanceFunc maneja la solicitud de balance de un token
func (s *Tokens) TokenBalanceFunc(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	tokenName := vars["token"]

	res, err := s.service.Details(r.Context(), tokenName, address)

	if err != nil {
		handleError(rw, r, err)
		return
	}

	handleJsonResponse(rw, http.StatusOK, res)
}

// TokenTransferFunc maneja la solicitud de transferencia de un token
func (s *Tokens) TokenTransferFunc(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	tokenName := vars["token"]

	var transferRequest tokens.WithdrawalRequest
	err := json.NewDecoder(r.Body).Decode(&transferRequest)
	if err != nil {
		handleError(rw, r, &errors.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("invalid body: %w", err),
		})
		return
	}

	// Añadir el nombre del token a la solicitud
	transferRequest.TokenName = tokenName

	// Decidir si la operación es sincrónica o asincrónica
	sync := r.FormValue(SyncQueryParameter) != ""
	job, transaction, err := s.service.CreateWithdrawal(r.Context(), sync, address, transferRequest)

	if err != nil {
		handleError(rw, r, err)
		return
	}

	var res interface{}
	if sync {
		res = transaction.ToJSONResponse()
	} else {
		res = job.ToJSONResponse()
	}

	handleJsonResponse(rw, http.StatusCreated, res)
}

// DeployTokenFunc maneja la solicitud para desplegar un contrato de token
func (s *Tokens) DeployTokenFunc(rw http.ResponseWriter, r *http.Request) {
	var req struct {
		TokenName         string            `json:"tokenName"`         // Nombre del token a desplegar
		Address           string            `json:"address"`           // Dirección de la cuenta donde se desplegará
		DeployBasicDeps   bool              `json:"deployBasicDeps"`   // Si es true, despliega FungibleToken antes
		DeployCustomDeps  []string          `json:"deployCustomDeps"`  // Lista de tokens adicionales a desplegar antes
		ContractAddresses map[string]string `json:"contractAddresses"` // Mapa de direcciones de contratos base
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(rw, r, &errors.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("invalid body: %w", err),
		})
		return
	}

	if req.TokenName == "" {
		handleError(rw, r, &errors.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("tokenName is required"),
		})
		return
	}

	if req.Address == "" {
		handleError(rw, r, &errors.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("address is required"),
		})
		return
	}

	// Decidir si la operación es sincrónica o asincrónica
	sync := r.FormValue(SyncQueryParameter) != ""

	// Configurar variables de entorno temporales con las direcciones de los contratos
	if req.ContractAddresses != nil {
		for k, v := range req.ContractAddresses {
			envVarName := fmt.Sprintf("FLOW_WALLET_%s_ADDRESS", strings.ToUpper(k))
			oldValue := os.Getenv(envVarName)
			os.Setenv(envVarName, v)
			defer os.Setenv(envVarName, oldValue) // Restaurar valor original al finalizar

			logrus.WithFields(logrus.Fields{
				"contract": k,
				"address":  v,
				"envVar":   envVarName,
			}).Info("Configurando dirección de contrato")
		}
	}

	// Primero desplegar las dependencias básicas si se solicita
	if req.DeployBasicDeps {
		// Intenta desplegar FungibleToken primero
		err = s.service.DeployTokenContractForAccount(r.Context(), sync, "FungibleToken", req.Address)
		if err != nil {
			// Solo reportar, no fallar
			logrus.WithFields(logrus.Fields{
				"error": err,
				"token": "FungibleToken",
			}).Warn("Error al desplegar dependencia básica")
		}

		// Intenta desplegar FlowToken después
		err = s.service.DeployTokenContractForAccount(r.Context(), sync, "FlowToken", req.Address)
		if err != nil {
			// Solo reportar, no fallar
			logrus.WithFields(logrus.Fields{
				"error": err,
				"token": "FlowToken",
			}).Warn("Error al desplegar dependencia básica")
		}
	}

	// Desplegar dependencias personalizadas si se especifican
	for _, depToken := range req.DeployCustomDeps {
		err = s.service.DeployTokenContractForAccount(r.Context(), sync, depToken, req.Address)
		if err != nil {
			// Solo reportar, no fallar
			logrus.WithFields(logrus.Fields{
				"error": err,
				"token": depToken,
			}).Warn("Error al desplegar dependencia personalizada")
		}
	}

	// Finalmente, desplegar el token principal
	err = s.service.DeployTokenContractForAccount(r.Context(), sync, req.TokenName, req.Address)
	if err != nil {
		handleError(rw, r, err)
		return
	}

	result := map[string]interface{}{
		"status":    "success",
		"message":   fmt.Sprintf("Token %s deployed to address %s", req.TokenName, req.Address),
		"tokenName": req.TokenName,
		"address":   req.Address,
	}

	handleJsonResponse(rw, http.StatusOK, result)
}
