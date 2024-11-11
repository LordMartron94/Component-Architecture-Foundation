package peer

import (
	"errors"
	"io"
	"net"
	"time"

	"github.com/component-architecture-foundation/networking/transport"
)

type Peer struct {
	Address             net.Addr
	Connection          net.Conn
	Outbound            bool
	AssociatedComponent transport.ComponentID
	Active              bool
}

func (p *Peer) CheckConnection() bool {
	p.Connection.SetReadDeadline(time.Now().Add(1 * time.Second))
	_, err := p.Connection.Read([]byte{})
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false // Connection closed
		}

		var netErr net.Error
		ok := errors.As(err, &netErr)
		if ok && (netErr.Timeout() || netErr.Temporary()) {
			return false // Timeout or temporary error
		}
		return false // Other errors
	}
	return true // Connection seems okay
}
