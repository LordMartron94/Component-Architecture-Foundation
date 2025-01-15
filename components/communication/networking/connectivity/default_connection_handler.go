package connectivity

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/coding"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/connectivity/communication"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/connectivity/peer"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/message_handling"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/scanning"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
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

	listeners map[*peer.Peer]chan struct{}
}

func NewDefaultConnectionHandler(logger *logging.HoornLogger, listenAddress string, messageCoder coding.MessageCoderInterface, peerHandler peer.PeerHandlerInterface, messageUtility message_handling.MessageUtilityInterface, communicationHandler communication.CommunicationHandlerInterface, messageChannel chan transport.Message, shutdownChan chan struct{}) *DefaultConnectionHandler {
	stopListeningChan := make(chan struct{})

	handler := &DefaultConnectionHandler{
		Logger:        logger,
		ListenAddress: listenAddress,

		MessageCoder:         messageCoder,
		PeerHandler:          peerHandler,
		MessageUtility:       messageUtility,
		CommunicationHandler: communicationHandler,

		MessageChannel: messageChannel,
		shutdownChan:   shutdownChan,

		listeners: make(map[*peer.Peer]chan struct{}),
	}

	dataListener := DataListener{
		Logger:         logger,
		PeerHandler:    peerHandler,
		MessageUtility: messageUtility,
		MessageChannel: messageChannel,
		sendResponse:   handler.sendResponse,
		shutdownChan:   stopListeningChan,
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
			close(d.DataListener.shutdownChan)
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

func (d *DefaultConnectionHandler) StopConnection(peer *peer.Peer) error {
	channelToClose := d.listeners[peer]

	if channelToClose != nil {
		close(channelToClose)
	}

	delete(d.listeners, peer)
	return nil
}

func (d *DefaultConnectionHandler) handleConnection(conn net.Conn) {
	scanner := scanning.NewScanner(conn, shared.EndOfMessageToken, d.Logger)

	shutdownChan := make(chan struct{})
	scanned, _ := scanner.Scan(shutdownChan)

	if !scanned {
		d.Logger.Error(fmt.Sprintf("Failed to read registration data from connection"), false, shared.NetworkingComponentName)
		conn.Close()
		return
	}

	//d.Logger.Debug("Scanned!", false, shared.NetworkingComponentName)

	data := scanner.Bytes()
	//dataString := scanner.Text()

	//d.Logger.Debug(fmt.Sprintf("Content: %s", dataString), false, shared.NetworkingComponentName)

	decodedData, err := d.MessageUtility.DecodeMessage(data)

	if err != nil {
		if decodedData.SenderID == nil {
			d.Logger.Warn(fmt.Sprintf("Cannot send response because requester id is missing: %s", err.Error()), false, shared.NetworkingComponentName)
			conn.Close()
			return
		}

		var missingValueError coding.MissingValueError
		if errors.As(err, &missingValueError) {
			componentID, _ := decodedData.GetComponentIDFromRegistrationMessage(d.Logger)

			peer := d.PeerHandler.AddPeer(conn, *componentID) // Necessary to add peer before sending response
			d.sendResponse(*decodedData.SenderID, []byte(shared.InvalidRequestResponsePayload), "__not a uuid because uuid field can be missing__")
			d.PeerHandler.RemovePeer(peer.Address)
		}
		conn.Close()
		return
	}

	actionRequested := decodedData.Payload.Action

	if actionRequested != "register" {
		d.Logger.Error(fmt.Sprintf("Invalid first action requested: '%s'", actionRequested), false, shared.NetworkingComponentName)
		d.sendResponse(*decodedData.SenderID, []byte(shared.InvalidFirstActionPayload), *decodedData.UniqueID)

		conn.Close()
		return
	}

	//d.Logger.Debug(fmt.Sprintf("Pushing data to channel: %s", decodedData.Payload), false, shared.NetworkingComponentName)
	d.MessageChannel <- decodedData

	componentID, err := decodedData.GetComponentIDFromRegistrationMessage(d.Logger)

	if componentID == nil {
		d.Logger.Error("Failed to get component ID from message during registration; closing connection", false, shared.NetworkingComponentName)

		payload, err1 := transport.MessagePayloadFromBytes([]byte(shared.InvalidRequestResponsePayload))
		payload.Args[0].Value = payload.Args[0].Value + fmt.Sprintf("; something went wrong with component ID extraction: %s", err)

		payloadBytes, err2 := transport.MessagePayloadToBytes(payload)

		if err1 != nil || err2 != nil {
			d.Logger.Error(fmt.Sprintf("Failed to create message payload: '%s/%s'", err1.Error(), err2.Error()), false, shared.NetworkingComponentName)
			time.Sleep(time.Second)
			conn.Close()
			return
		}

		d.sendRawResponse(payloadBytes, conn, *decodedData.UniqueID)
		time.Sleep(time.Second)
		conn.Close()
		return
	}

	peer := d.PeerHandler.AddPeer(conn, *componentID)

	d.Logger.Info(fmt.Sprintf("New connection from '%s'", conn.RemoteAddr()), false, shared.NetworkingComponentName)

	stopListeningChan := make(chan struct{})
	d.listeners[peer] = stopListeningChan

	go d.DataListener.ListenForData(*peer, scanner, stopListeningChan)
}

func (d *DefaultConnectionHandler) sendResponse(component string, payload []byte, targetUUID string) {
	_, err := d.CommunicationHandler.SendResponse(component, payload, targetUUID)

	if err != nil {
		d.Logger.Warn(fmt.Sprintf("Failed to send response: '%s'", err.Error()), false, shared.NetworkingComponentName)
	}
}

func (d *DefaultConnectionHandler) sendRawResponse(payload []byte, connection net.Conn, targetUUID string) {
	parsedPayload, err := transport.MessagePayloadFromBytes(payload)

	if err != nil {
		d.Logger.Error(fmt.Sprintf("Failed to create message payload: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return
	}

	message := d.MessageUtility.CreateMessage(parsedPayload, targetUUID)
	encodedMessage, err := d.MessageCoder.Encode(message)
	encodedMessage = append(encodedMessage, []byte(shared.EndOfMessageToken)...)
	connection.Write(encodedMessage)
}
