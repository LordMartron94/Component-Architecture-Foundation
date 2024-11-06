package message_handling

import (
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type RegisterMessageHandler struct {
	Logger            logging.HoornLogger
	Listener          ListenerInterface
	RegisterComponent func(message transport.Message)
}

func (r RegisterMessageHandler) ProcessMessage(message transport.Message) error {
	r.Logger.Info(fmt.Sprintf("Received registration request from '%s@%s'", message.Requester.Title, message.Requester.Version), false, shared.MainComponentName)
	r.RegisterComponent(message)

	err := r.Listener.SendResponse(message.Requester, []byte(shared.RegisterSuccessResponsePayload))
	if err != nil {
		r.Logger.Error(fmt.Sprintf("Failed to send response to '%s@%s': %s", message.Requester.Title, message.Requester.Version, err.Error()), false, shared.MainComponentName)
		return err
	}

	return nil
}
