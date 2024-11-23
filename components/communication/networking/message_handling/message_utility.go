package message_handling

import (
	"fmt"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/coding"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
)

type MessageUtility struct {
	Logger       *logging.HoornLogger
	Server       transport.ComponentID
	MessageCoder coding.MessageCoderInterface
}

// DecodeMessage decodes a given byte slice into a transport.Message.
// If the decoding fails, an error is returned.
// It does return the decoded message.
func (m *MessageUtility) DecodeMessage(data []byte) (transport.Message, error) {
	decodedMessage, err := m.MessageCoder.Decode(data)
	if err != nil {
		m.Logger.Error(fmt.Sprintf("Failed to decode message: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return decodedMessage, err
	}

	return decodedMessage, nil
}

func (m *MessageUtility) CreateMessage(payload transport.MessagePayload, targetID string) transport.Message {
	message := transport.NewMessage(m.Server, payload, targetID)
	return *message
}
