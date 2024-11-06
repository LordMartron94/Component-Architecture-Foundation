package handlers

import (
	"net"

	"github.com/component-architecture-foundation/networking/transport"
)

type Peer struct {
	Address             net.Addr
	Connection          net.Conn
	Outbound            bool
	AssociatedComponent transport.ComponentID
}
