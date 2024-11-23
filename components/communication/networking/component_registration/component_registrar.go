package component_registration

import (
	"fmt"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/message_handling"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
)

type ComponentRegistrar struct {
	Logger               *logging.HoornLogger
	registeredComponents []transport.ComponentID
	Listener             message_handling.ListenerInterface
}

func (c *ComponentRegistrar) RegisterComponent(message transport.Message) (transport.ComponentID, error) {
	componentID, err := message.GetComponentIDFromRegistrationMessage(c.Logger)

	if err != nil {
		c.Logger.Error(fmt.Sprintf("Error getting component ID from message: %s (on registration)", err.Error()), false, shared.NetworkingComponentName)
		_, err := c.Listener.SendResponse(shared.ServerUUID, []byte(shared.InvalidRequestResponsePayload), componentID.ComponentUniqueID, false)
		if err != nil {
			c.Logger.Error(fmt.Sprintf("Failed to send response to '%s@%s': %s", componentID.Title, componentID.Version, err.Error()), false, shared.NetworkingComponentName)
			return transport.ComponentID{}, err
		}
		return transport.ComponentID{}, err
	}

	if containsComponentID(c.registeredComponents, *componentID) {
		c.Logger.Info(fmt.Sprintf("Component %s is already registered.", componentID.Title), false, shared.NetworkingComponentName)
		return *componentID, nil
	}

	c.registeredComponents = append(c.registeredComponents, *componentID)

	c.Logger.Info(fmt.Sprintf("Router is running. Components Registered: %d", len(c.GetRegisteredComponents())), false, shared.NetworkingComponentName)

	return *componentID, nil
}

func (c *ComponentRegistrar) RemoveRegisteredComponent(id string) {
	for i, component := range c.registeredComponents {
		if component.ComponentUniqueID == id {
			c.Logger.Debug(fmt.Sprintf("Component %s is being unregistered.", component.Title), false, shared.NetworkingComponentName)
			c.registeredComponents = append(c.registeredComponents[:i], c.registeredComponents[i+1:]...)
			break
		}
	}

	c.Logger.Info(fmt.Sprintf("Router is running. Components Registered: %d", len(c.GetRegisteredComponents())), false, shared.NetworkingComponentName)
}

// GetRegisteredComponents returns a copy of the current registered components.
func (c *ComponentRegistrar) GetRegisteredComponents() []transport.ComponentID {
	num := len(c.registeredComponents)
	components := make([]transport.ComponentID, num)
	copy(components, c.registeredComponents)
	return components
}

func containsComponentID(s []transport.ComponentID, e transport.ComponentID) bool {
	for _, a := range s {
		if a.Equal(e) {
			return true
		}
	}
	return false
}

func (c *ComponentRegistrar) GetComponentByID(id string) (transport.ComponentID, error) {
	for _, component := range c.registeredComponents {
		if component.ComponentUniqueID == id {
			return component, nil
		}
	}
	return transport.ComponentID{}, fmt.Errorf("component with ID '%s' not found", id)
}
