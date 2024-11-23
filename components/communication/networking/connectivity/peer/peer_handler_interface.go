package peer

import (
	"net"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
)

type PeerHandlerInterface interface {
	AddPeer(conn net.Conn, component transport.ComponentID) *Peer
	RemovePeer(addr net.Addr)
	FindPeerByComponentID(id string) (*Peer, error)
	ClosePeerConnections()
	GetActiveConnectionsNumber() int
}
