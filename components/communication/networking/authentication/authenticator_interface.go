package authentication

import (
	"net"
)

type AuthenticatorInterface interface {
	Authenticate(conn net.Conn) error
}
