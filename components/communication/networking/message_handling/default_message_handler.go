package message_handling

import (
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type DefaultMessageHandler struct {
	Logger                       *logging.HoornLogger
	Listener                     ListenerInterface
	GetTargetComponentForMessage func(message transport.Message) (transport.ComponentID, error)
	GetComponentByID             func(componentID string) (transport.ComponentID, error)
}

func (d DefaultMessageHandler) ProcessMessage(message transport.Message) error {
	originalComponent, err := d.GetComponentByID(*message.RequesterID)
	targetComponent, err := d.GetTargetComponentForMessage(message)

	if err != nil {
		d.Logger.Error(fmt.Sprintf("Failed to find target component for message: %s", err.Error()), false, shared.MainComponentName)
		err := d.Listener.SendResponse(*message.RequesterID, []byte(shared.NoMatchFoundResponsePayload), *message.UniqueID)

		if err != nil {
			d.Logger.Error(fmt.Sprintf("Failed to send response to '%s@%s': %s", originalComponent.Title, originalComponent.Version, err.Error()), false, shared.MainComponentName)
			return err
		}

		return nil
	}

	err = d.Listener.SendRequest(targetComponent.ComponentUniqueID, *message.Payload)

	if err != nil {
		d.Logger.Error(fmt.Sprintf("Failed to send request to '%s@%s': %s", targetComponent.Title, targetComponent.Version, err.Error()), false, shared.MainComponentName)
		return err
	}

	return nil
}
