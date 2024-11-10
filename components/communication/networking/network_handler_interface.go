package networking

import "github.com/component-architecture-foundation/networking/transport"

type NetworkHandlerInterface interface {
	Shutdown() error
	StartListenLoop() error
	SendResponse(id transport.ComponentID, message []byte, targetUUID string) error
	SendRequest(id transport.ComponentID, payload transport.MessagePayload) error
	GetActiveConnectionsNumber() int
}
