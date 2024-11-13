package peer

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type PeerHandler struct {
	Logger            *logging.HoornLogger
	activeConnections []*Peer
	allPeers          []*Peer
	lastAliveMessages map[string]time.Time
	lastRemovedPeerAt map[string]time.Time
	mu                sync.Mutex
}

func NewPeerHandler(logger *logging.HoornLogger) *PeerHandler {
	peerHandler := &PeerHandler{
		Logger:            logger,
		activeConnections: make([]*Peer, 0),
		allPeers:          make([]*Peer, 0),
		lastAliveMessages: make(map[string]time.Time),
		lastRemovedPeerAt: make(map[string]time.Time),
		mu:                sync.Mutex{},
	}

	go peerHandler.ClassifyConnections(30 * time.Second)

	return peerHandler
}

func (p *PeerHandler) ClassifyConnections(keepAliveInterval time.Duration) {
	ticker := time.NewTicker(keepAliveInterval)
	defer ticker.Stop()

	for range ticker.C {
		p.detectInactiveConnections(keepAliveInterval)

		// Remove inactive peers from the activeConnections list and put active peers to the activeConnections list
		for _, peer := range p.allPeers {
			if peer.Active && !p.peerInActiveConnections(peer) {
				p.activeConnections = append(p.activeConnections, peer)
			}
			if !peer.Active && p.peerInActiveConnections(peer) {
				p.Logger.Debug(fmt.Sprintf("Removing inactive peer from active connections: %v", peer.Address.String()), false, shared.NetworkingComponentName)
				p.RemovePeerFromActiveConnections(peer)
			}
		}
	}
}

func (p *PeerHandler) RemovePeerFromActiveConnections(peer *Peer) {
	p.mu.Lock()         // Acquire the lock
	defer p.mu.Unlock() // Release the lock when the function exits

	for i, activePeer := range p.activeConnections {
		p.Logger.Debug(fmt.Sprintf("Checking address '%s' against '%s'", activePeer.Address.String(), peer.Address.String()), false, shared.NetworkingComponentName)

		if activePeer.Address.String() == peer.Address.String() {
			p.Logger.Debug(fmt.Sprintf("Removing peer from active connections: %v", peer.Address.String()), false, shared.NetworkingComponentName)
			p.activeConnections = append(p.activeConnections[:i], p.activeConnections[i+1:]...)
			break
		}
	}
}

func (p *PeerHandler) peerInActiveConnections(peer *Peer) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, activePeer := range p.activeConnections {
		if activePeer.Address.String() == peer.Address.String() {
			return true
		}
	}

	return false
}

func (p *PeerHandler) detectInactiveConnections(keepAliveInterval time.Duration) {
	currentTime := time.Now()

	// Detect dead connections
	for i, peer := range p.activeConnections {
		if currentTime.Sub(p.lastAliveMessages[peer.Address.String()]) > keepAliveInterval+(15*time.Second) {
			p.Logger.Warn(fmt.Sprintf("No successful keep-alive message from peer in required interval; removing dead connection from peer: '%v'", peer.Address.String()), false, shared.NetworkingComponentName)

			p.activeConnections[i].Active = false
			p.lastRemovedPeerAt[peer.Address.String()] = currentTime
		}
	}
}

func (p *PeerHandler) AddPeer(conn net.Conn, component transport.ComponentID) *Peer {
	p.mu.Lock()

	peer := &Peer{
		Address:             conn.RemoteAddr(),
		Connection:          conn,
		Outbound:            false,
		AssociatedComponent: component,
		Active:              true,
	}

	p.activeConnections = append(p.activeConnections, peer)
	p.allPeers = append(p.allPeers, peer)
	p.lastAliveMessages[peer.Address.String()] = time.Now()

	p.mu.Unlock()

	p.Logger.Debug(fmt.Sprintf("Adding new peer: %v", conn.RemoteAddr().String()), false, shared.NetworkingComponentName)
	p.Logger.Info(fmt.Sprintf("Router is running. Active peers: %d", p.GetActiveConnectionsNumber()), false, shared.NetworkingComponentName)
	return peer
}

func (p *PeerHandler) GetActiveConnectionsNumber() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.activeConnections)
}

func (p *PeerHandler) RemovePeer(addr net.Addr) {
	p.mu.Lock()

	p.Logger.Debug(fmt.Sprintf("Removing peer by address: %v", addr.String()), false, shared.NetworkingComponentName)

	for i, peer := range p.activeConnections {
		if peer.Address.String() == addr.String() {
			p.activeConnections = append(p.activeConnections[:i], p.activeConnections[i+1:]...)
			p.allPeers = append(p.allPeers[:i], p.allPeers[i+1:]...)
			break
		}
	}

	p.mu.Unlock()

	p.Logger.Info(fmt.Sprintf("Router is running. Active peers: %d", p.GetActiveConnectionsNumber()), false, shared.NetworkingComponentName)
}

func (p *PeerHandler) FindPeerByComponentID(id string) (*Peer, error) {
	for _, peer := range p.allPeers {
		if peer.AssociatedComponent.ComponentUniqueID == id {
			return peer, nil
		}
	}

	return &Peer{}, fmt.Errorf("peer not found by component ID: %v", id)
}

func (p *PeerHandler) ClosePeerConnections() {
	p.Logger.Info("Closing all peer connections", false, shared.NetworkingComponentName)

	for _, peer := range p.allPeers {
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

func (p *PeerHandler) KeepAlive(peer *Peer) {
	currentTime := time.Now()
	lastSuccessfulAliveTime := p.lastAliveMessages[peer.Address.String()]

	if currentTime.Sub(lastSuccessfulAliveTime) > 5*time.Minute {
		p.Logger.Debug(fmt.Sprintf("Gotten keep alive peer connection: %v", peer.Address), false, shared.NetworkingComponentName)
	}

	p.lastAliveMessages[peer.Address.String()] = time.Now()
}
