package connectivity

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/connectivity/peer"
	"github.com/component-architecture-foundation/networking/message_handling"
	"github.com/component-architecture-foundation/networking/scanning"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type DataListener struct {
	Logger         *logging.HoornLogger
	PeerHandler    peer.PeerHandlerInterface
	MessageUtility message_handling.MessageUtilityInterface
	MessageChannel chan transport.Message

	sendResponse func(component transport.ComponentID, payload []byte, targetUUID string)
	shutdownChan chan struct{}
}

func (d *DataListener) ListenForData(conn net.Conn, peer peer.Peer, scanner *scanning.Scanner) {
	for {
		select {
		case <-d.shutdownChan:
			err := conn.Close()
			if err != nil {
				d.Logger.Error(fmt.Sprintf("There was an error closing the connection: '%s'", err.Error()), false, shared.NetworkingComponentName)
			}
			return
		default:
			for {
				scanned, err := scanner.Scan()

				if err != nil {
					if err == io.EOF {
						d.Logger.Info(fmt.Sprintf("Connection closed by peer: '%s'", conn.RemoteAddr()), false, shared.NetworkingComponentName)
						d.PeerHandler.RemovePeer(peer.Address)
						return
					}

					var operr *net.OpError
					if errors.As(err, &operr) {
						if operr.Op == "read" && strings.Contains(operr.Err.Error(), "wsarecv") {
							d.Logger.Info(fmt.Sprintf("Connection closed by peer: '%s'", conn.RemoteAddr()), false, shared.NetworkingComponentName)
							d.PeerHandler.RemovePeer(peer.Address)
							return
						}
					}

					d.Logger.Error(fmt.Sprintf("Error reading data: %s", err), false, shared.NetworkingComponentName)
					continue
				}

				if !scanned {
					continue
				}

				data := scanner.Bytes()

				if err != nil {
					d.Logger.Error(fmt.Sprintf("Failed to read data: '%s'", err.Error()), false, shared.NetworkingComponentName)
					continue
				}

				decodedMessage, err := d.MessageUtility.DecodeMessage(data)

				if err != nil {
					err = d.sendResponse(decodedMessage.Requester, []byte(shared.InvalidRequestResponsePayload))
					if err != nil {
						d.Logger.Warn(fmt.Sprintf("Failed to send failure response: '%s'", err.Error()), false, shared.NetworkingComponentName)
					}
				}

				err = d.sendResponse(decodedMessage.Requester, []byte(shared.DefaultSuccessResponsePayload))
				if err != nil {
					d.Logger.Warn(fmt.Sprintf("Failed to send success response: '%s'", err.Error()), false, shared.NetworkingComponentName)
				}

				d.MessageChannel <- decodedMessage
			}
		}
	}
}
