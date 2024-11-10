package networking

import "github.com/component-architecture-foundation/networking/transport"

type NetworkHandlerInterface interface {
	Shutdown() error
	StartListenLoop() error
	SendResponse(id string, message []byte, targetUUID string) error
	SendRequest(id string, payload transport.MessagePayload) error
	GetActiveConnectionsNumber() int
}
