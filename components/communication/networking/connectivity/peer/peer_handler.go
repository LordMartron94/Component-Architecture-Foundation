package peer

import (
	"fmt"
	"net"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
)

type PeerHandler struct {
	Logger            *logging.HoornLogger
	activeConnections []Peer
}

func (p *PeerHandler) AddPeer(conn net.Conn, component transport.ComponentID) Peer {
	peer := Peer{
		Address:             conn.RemoteAddr(),
		Connection:          conn,
		Outbound:            false,
		AssociatedComponent: component,
	}

	p.activeConnections = append(p.activeConnections, peer)
	return peer
}

func (p *PeerHandler) RemovePeer(addr net.Addr) {
	for i, peer := range p.activeConnections {
		if peer.Address.String() == addr.String() {
			p.activeConnections = append(p.activeConnections[:i], p.activeConnections[i+1:]...)
			break
		}
	}
}

func (p *PeerHandler) FindPeerByComponentID(id transport.ComponentID) (Peer, error) {
	for _, peer := range p.activeConnections {
		if peer.AssociatedComponent.Equal(id) {
			return peer, nil
		}
	}

	return Peer{}, fmt.Errorf("peer not found by component ID: %v", id)
}
