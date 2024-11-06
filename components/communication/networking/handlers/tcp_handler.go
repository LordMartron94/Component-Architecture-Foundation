package handlers

import (
	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type TCPHandler struct {
	Logger               logging.HoornLogger
	MessageChannel       chan transport.Message
	ConnectionHandler    ConnectionHandlerInterface
	CommunicationHandler CommunicationHandlerInterface

	shutdownChan chan struct{}
}

func NewTCPHandler(logger logging.HoornLogger, messageChannel chan transport.Message, handlerInterface ConnectionHandlerInterface, communicationHandlerInterface CommunicationHandlerInterface, shutdownChan chan struct{}) *TCPHandler {
	return &TCPHandler{
		Logger:               logger,
		MessageChannel:       messageChannel,
		ConnectionHandler:    handlerInterface,
		CommunicationHandler: communicationHandlerInterface,
		shutdownChan:         shutdownChan,
	}
}

func (tcp *TCPHandler) Shutdown() error {
	tcp.Logger.Info("Shutting down Listener", false, shared.NetworkingComponentName)
	close(tcp.MessageChannel)
	close(tcp.shutdownChan)
	return nil
}

func (tcp *TCPHandler) StartListenLoop() error {
	return tcp.ConnectionHandler.StartListenLoop()
}

func (tcp *TCPHandler) SendResponse(id transport.ComponentID, payload []byte) error {
	return tcp.CommunicationHandler.SendResponse(id, payload)
}

func (tcp *TCPHandler) SendRequest(id transport.ComponentID, payload transport.MessagePayload) error {
	return tcp.CommunicationHandler.SendRequest(id, payload)
}
