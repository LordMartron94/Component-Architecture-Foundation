package strategy

import (
	"github.com/component-architecture-foundation/protocols/payloads"
)

type DeserializationStrategyInterface interface {
	Deserialize(data []byte) (payloads.PayloadInterface, error)
}
