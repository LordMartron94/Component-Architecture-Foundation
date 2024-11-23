package processors

import "github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"

type MessageProcessorInterface interface {
	ProcessMessage(message transport.Message) error
}
