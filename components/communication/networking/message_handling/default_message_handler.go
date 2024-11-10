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
}

func (d DefaultMessageHandler) ProcessMessage(message transport.Message) error {
	targetComponent, err := d.GetTargetComponentForMessage(message)

	if err != nil {
		d.Logger.Error(fmt.Sprintf("Failed to find target component for message: %s", err.Error()), false, shared.MainComponentName)
		err := d.Listener.SendResponse(*message.Requester, []byte(shared.NoMatchFoundResponsePayload), *message.UniqueID)

		if err != nil {
			d.Logger.Error(fmt.Sprintf("Failed to send response to '%s@%s': %s", message.Requester.Title, message.Requester.Version, err.Error()), false, shared.MainComponentName)
			return err
		}

		return nil
	}

	err = d.Listener.SendRequest(targetComponent, *message.Payload)

	if err != nil {
		d.Logger.Error(fmt.Sprintf("Failed to send request to '%s@%s': %s", targetComponent.Title, targetComponent.Version, err.Error()), false, shared.MainComponentName)
		return err
	}

	return nil
}
