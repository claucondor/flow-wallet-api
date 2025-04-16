package tokens

import (
	"context"
	"fmt"
	"strings"

	"github.com/flow-hydraulics/flow-wallet-api/accounts"
	"github.com/flow-hydraulics/flow-wallet-api/flow_helpers"
	"github.com/flow-hydraulics/flow-wallet-api/templates"
	"github.com/flow-hydraulics/flow-wallet-api/templates/template_strings"
	"github.com/onflow/flow-go-sdk"
	flow_templates "github.com/onflow/flow-go-sdk/templates"
	log "github.com/sirupsen/logrus"
)

// DeployTokenContractForAccount is used for testing purposes.
func (s *ServiceImpl) DeployTokenContractForAccount(ctx context.Context, runSync bool, tokenName, address string) error {
	// Check if the input is a valid address
	address, err := flow_helpers.ValidateAddress(address, s.cfg.ChainID)
	if err != nil {
		return err
	}

	// Verificar si los contratos son contratos estándar que ya están desplegados en la red
	if s.cfg.ChainID == flow.Testnet {
		// En testnet, verificamos si el token es uno de los contratos estándar que ya están desplegados
		if tokenName == "FlowToken" {
			log.Info("FlowToken ya está desplegado en testnet en la dirección 0x7e60df042a9c0868, omitiendo despliegue")
			return nil
		}
		if tokenName == "FungibleToken" {
			log.Info("FungibleToken ya está desplegado en testnet en la dirección 0x9a0766d93b6608b7, omitiendo despliegue")
			return nil
		}
		if tokenName == "NonFungibleToken" {
			log.Info("NonFungibleToken ya está desplegado en testnet en la dirección 0x631e88ae7f1d7c20, omitiendo despliegue")
			return nil
		}
	} else if s.cfg.ChainID == flow.Emulator && tokenName == "FlowToken" {
		// Omitir FlowToken en el emulador ya que ya está desplegado internamente
		log.Info("FlowToken ya está desplegado por defecto en el emulador, omitiendo despliegue")
		return nil
	}

	token, err := s.templates.GetTokenByName(tokenName)
	if err != nil {
		return err
	}

	n := token.Name

	// Verificar si estamos en el emulador y si el token es FUSD
	if s.cfg.ChainID == flow.Emulator && n == "FUSD" {
		log.Info("Desplegando contrato FUSD en el emulador")

		// Si hay algún error, podemos ignorarlo si es porque el contrato ya existe
		tmplStr, err := template_strings.GetByName(n)
		if err != nil {
			return err
		}

		src, err := templates.TokenCode(s.cfg.ChainID, token, tmplStr)
		if err != nil {
			return err
		}

		// Usar la dirección de la cuenta del administrador en el emulador
		c := flow_templates.Contract{Name: n, Source: src}

		err = accounts.AddContract(ctx, s.fc, s.km, address, c, s.cfg.TransactionTimeout)
		if err != nil && !strings.Contains(err.Error(), "cannot overwrite existing contract") {
			return err
		}

		return nil
	}

	// Para cualquier otro token que no sea un contrato estándar
	log.WithFields(log.Fields{
		"token":   tokenName,
		"network": s.cfg.ChainID.String(),
		"address": address,
	}).Info("Desplegando contrato de token personalizado")

	tmplStr, err := template_strings.GetByName(n)
	if err != nil {
		return fmt.Errorf("error obteniendo template para %s: %w", n, err)
	}

	src, err := templates.TokenCode(s.cfg.ChainID, token, tmplStr)
	if err != nil {
		return fmt.Errorf("error generando código para %s: %w", n, err)
	}

	c := flow_templates.Contract{Name: n, Source: src}

	err = accounts.AddContract(ctx, s.fc, s.km, address, c, s.cfg.TransactionTimeout)
	if err != nil && !strings.Contains(err.Error(), "cannot overwrite existing contract") {
		return fmt.Errorf("error desplegando contrato %s: %w", n, err)
	}

	log.WithFields(log.Fields{
		"token":   tokenName,
		"network": s.cfg.ChainID.String(),
		"address": address,
	}).Info("Contrato desplegado exitosamente")

	return nil
}
