package peer

import (
	"net"

	"github.com/component-architecture-foundation/networking/transport"
)

type PeerHandlerInterface interface {
	AddPeer(conn net.Conn, component transport.ComponentID) Peer
	RemovePeer(addr net.Addr)
	FindPeerByComponentID(id transport.ComponentID) (Peer, error)
	ClosePeerConnections()
}
