package communication

import "github.com/component-architecture-foundation/networking/transport"

type CommunicationHandlerInterface interface {
	SendResponse(id string, payload []byte, targetUUID string) (string, error)
	SendRequest(id string, payload transport.MessagePayload) (string, error)
}
