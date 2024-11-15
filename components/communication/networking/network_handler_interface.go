package networking

import "github.com/component-architecture-foundation/networking/transport"

type NetworkHandlerInterface interface {
	Shutdown() error
	StartListenLoop() error
	SendResponse(id string, message []byte, targetUUID string, isClientResponse bool) (string, error)
	SendRequest(id string, payload transport.MessagePayload) (string, error)
	GetActiveConnectionsNumber() int
}
