package message_handling

import "github.com/component-architecture-foundation/networking/transport"

type MessageUtilityInterface interface {
	DecodeMessage(data []byte) (transport.Message, error)
	CreateMessage(payload transport.MessagePayload, target transport.ComponentID) transport.Message
}
