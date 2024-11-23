package connectivity

import (
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/connectivity/peer"
)

type ConnectionHandlerInterface interface {
	StartListenLoop() error
	CloseConnections()
	GetActiveConnectionsNumber() int
	StopConnection(peer *peer.Peer) error
}
