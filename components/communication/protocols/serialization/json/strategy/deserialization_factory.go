package strategy

import (
	"reflect"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/protocols/payloads"
	"github.com/component-architecture-foundation/protocols/serialization/json/common"
)

func newDeserializer(payloadType reflect.Type, logger logging.HoornLogger) DeserializationStrategyInterface {
	return &genericDeserializer{
		deserializationUtility: common.DeserializationUtility{Logger: logger},
		payloadType:            payloadType,
	}
}

type genericDeserializer struct {
	deserializationUtility common.DeserializationUtility
	payloadType            reflect.Type
}

func (g genericDeserializer) Deserialize(data []byte) (payloads.PayloadInterface, error) {
	return g.deserializationUtility.Deserialize(data, g.payloadType)
}
