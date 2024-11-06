package message_handling

import (
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/coding"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type MessageUtility struct {
	Logger       logging.HoornLogger
	Server       transport.ComponentID
	MessageCoder coding.MessageCoderInterface
}

func (m *MessageUtility) DecodeMessage(data []byte) (transport.Message, error) {
	decodedMessage, err := m.MessageCoder.Decode(data)
	if err != nil {
		m.Logger.Error(fmt.Sprintf("Failed to decode message: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return transport.Message{}, err
	}

	return decodedMessage, nil
}

func (m *MessageUtility) CreateMessage(payload transport.MessagePayload, target transport.ComponentID) transport.Message {
	message := transport.NewMessage(m.Server, target, payload)
	return *message
}
