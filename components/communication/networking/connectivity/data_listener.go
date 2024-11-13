package connectivity

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

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

	sendResponse func(component string, payload []byte, targetUUID string)
	shutdownChan chan struct{}
}

func (d *DataListener) ListenForData(peer peer.Peer, scanner *scanning.Scanner) {
	d.Logger.Debug(fmt.Sprintf("Started listening for data from peer: '%s'", peer.Address), false, shared.NetworkingComponentName)
	defer d.Logger.Debug(fmt.Sprintf("Stopped listening for data from peer: '%s'", peer.Address), false, shared.NetworkingComponentName)

	dataChan := make(chan []byte, 1) // Buffered channel
	stopScanning := make(chan struct{})
	var bufferMutex sync.Mutex // Mutex for buffer access

	go func() {
		for {
			select {
			case <-stopScanning:
				break
			default:
				scanned, err := scanner.Scan(stopScanning)
				if err != nil {
					if d.handleScanError(err, peer) {
						close(stopScanning)
					}
					continue
				}

				if !scanned {
					continue
				}

				bufferMutex.Lock()
				dataChan <- scanner.Bytes()
				bufferMutex.Unlock()
			}
		}
	}()

	for {
		select {
		case <-d.shutdownChan:
			d.Logger.Debug("Stopped listening for data because of shutdown signal.", false, shared.NetworkingComponentName)

			// Check if stopScanning is already closed
			if _, ok := <-stopScanning; !ok {
				d.Logger.Debug("StopScanning is already closed.", false, shared.NetworkingComponentName)
				return
			}

			// If stopScanning is not closed, close it and return
			d.Logger.Debug("Closing stopScanning channel.", false, shared.NetworkingComponentName)
			close(stopScanning)
			return
		case data := <-dataChan:
			decodedMessage, err := d.MessageUtility.DecodeMessage(data)
			d.sendHandleResponse(decodedMessage, err)
			d.MessageChannel <- decodedMessage
		}
	}
}

func (d *DataListener) sendHandleResponse(decodedMessage transport.Message, decodeError error) {
	if decodeError != nil {
		d.sendResponse(*decodedMessage.RequesterID, []byte(shared.InvalidRequestResponsePayload), *decodedMessage.UniqueID)
		return
	}

	d.sendResponse(*decodedMessage.RequesterID, []byte(shared.DefaultSuccessResponsePayload), *decodedMessage.UniqueID)

	return
}

// handleScanError handles errors during scanning and logs them. It also removes the peer from the peer handler if the connection is closed.
// Returns true if the scanning should stop for this connection.
func (d *DataListener) handleScanError(err error, peer peer.Peer) bool {
	if err == io.EOF {
		d.Logger.Info(fmt.Sprintf("Connection closed by peer: '%s'", peer.Address), false, shared.NetworkingComponentName)
		d.PeerHandler.RemovePeer(peer.Address)
		return true
	}

	var operr *net.OpError
	if errors.As(err, &operr) {
		if operr.Op == "read" && strings.Contains(operr.Err.Error(), "wsarecv") {
			d.Logger.Info(fmt.Sprintf("Connection closed by peer: '%s'", peer.Address), false, shared.NetworkingComponentName)
			d.PeerHandler.RemovePeer(peer.Address)
			return true
		}
	}

	d.Logger.Error(fmt.Sprintf("Error reading data: %s", err), false, shared.NetworkingComponentName)
	return false
}
