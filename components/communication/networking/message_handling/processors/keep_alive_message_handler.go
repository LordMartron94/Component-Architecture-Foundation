package processors

import (
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/connectivity/peer"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type KeepAliveMessageHandler struct {
	Logger      *logging.HoornLogger
	PeerHandler *peer.PeerHandler
}

func (k *KeepAliveMessageHandler) ProcessMessage(message transport.Message) error {
	associatedPeer, err := k.PeerHandler.FindPeerByComponentID(*message.SenderID)

	if err != nil {
		k.Logger.Warn(fmt.Sprintf("Failed to find peer by component ID: '%s'", *message.SenderID), false, shared.RoutingComponentName)
		return err
	}

	k.PeerHandler.KeepAlive(associatedPeer)
	return nil
}
