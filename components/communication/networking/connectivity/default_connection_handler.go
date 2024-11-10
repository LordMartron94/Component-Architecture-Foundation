package connectivity

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/coding"
	"github.com/component-architecture-foundation/networking/connectivity/communication"
	"github.com/component-architecture-foundation/networking/connectivity/peer"
	"github.com/component-architecture-foundation/networking/message_handling"
	"github.com/component-architecture-foundation/networking/scanning"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type DefaultConnectionHandler struct {
	Logger        *logging.HoornLogger
	ListenAddress string

	MessageCoder         coding.MessageCoderInterface
	PeerHandler          peer.PeerHandlerInterface
	MessageUtility       message_handling.MessageUtilityInterface
	CommunicationHandler communication.CommunicationHandlerInterface
	DataListener         DataListener

	MessageChannel chan transport.Message
	shutdownChan   chan struct{}
}

func NewDefaultConnectionHandler(logger *logging.HoornLogger, listenAddress string, messageCoder coding.MessageCoderInterface, peerHandler peer.PeerHandlerInterface, messageUtility message_handling.MessageUtilityInterface, communicationHandler communication.CommunicationHandlerInterface, messageChannel chan transport.Message, shutdownChan chan struct{}) *DefaultConnectionHandler {
	handler := &DefaultConnectionHandler{
		Logger:        logger,
		ListenAddress: listenAddress,

		MessageCoder:         messageCoder,
		PeerHandler:          peerHandler,
		MessageUtility:       messageUtility,
		CommunicationHandler: communicationHandler,

		MessageChannel: messageChannel,
		shutdownChan:   shutdownChan,
	}

	dataListener := DataListener{
		Logger:         logger,
		PeerHandler:    peerHandler,
		MessageUtility: messageUtility,
		MessageChannel: messageChannel,
		sendResponse:   handler.sendResponse,
		shutdownChan:   shutdownChan,
	}

	handler.DataListener = dataListener

	return handler
}

func (d *DefaultConnectionHandler) StartListenLoop() error {
	listener, err := net.Listen("tcp", d.ListenAddress)

	if err != nil {
		d.Logger.Error(fmt.Sprintf("There was an error listening on address: '%s', error: '%s'", d.ListenAddress, err.Error()), false, shared.NetworkingComponentName)
		return err
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			d.Logger.Error(fmt.Sprintf("There was an error closing the Listener: '%s'", err.Error()), false, shared.NetworkingComponentName)
		}
	}(listener)

	d.Logger.Info(fmt.Sprintf("Listening for connections on address: '%s'", d.ListenAddress), false, shared.NetworkingComponentName)

	for {
		select {
		case <-d.shutdownChan:
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				d.Logger.Error(fmt.Sprintf("Failed to accept connection: '%s'", err.Error()), false, shared.NetworkingComponentName)
				continue
			}

			go d.handleConnection(conn)
		}
	}
}

func (d *DefaultConnectionHandler) GetActiveConnectionsNumber() int {
	return d.PeerHandler.GetActiveConnectionsNumber()
}

func (d *DefaultConnectionHandler) CloseConnections() {
	d.PeerHandler.ClosePeerConnections()
}

func (d *DefaultConnectionHandler) handleConnection(conn net.Conn) {
	scanner := scanning.NewScanner(conn, shared.EndOfMessageToken)

	shutdownChan := make(chan struct{})
	scanned, _ := scanner.Scan(shutdownChan)

	if !scanned {
		d.Logger.Error(fmt.Sprintf("Failed to read registration data from connection"), false, shared.NetworkingComponentName)
		conn.Close()
		return
	}

	d.Logger.Debug("Scanned!", false, shared.NetworkingComponentName)

	data := scanner.Bytes()
	dataString := scanner.Text()

	d.Logger.Debug(fmt.Sprintf("Content: %s", dataString), false, shared.NetworkingComponentName)

	decodedData, err := d.MessageUtility.DecodeMessage(data)

	if err != nil {
		if decodedData.RequesterID == nil {
			d.Logger.Warn(fmt.Sprintf("Cannot send response because requester id is missing: %s", err.Error()), false, shared.NetworkingComponentName)
			conn.Close()
			return
		}

		var missingValueError coding.MissingValueError
		if errors.As(err, &missingValueError) {
			componentID, _ := decodedData.GetComponentIDFromRegistrationMessage(d.Logger)

			peer := d.PeerHandler.AddPeer(conn, *componentID) // Necessary to add peer before sending response
			d.sendResponse(*decodedData.RequesterID, []byte(shared.InvalidRequestResponsePayload), "__not a uuid because uuid field can be missing__")
			d.PeerHandler.RemovePeer(peer.Address)
		}
		conn.Close()
		return
	}

	actionRequested := decodedData.Payload.Action

	if actionRequested != "register" {
		d.Logger.Error(fmt.Sprintf("Invalid first action requested: '%s'", actionRequested), false, shared.NetworkingComponentName)
		d.sendResponse(*decodedData.RequesterID, []byte(shared.InvalidFirstActionPayload), *decodedData.UniqueID)

		conn.Close()
		return
	}

	d.Logger.Debug(fmt.Sprintf("Pushing data to channel: %s", decodedData.Payload), false, shared.NetworkingComponentName)
	d.MessageChannel <- decodedData

	componentID, _ := decodedData.GetComponentIDFromRegistrationMessage(d.Logger)

	if componentID == nil {
		d.Logger.Error("Failed to get component ID from message", false, shared.NetworkingComponentName)
		d.sendRawResponse([]byte(shared.InvalidRequestResponsePayload), conn, *decodedData.UniqueID)
		time.Sleep(time.Second)
		conn.Close()
		return
	}

	peer := d.PeerHandler.AddPeer(conn, *componentID)

	d.Logger.Info(fmt.Sprintf("New connection from '%s'", conn.RemoteAddr()), false, shared.NetworkingComponentName)

	go d.DataListener.ListenForData(peer, scanner)
}

func (d *DefaultConnectionHandler) sendResponse(component string, payload []byte, targetUUID string) {
	err := d.CommunicationHandler.SendResponse(component, payload, targetUUID)

	if err != nil {
		d.Logger.Warn(fmt.Sprintf("Failed to send response: '%s'", err.Error()), false, shared.NetworkingComponentName)
	}
}

func (d *DefaultConnectionHandler) sendRawResponse(payload []byte, connection net.Conn, messageUUID string) {
	parsedPayload, err := transport.MessagePayloadFromBytes(payload)

	if err != nil {
		d.Logger.Error(fmt.Sprintf("Failed to create message payload: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return
	}

	message := d.MessageUtility.CreateMessage(parsedPayload)
	message.Payload.Args = append(message.Payload.Args, transport.Argument{
		Type:  "string",
		Value: messageUUID,
	})
	encodedMessage, err := d.MessageCoder.Encode(message)
	encodedMessage = append(encodedMessage, []byte(shared.EndOfMessageToken)...)
	connection.Write(encodedMessage)
}
