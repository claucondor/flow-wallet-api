package tokens

import (
	"encoding/json"

	"github.com/onflow/cadence"
)

type Balance struct {
	CadenceValue cadence.Value
}

func (b *Balance) MarshalJSON() ([]byte, error) {
	if b.CadenceValue == nil {
		// Not able to omit the balance field as it is in a "parent" struct
		// So using JSON null here
		return json.Marshal(nil)
	}

	// Only handle fixed point numbers differently, rest can use the default
	typeID := b.CadenceValue.Type().ID()
	if typeID == "UFix64" || typeID == "Fix64" {
		return json.Marshal(b.CadenceValue.String())
	}

	// Convertir a valor JSON usando String()
	return json.Marshal(b.CadenceValue.String())
}
