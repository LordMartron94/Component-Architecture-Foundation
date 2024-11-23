package message_handling

import "github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"

type ListenerInterface interface {
	SendResponse(requester string, response []byte, targetUUID string, isClientResponse bool) (string, error)
	SendRequest(id string, payload transport.MessagePayload) (string, error)
}
