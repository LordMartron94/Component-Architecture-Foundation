package serialization

import "github.com/component-architecture-foundation/protocols/payloads"

type SerializerInterface interface {
	Serialize(data payloads.PayloadInterface) ([]byte, error)
	Deserialize(data []byte) (payloads.PayloadInterface, error)
}
