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
	RegisterComponent func(message transport.Message) transport.ComponentID
}

func (r *RegisterMessageHandler) ProcessMessage(message transport.Message) error {
	componentID, err := message.GetComponentIDFromRegistrationMessage(r.Logger)

	if err != nil {
		r.Logger.Error(fmt.Sprintf("Failed to parse component ID from message: %s", err.Error()), false, shared.MainComponentName)
		return err
	}

	r.Logger.Info(fmt.Sprintf("Received registration request from '%s@%s'", componentID.Title, componentID.Version), false, shared.MainComponentName)
	capabilities := r.convertCapabilitiesToString(componentID.Capabilities)
	r.Logger.Debug(fmt.Sprintf("Capabilities: %s", capabilities), false, shared.MainComponentName)
	component := r.RegisterComponent(message)

	_, err = r.Listener.SendResponse(*message.SenderID, []byte(shared.RegisterSuccessResponsePayload), *message.UniqueID)
	if err != nil {
		r.Logger.Error(fmt.Sprintf("Failed to send response to '%s@%s': %s", component.Title, component.Version, err.Error()), false, shared.MainComponentName)
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
