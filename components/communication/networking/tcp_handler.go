package networking

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/authentication"
	"github.com/component-architecture-foundation/networking/coding"
	"github.com/component-architecture-foundation/networking/scanning"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type TCPHandler struct {
	Logger         logging.HoornLogger
	ListenAddress  string
	Authenticator  authentication.AuthenticatorInterface
	MessageChannel chan transport.Message
	MessageCoder   coding.MessageCoderInterface

	activeConnections []Peer
	shutdownChan      chan struct{}
	server            transport.ComponentID
}

func NewTCPHandler(logger logging.HoornLogger, address string, authenticator authentication.AuthenticatorInterface, messageChannel chan transport.Message, coder coding.MessageCoderInterface) *TCPHandler {
	return &TCPHandler{
		Logger:            logger,
		ListenAddress:     address,
		Authenticator:     authenticator,
		MessageChannel:    messageChannel,
		MessageCoder:      coder,
		activeConnections: make([]Peer, 0),
		shutdownChan:      make(chan struct{}),
		server: transport.ComponentID{
			Title:        "Middleman",
			Version:      "1.0.0",
			Capabilities: make([]transport.Capability, 0),
		},
	}
}

func (tcp *TCPHandler) Shutdown() error {
	tcp.Logger.Info("Shutting down Listener", false, shared.NetworkingComponentName)
	close(tcp.MessageChannel)
	close(tcp.shutdownChan)
	return nil
}

// StartListenLoop starts a loop that listens for inbound connections.
func (tcp *TCPHandler) StartListenLoop() error {
	listener, err := net.Listen("tcp", tcp.ListenAddress)

	if err != nil {
		tcp.Logger.Error(fmt.Sprintf("There was an error listening on address: '%s', error: '%s'", tcp.ListenAddress, err.Error()), false, shared.NetworkingComponentName)
		return err
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			tcp.Logger.Error(fmt.Sprintf("There was an error closing the Listener: '%s'", err.Error()), false, shared.NetworkingComponentName)
		}
	}(listener)

	tcp.Logger.Info(fmt.Sprintf("Listening for connections on address: '%s'", tcp.ListenAddress), false, shared.NetworkingComponentName)

	for {
		select {
		case <-tcp.shutdownChan:
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				tcp.Logger.Error(fmt.Sprintf("Failed to accept connection: '%s'", err.Error()), false, shared.NetworkingComponentName)
				continue
			}

			go tcp.handleConnection(conn)
		}
	}
}

func (tcp *TCPHandler) handleConnection(conn net.Conn) {
	scanner := scanning.NewScanner(conn, shared.EndOfMessageToken)
	scanned, _ := scanner.Scan()

	if !scanned {
		tcp.Logger.Error(fmt.Sprintf("Failed to read registration data from connection"), false, shared.NetworkingComponentName)
		conn.Close()
		return
	}

	tcp.Logger.Debug("Scanned!", false, shared.NetworkingComponentName)

	data := scanner.Bytes()
	dataString := scanner.Text()

	tcp.Logger.Debug(fmt.Sprintf("Content: %s", dataString), false, shared.NetworkingComponentName)

	decodedData, err := tcp.decodeReadData(data)

	if err != nil {
		tcp.Logger.Error(fmt.Sprintf("Failed to decode registration data: '%s'", err.Error()), false, shared.NetworkingComponentName)
		conn.Close()
		return
	}

	actionRequested := decodedData.Payload.Action

	if actionRequested != "register" {
		tcp.Logger.Error(fmt.Sprintf("Invalid first action requested: '%s'", actionRequested), false, shared.NetworkingComponentName)
		conn.Close()
		return
	}

	tcp.Logger.Debug(fmt.Sprintf("Pushing data to channel: %s", decodedData.Payload), false, shared.NetworkingComponentName)
	tcp.MessageChannel <- decodedData

	peer := Peer{
		Address:             conn.RemoteAddr(),
		Connection:          conn,
		Outbound:            false,
		AssociatedComponent: decodedData.Requester,
	}

	tcp.activeConnections = append(tcp.activeConnections, peer)

	tcp.Logger.Info(fmt.Sprintf("New connection from '%s'", conn.RemoteAddr()), false, shared.NetworkingComponentName)

	go tcp.listenForData(conn, peer, scanner)
}

func (tcp *TCPHandler) listenForData(conn net.Conn, peer Peer, scanner *scanning.Scanner) {
	for {
		select {
		case <-tcp.shutdownChan:
			err := conn.Close()
			if err != nil {
				tcp.Logger.Error(fmt.Sprintf("There was an error closing the connection: '%s'", err.Error()), false, shared.NetworkingComponentName)
			}
			return
		default:
			// Set a read deadline to prevent indefinite blocking
			//err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			//if err != nil {
			//	tcp.Logger.Error(fmt.Sprintf("Error setting read deadline: %s", err), false, shared.NetworkingComponentName)
			//	return
			//}

			for {
				scanned, err := scanner.Scan()

				if err != nil {
					//if os.IsTimeout(err) {
					//	tcp.Logger.Debug(fmt.Sprintf("Read deadline exceeded for peer: '%s'", conn.RemoteAddr()), false, shared.NetworkingComponentName)
					//	conn.Close()
					//	tcp.removePeer(peer.Address)
					//	return
					//}

					if err == io.EOF {
						tcp.Logger.Info(fmt.Sprintf("Connection closed by peer: '%s'", conn.RemoteAddr()), false, shared.NetworkingComponentName)
						tcp.removePeer(peer.Address)
						return
					}

					var operr *net.OpError
					if errors.As(err, &operr) {
						if operr.Op == "read" && strings.Contains(operr.Err.Error(), "wsarecv") {
							tcp.Logger.Info(fmt.Sprintf("Connection closed by peer: '%s'", conn.RemoteAddr()), false, shared.NetworkingComponentName)
							tcp.removePeer(peer.Address)
							return
						}
					}

					tcp.Logger.Error(fmt.Sprintf("Error reading data: %s", err), false, shared.NetworkingComponentName)
					continue
				}

				if !scanned {
					continue
				}

				data := scanner.Bytes()

				if err != nil {
					tcp.Logger.Error(fmt.Sprintf("Failed to read data: '%s'", err.Error()), false, shared.NetworkingComponentName)
					continue
				}

				//Reset the deadline if a successful read occurs
				//conn.SetReadDeadline(time.Time{})

				decodedMessage, err := tcp.decodeReadData(data)

				if err != nil {
					err = tcp.sendResponse(decodedMessage.Target, peer, []byte(shared.InvalidRequestResponse))
					if err != nil {
						tcp.Logger.Error(fmt.Sprintf("Failed to send failure response: '%s'", err.Error()), false, shared.NetworkingComponentName)
						continue
					}
				}

				err = tcp.sendResponse(decodedMessage.Target, peer, []byte(shared.DefaultSuccessReponse))
				if err != nil {
					tcp.Logger.Error(fmt.Sprintf("Failed to send success response: '%s'", err.Error()), false, shared.NetworkingComponentName)
					continue
				}

				tcp.MessageChannel <- decodedMessage
			}
		}
	}
}

func (tcp *TCPHandler) decodeReadData(data []byte) (transport.Message, error) {
	decodedMessage, err := tcp.MessageCoder.Decode(data)
	if err != nil {
		tcp.Logger.Error(fmt.Sprintf("Failed to decode message: '%s'", err.Error()), false, shared.NetworkingComponentName)

		return transport.Message{}, err
	}

	return decodedMessage, nil
}

func (tcp *TCPHandler) sendResponse(component transport.ComponentID, peer Peer, response []byte) error {
	codedResponse, err := tcp.MessageCoder.Decode(response)

	if err != nil {
		tcp.Logger.Error(fmt.Sprintf("Failed to encode response: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return err
	}

	codedResponse.Target = component
	codedResponse.TimeSent = time.Now()

	err = tcp.sendMessage(peer, codedResponse)

	if err != nil {
		tcp.Logger.Error(fmt.Sprintf("Failed to send response: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return err
	}

	return nil
}

func (tcp *TCPHandler) removePeer(addr net.Addr) {
	for i, p := range tcp.activeConnections {
		if p.Address.String() == addr.String() {
			tcp.activeConnections = append(tcp.activeConnections[:i], tcp.activeConnections[i+1:]...)
			break //Peer removed, exit the loop
		}
	}
}

func (tcp *TCPHandler) sendMessage(target Peer, message transport.Message) error {
	encodedMessage, err := tcp.MessageCoder.Encode(message)

	if err != nil {
		return err
	}

	encodedMessage = append(encodedMessage, []byte(shared.EndOfMessageToken)...)

	conn := target.Connection
	_, err = conn.Write(encodedMessage)

	if err != nil {
		return err
	}

	return nil
}

func (tcp *TCPHandler) decodeMessage(message []byte, target transport.ComponentID) (transport.Message, error) {
	decodedMessage, err := tcp.MessageCoder.Decode(message)
	if err != nil {
		return transport.Message{}, err
	}

	decodedMessage.TimeSent = time.Now()
	decodedMessage.Target = target

	return decodedMessage, nil
}

func (tcp *TCPHandler) createMessage(payload transport.MessagePayload, target transport.ComponentID) transport.Message {
	message := transport.NewMessage(tcp.server, target, payload)
	return *message
}

func (tcp *TCPHandler) SendResponse(id transport.ComponentID, message []byte) error {
	associatedPeer, err := tcp.findPeerByComponentID(id)
	if err != nil {
		tcp.Logger.Error(fmt.Sprintf("Error finding peer by component ID: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return err
	}

	decodedMessage, err := tcp.decodeMessage(message, id)
	return tcp.sendMessage(associatedPeer, decodedMessage)
}

func (tcp *TCPHandler) SendRequest(id transport.ComponentID, payload transport.MessagePayload) error {
	associatedPeer, err := tcp.findPeerByComponentID(id)
	if err != nil {
		tcp.Logger.Error(fmt.Sprintf("Error finding peer by component ID: '%s'", err.Error()), false, shared.NetworkingComponentName)
		return err
	}

	createdMessage := tcp.createMessage(payload, id)
	return tcp.sendMessage(associatedPeer, createdMessage)
}

func (tcp *TCPHandler) findPeerByComponentID(id transport.ComponentID) (Peer, error) {
	for _, peer := range tcp.activeConnections {
		if peer.AssociatedComponent.Equal(id) {
			return peer, nil
		}
	}

	return Peer{}, fmt.Errorf("peer not found by component ID: %v", id)
}
