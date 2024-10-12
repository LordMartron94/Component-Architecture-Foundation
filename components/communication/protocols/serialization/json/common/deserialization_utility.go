package common

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/protocols/payloads"
	"github.com/component-architecture-foundation/shared"
)

type DeserializationUtility struct {
	Logger logging.HoornLogger
}

func (d *DeserializationUtility) Deserialize(data []byte, payloadType reflect.Type) (payloads.PayloadInterface, error) {
	payload := reflect.New(payloadType).Interface()

	err := json.Unmarshal(data, &payload)
	if err != nil {
		d.Logger.Error("Failed to deserialize JSON data: "+err.Error(), false, shared.MainComponentName)
		return nil, err
	}

	// Type assertion to PayloadInterface
	payloadInterface, ok := payload.(payloads.PayloadInterface)
	if !ok {
		d.Logger.Error("Deserialized payload does not implement PayloadInterface", false, shared.MainComponentName)
		return nil, fmt.Errorf("deserialized payload does not implement PayloadInterface")
	}

	return payloadInterface, nil
}
