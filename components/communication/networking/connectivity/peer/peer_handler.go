package peer

import (
	"fmt"
	"net"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
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

func (p *PeerHandler) ClosePeerConnections() {
	p.Logger.Info("Closing all peer connections", false, shared.NetworkingComponentName)

	for _, peer := range p.activeConnections {
		err := peer.Connection.Close()
		if err != nil {
			p.Logger.Error(fmt.Sprintf("Error closing peer connection: %v", err), false, shared.NetworkingComponentName)
		}
	}
}
