package communication

import (
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/coding"
	"github.com/component-architecture-foundation/networking/connectivity/peer"
	"github.com/component-architecture-foundation/networking/message_handling"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type CommunicationHandler struct {
	Logger         logging.HoornLogger
	PeerHandler    peer.PeerHandlerInterface
	MessageUtility message_handling.MessageUtilityInterface
	MessageCoder   coding.MessageCoderInterface
}

func (c *CommunicationHandler) SendResponse(id transport.ComponentID, payload []byte) error {
	associatedPeer, err := c.PeerHandler.FindPeerByComponentID(id)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Error finding peer by component ID: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return err
	}

	createdPayload, err := transport.MessagePayloadFromBytes(payload)

	if err != nil {
		c.Logger.Error(fmt.Sprintf("Error creating message payload: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return err
	}

	createdMessage := c.MessageUtility.CreateMessage(createdPayload, id)
	return c.sendMessage(associatedPeer, createdMessage)
}

func (c *CommunicationHandler) SendRequest(id transport.ComponentID, payload transport.MessagePayload) error {
	associatedPeer, err := c.PeerHandler.FindPeerByComponentID(id)
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Error finding peer by component ID: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return err
	}

	createdMessage := c.MessageUtility.CreateMessage(payload, id)
	return c.sendMessage(associatedPeer, createdMessage)
}

func (c *CommunicationHandler) sendMessage(target peer.Peer, message transport.Message) error {
	encodedMessage, err := c.MessageCoder.Encode(message)

	if err != nil {
		return err
	}

	encodedMessage = append(encodedMessage, []byte(shared.EndOfMessageToken)...)

	conn := target.Connection
	_, err = conn.Write(encodedMessage)

	if err != nil {
		return err
	}

	return nil
}
