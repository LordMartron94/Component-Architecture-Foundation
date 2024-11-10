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

func (d *DataListener) ListenForData(peer peer.Peer, scanner *scanning.Scanner) {
	dataChan := make(chan []byte)
	go func() {
		for {
			if !peer.CheckConnection() {
				continue
			}

			scanned, err := scanner.Scan(d.shutdownChan)

			if err != nil {
				d.handleScanError(err, peer)
				continue
			}

			if !scanned {
				continue
			}

			dataChan <- scanner.Bytes()
		}
	}()

	for {
		select {
		case <-d.shutdownChan:
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
		d.sendResponse(*decodedMessage.Requester, []byte(shared.InvalidRequestResponsePayload), *decodedMessage.UniqueID)
		return
	}

	d.sendResponse(*decodedMessage.Requester, []byte(shared.DefaultSuccessResponsePayload), *decodedMessage.UniqueID)

	return
}

func (d *DataListener) handleScanError(err error, peer peer.Peer) {
	if err == io.EOF {
		d.Logger.Info(fmt.Sprintf("Connection closed by peer: '%s'", peer.Address), false, shared.NetworkingComponentName)
		d.PeerHandler.RemovePeer(peer.Address)
		return
	}

	var operr *net.OpError
	if errors.As(err, &operr) {
		if operr.Op == "read" && strings.Contains(operr.Err.Error(), "wsarecv") {
			d.Logger.Info(fmt.Sprintf("Connection closed by peer: '%s'", peer.Address), false, shared.NetworkingComponentName)
			d.PeerHandler.RemovePeer(peer.Address)
			return
		}
	}

	d.Logger.Error(fmt.Sprintf("Error reading data: %s", err), false, shared.NetworkingComponentName)
}
