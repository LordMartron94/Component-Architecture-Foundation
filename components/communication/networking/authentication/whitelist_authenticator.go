package authentication

import (
	"fmt"
	"net"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/shared"
)

type WhitelistAuthenticator struct {
	Logger             logging.HoornLogger
	allowedConnections []string
}

func (wA *WhitelistAuthenticator) Authenticate(conn net.Conn) error {
	wA.Logger.Debug(fmt.Sprintf("Attempting authentication with '%s' from '%s'", conn.RemoteAddr(), conn.LocalAddr()), false, shared.NetworkingComponentName)
	return nil
}
