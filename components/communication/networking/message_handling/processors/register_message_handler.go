package processors

import (
	"fmt"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/message_handling"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
)

type RegisterMessageHandler struct {
	Logger            *logging.HoornLogger
	Listener          message_handling.ListenerInterface
	RegisterComponent func(message transport.Message) (transport.ComponentID, error)
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
	component, err := r.RegisterComponent(message)

	if err != nil {
		r.Logger.Warn(fmt.Sprintf("Failed to register component '%s@%s': %s", componentID.Title, componentID.Version, err.Error()), false, shared.MainComponentName)

		_, err1 := r.Listener.SendResponse(*message.SenderID, []byte(shared.InvalidRequestResponsePayload), *message.UniqueID, false)
		if err1 != nil {
			r.Logger.Error(fmt.Sprintf("Failed to send response to '%s@%s': %s", component.Title, component.Version, err.Error()), false, shared.MainComponentName)
			return err1
		}

		return err
	}

	_, err = r.Listener.SendResponse(*message.SenderID, []byte(shared.RegisterSuccessResponsePayload), *message.UniqueID, false)
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
