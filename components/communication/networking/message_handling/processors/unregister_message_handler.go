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
	GetComponentByID   func(id string) (transport.ComponentID, error)
}

func (u UnregisterMessageHandler) ProcessMessage(message transport.Message) error {
	componentID, err := u.GetComponentByID(*message.RequesterID)

	if err != nil {
		u.Logger.Error(fmt.Sprintf("Failed to find component with ID '%s': %s", *message.RequesterID, err.Error()), false, shared.MainComponentName)
		return err // Returning error here will prevent the message from being processed further. Returning nil will allow the message to be processed.
	}

	u.Logger.Info(fmt.Sprintf("Received unregister message from '%s@%s'", componentID.Title, componentID.Version), false, shared.RoutingComponentName)
	u.ComponentRegistrar.RemoveRegisteredComponent(*message.RequesterID)
	return nil
}
