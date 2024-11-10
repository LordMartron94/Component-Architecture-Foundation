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

func NewPeerHandler(logger *logging.HoornLogger) *PeerHandler {
	return &PeerHandler{
		Logger:            logger,
		activeConnections: make([]Peer, 0),
	}
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

func (p *PeerHandler) GetActiveConnectionsNumber() int {
	return len(p.activeConnections)
}

func (p *PeerHandler) RemovePeer(addr net.Addr) {
	p.Logger.Debug(fmt.Sprintf("Removing peer by address: %v", addr.String()), false, shared.NetworkingComponentName)

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
		// Check if the connection is already closed
		if !peer.CheckConnection() {
			p.Logger.Debug(fmt.Sprintf("Peer connection already closed: %v", peer.Connection.RemoteAddr()), false, shared.NetworkingComponentName)
			continue // Skip to the next connection
		}

		err := peer.Connection.Close()
		if err != nil {
			p.Logger.Error(fmt.Sprintf("Error closing peer connection: %v", err), false, shared.NetworkingComponentName)
		}
	}
}
