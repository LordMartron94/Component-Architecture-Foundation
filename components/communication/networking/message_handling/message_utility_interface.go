package message_handling

import "github.com/component-architecture-foundation/networking/transport"

type MessageUtilityInterface interface {
	// DecodeMessage decodes a bytes json message object to an actual transport.Message object.
	// Returns error if there are missing values for pointers.
	//But still returns the decoded message.
	DecodeMessage(data []byte) (transport.Message, error)
	CreateMessage(payload transport.MessagePayload) transport.Message
}
