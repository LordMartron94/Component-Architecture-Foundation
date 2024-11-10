package processors

import (
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/message_handling"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type RegisterMessageHandler struct {
	Logger            *logging.HoornLogger
	Listener          message_handling.ListenerInterface
	RegisterComponent func(message transport.Message)
}

func (r *RegisterMessageHandler) ProcessMessage(message transport.Message) error {
	r.Logger.Info(fmt.Sprintf("Received registration request from '%s@%s'", message.Requester.Title, message.Requester.Version), false, shared.MainComponentName)
	capabilities := r.convertCapabilitiesToString(message.Requester.Capabilities)
	r.Logger.Debug(fmt.Sprintf("Capabilities: %s", capabilities), false, shared.MainComponentName)
	r.RegisterComponent(message)

	err := r.Listener.SendResponse(*message.Requester, []byte(shared.RegisterSuccessResponsePayload), *message.UniqueID)
	if err != nil {
		r.Logger.Error(fmt.Sprintf("Failed to send response to '%s@%s': %s", message.Requester.Title, message.Requester.Version, err.Error()), false, shared.MainComponentName)
		return err
	}

	return nil
}

func (r *RegisterMessageHandler) convertCapabilitiesToString(capabilities []transport.Capability) string {
	var finalString = ""

	for i, capability := range capabilities {
		capabilityString := fmt.Sprintf("Capability [%d] - [Name: %s, numArgs: %d]", i, capability.Name, capability.Signature.NumOfArgs)
		finalString += capabilityString + " | "
	}

	return finalString
}
