package communication

import "github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"

type CommunicationHandlerInterface interface {
	SendResponse(id string, payload []byte, targetUUID string) (string, error)
	SendRequest(id string, payload transport.MessagePayload) (string, error)
}
