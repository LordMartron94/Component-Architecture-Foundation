package coding

import "github.com/component-architecture-foundation/networking/transport"

type MessageCoderInterface interface {
	Decode(message []byte) (transport.Message, error)
	Encode(message transport.Message) ([]byte, error)
}
