package handlers

import (
	"sync"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/connectivity"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/connectivity/communication"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/connectivity/peer"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
)

type TCPHandler struct {
	Logger               *logging.HoornLogger
	ConnectionHandler    connectivity.ConnectionHandlerInterface
	CommunicationHandler communication.CommunicationHandlerInterface
	WaitGroup            *sync.WaitGroup

	shutdownChan chan struct{}
}

func NewTCPHandler(logger *logging.HoornLogger, handlerInterface connectivity.ConnectionHandlerInterface, communicationHandlerInterface communication.CommunicationHandlerInterface, shutdownChan chan struct{}, waitgroup *sync.WaitGroup) *TCPHandler {
	waitgroup.Add(1)

	return &TCPHandler{
		Logger:               logger,
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

func (tcp *TCPHandler) SendResponse(id string, payload []byte, targetUUID string, isClientResponse bool) (string, error) {
	payloadObject, _ := transport.MessagePayloadFromBytes(payload)

	isClientResponseString := "false"

	if isClientResponse {
		isClientResponseString = "true"
	}

	payloadObject.Args = append(payloadObject.Args, transport.Argument{
		Type:  "bool",
		Value: isClientResponseString,
	})

	payloadBytes, _ := transport.MessagePayloadToBytes(payloadObject)

	return tcp.CommunicationHandler.SendResponse(id, payloadBytes, targetUUID)
}

func (tcp *TCPHandler) SendRequest(id string, payload transport.MessagePayload) (string, error) {
	return tcp.CommunicationHandler.SendRequest(id, payload)
}

func (tcp *TCPHandler) StopConnection(peer *peer.Peer) error {
	return tcp.ConnectionHandler.StopConnection(peer)
}
