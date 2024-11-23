package communication

import (
	"fmt"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/coding"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/connectivity/peer"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/message_handling"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
)

type CommunicationHandler struct {
	Logger         *logging.HoornLogger
	PeerHandler    peer.PeerHandlerInterface
	MessageUtility message_handling.MessageUtilityInterface
	MessageCoder   coding.MessageCoderInterface
}

func (c *CommunicationHandler) SendResponse(id string, payload []byte, targetUUID string) (string, error) {
	associatedPeer, err := c.PeerHandler.FindPeerByComponentID(id)
	if err != nil {
		c.Logger.Warn(fmt.Sprintf("Error finding peer by component ID: '%s' | Can't send Response.", err.Error()), false, shared.NetworkingComponentName)
		return "", err
	}

	createdPayload, err := transport.MessagePayloadFromBytes(payload)
	if err != nil {
		c.Logger.Warn(fmt.Sprintf("Error creating message payload: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return "", err
	}

	createdMessage := c.MessageUtility.CreateMessage(createdPayload, targetUUID)
	return *createdMessage.UniqueID, c.sendMessage(*associatedPeer, createdMessage)
}

func (c *CommunicationHandler) SendRequest(id string, payload transport.MessagePayload) (string, error) {
	associatedPeer, err := c.PeerHandler.FindPeerByComponentID(id)
	if err != nil {
		c.Logger.Warn(fmt.Sprintf("Error finding peer by component ID: '%s' | Can't send Request.", err.Error()), false, shared.NetworkingComponentName)
		return "", err
	}

	createdMessage := c.MessageUtility.CreateMessage(payload, "")
	return *createdMessage.UniqueID, c.sendMessage(*associatedPeer, createdMessage)
}

func (c *CommunicationHandler) sendMessage(target peer.Peer, message transport.Message) error {
	encodedMessage, err := c.MessageCoder.Encode(message)

	if err != nil {
		return err
	}

	//if message.Payload.Action != "response" {
	c.Logger.Info(fmt.Sprintf("Sending message to peer: '%s'; '%s'", target.Address, message.Payload.Action), false, shared.NetworkingComponentName)
	//}

	encodedMessage = append(encodedMessage, []byte(shared.EndOfMessageToken)...)

	conn := target.Connection
	_, err = conn.Write(encodedMessage)

	if err != nil {
		return err
	}

	return nil
}
