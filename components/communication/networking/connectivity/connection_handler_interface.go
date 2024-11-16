package connectivity

import (
	"github.com/component-architecture-foundation/networking/connectivity/peer"
)

type ConnectionHandlerInterface interface {
	StartListenLoop() error
	CloseConnections()
	GetActiveConnectionsNumber() int
	StopConnection(peer *peer.Peer) error
}
