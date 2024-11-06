package networking

import "github.com/component-architecture-foundation/networking/transport"

type NetworkHandlerInterface interface {
	Shutdown() error
	StartListenLoop() error
	SendMessage(Peer, transport.Message) error
	SendRequest(transport.ComponentID, transport.Message) error
}
