package processors

import (
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/component_registration"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type UnregisterMessageHandler struct {
	Logger             *logging.HoornLogger
	ComponentRegistrar *component_registration.ComponentRegistrar
}

func (u UnregisterMessageHandler) ProcessMessage(message transport.Message) error {
	u.Logger.Info(fmt.Sprintf("Received unregister message from '%s@%s'", message.Requester.Title, message.Requester.Version), false, shared.RoutingComponentName)
	u.ComponentRegistrar.RemoveRegisteredComponent(*message.Requester)
	return nil
}
