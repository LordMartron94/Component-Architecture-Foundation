package message_handling

import "github.com/component-architecture-foundation/networking/transport"

type MessageProcessorInterface interface {
	ProcessMessage(message transport.Message) error
}
