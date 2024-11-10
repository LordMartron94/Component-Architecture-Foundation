package message_handling

import "github.com/component-architecture-foundation/networking/transport"

type ListenerInterface interface {
	SendResponse(requester transport.ComponentID, response []byte, targetUUID string) error
	SendRequest(id transport.ComponentID, payload transport.MessagePayload) error
}
