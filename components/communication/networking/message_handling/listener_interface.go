package message_handling

import "github.com/component-architecture-foundation/networking/transport"

type ListenerInterface interface {
	SendResponse(requester string, response []byte, targetUUID string) error
	SendRequest(id string, payload transport.MessagePayload) error
}
