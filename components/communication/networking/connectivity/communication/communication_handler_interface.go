package communication

import "github.com/component-architecture-foundation/networking/transport"

type CommunicationHandlerInterface interface {
	SendResponse(id transport.ComponentID, payload []byte) error
	SendRequest(id transport.ComponentID, payload transport.MessagePayload) error
}
