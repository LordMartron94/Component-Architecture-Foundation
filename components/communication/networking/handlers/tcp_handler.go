package handlers

import (
	"sync"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/connectivity"
	"github.com/component-architecture-foundation/networking/connectivity/communication"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type TCPHandler struct {
	Logger               *logging.HoornLogger
	MessageChannel       chan transport.Message
	ConnectionHandler    connectivity.ConnectionHandlerInterface
	CommunicationHandler communication.CommunicationHandlerInterface
	WaitGroup            *sync.WaitGroup

	shutdownChan chan struct{}
}

func NewTCPHandler(logger *logging.HoornLogger, messageChannel chan transport.Message, handlerInterface connectivity.ConnectionHandlerInterface, communicationHandlerInterface communication.CommunicationHandlerInterface, shutdownChan chan struct{}, waitgroup *sync.WaitGroup) *TCPHandler {
	return &TCPHandler{
		Logger:               logger,
		MessageChannel:       messageChannel,
		ConnectionHandler:    handlerInterface,
		CommunicationHandler: communicationHandlerInterface,
		shutdownChan:         shutdownChan,
		WaitGroup:            waitgroup,
	}
}

func (tcp *TCPHandler) Shutdown() error {
	tcp.Logger.Info("Shutting down Listener", false, shared.NetworkingComponentName)
	close(tcp.shutdownChan)
	//tcp.ConnectionHandler.CloseConnections()  // Removed because this is for the clients to do.
	tcp.WaitGroup.Done()
	return nil
}

func (tcp *TCPHandler) GetActiveConnectionsNumber() int {
	return tcp.ConnectionHandler.GetActiveConnectionsNumber()
}

func (tcp *TCPHandler) StartListenLoop() error {
	return tcp.ConnectionHandler.StartListenLoop()
}

func (tcp *TCPHandler) SendResponse(id string, payload []byte, targetUUID string) error {
	return tcp.CommunicationHandler.SendResponse(id, payload, targetUUID)
}

func (tcp *TCPHandler) SendRequest(id string, payload transport.MessagePayload) error {
	return tcp.CommunicationHandler.SendRequest(id, payload)
}
