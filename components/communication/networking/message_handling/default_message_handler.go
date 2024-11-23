package message_handling

import (
	"fmt"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
)

type DefaultMessageHandler struct {
	Logger                                 *logging.HoornLogger
	Listener                               ListenerInterface
	GetTargetComponentForMessage           func(message transport.Message) (transport.ComponentID, error)
	GetComponentByID                       func(componentID string) (transport.ComponentID, error)
	GetExpectedClientResponses             func(message transport.Message) int
	unrepliedMessageWithRequester          map[string]transport.ComponentID
	mapSystemMessageIDToComponentMessageID map[string]string
	initialized                            bool
}

func (d *DefaultMessageHandler) Init() {
	d.unrepliedMessageWithRequester = make(map[string]transport.ComponentID)
	d.mapSystemMessageIDToComponentMessageID = make(map[string]string)
	d.initialized = true
}

func (d *DefaultMessageHandler) ProcessMessage(message transport.Message) error {
	if !d.initialized {
		d.Init()
	}

	d.Logger.Debug(fmt.Sprintf("Received message: %s", message.ToString()), false, shared.MainComponentName)

	if message.TargetID == "" {
		return d.processRequest(message)
	} else {
		return d.processResponse(message)
	}
}

func (d *DefaultMessageHandler) processRequest(message transport.Message) error {
	originalComponent, err := d.GetComponentByID(*message.SenderID)
	targetComponent, err := d.GetTargetComponentForMessage(message)
	expectedClientResponses := d.GetExpectedClientResponses(message)

	if err != nil {
		d.Logger.Error(fmt.Sprintf("Failed to find target component for action '%s' requested by '%s@%s' message: %s", message.Payload.Action, originalComponent.Title, originalComponent.Version, err.Error()), false, shared.MainComponentName)
		_, err := d.Listener.SendResponse(*message.SenderID, []byte(shared.NoMatchFoundResponsePayload), *message.UniqueID, false)

		if err != nil {
			d.Logger.Error(fmt.Sprintf("Failed to send response to '%s@%s': %s", originalComponent.Title, originalComponent.Version, err.Error()), false, shared.MainComponentName)
			return err
		}

		return nil
	}

	payload, err := transport.MessagePayloadFromBytes([]byte(shared.DefaultSuccessResponsePayload))
	payload.Args = append(payload.Args, transport.Argument{
		Type:  "int",
		Value: fmt.Sprintf("%d", expectedClientResponses),
	})

	bytes, err := transport.MessagePayloadToBytes(payload)

	if err != nil {
		d.Logger.Warn(fmt.Sprintf("Failed to create message payload: '%s'", err.Error()), false, shared.MainComponentName)
		return err
	}

	_, err = d.Listener.SendResponse(*message.SenderID, bytes, *message.UniqueID, false)

	uuid, err := d.Listener.SendRequest(targetComponent.ComponentUniqueID, *message.Payload)

	if err != nil {
		d.Logger.Error(fmt.Sprintf("Failed to send request to '%s@%s': %s", targetComponent.Title, targetComponent.Version, err.Error()), false, shared.MainComponentName)
		return err
	}

	d.mapSystemMessageIDToComponentMessageID[uuid] = *message.UniqueID
	d.unrepliedMessageWithRequester[uuid] = originalComponent

	return nil
}

func (d *DefaultMessageHandler) processResponse(message transport.Message) error {
	bytes, err := transport.MessagePayloadToBytes(*message.Payload)

	targetComponent := d.unrepliedMessageWithRequester[message.TargetID]
	targetID := d.mapSystemMessageIDToComponentMessageID[message.TargetID]

	_, err = d.Listener.SendResponse(targetComponent.ComponentUniqueID, bytes, targetID, true)

	if err != nil {
		d.Logger.Error(fmt.Sprintf("Failed to send request to '%s@%s': %s", targetComponent.Title, targetComponent.Version, err.Error()), false, shared.MainComponentName)
		return err
	}

	delete(d.unrepliedMessageWithRequester, targetID)
	delete(d.mapSystemMessageIDToComponentMessageID, targetID)

	return nil
}
