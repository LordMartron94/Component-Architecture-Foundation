package processors

import (
	"fmt"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/connectivity/peer"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
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
