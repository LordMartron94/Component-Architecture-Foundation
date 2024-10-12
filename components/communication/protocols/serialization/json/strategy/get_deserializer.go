package strategy

import (
	"fmt"
	"reflect"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/protocols/payloads"
	"github.com/component-architecture-foundation/shared"
)

type GetDeserializer struct {
	Logger logging.HoornLogger
}

func (gD GetDeserializer) Get(serializerType string) (DeserializationStrategyInterface, error) {
	gD.Logger.Debug("Getting deserialization strategy for serializer type: "+serializerType, false, shared.MainComponentName)

	switch serializerType {
	case "event":
		return newDeserializer(reflect.TypeOf(payloads.EventPayload{}), gD.Logger), nil
	case "request":
		return newDeserializer(reflect.TypeOf(payloads.RequestPayload{}), gD.Logger), nil
	case "response":
		return newDeserializer(reflect.TypeOf(payloads.ResponsePayload{}), gD.Logger), nil
	default:
		gD.Logger.Error("Unsupported serializer type: "+serializerType, false, shared.MainComponentName)
		return nil, fmt.Errorf("unsupported serializer type: %s", serializerType)
	}
}
