package component_registration

import (
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/message_handling"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type SpecialActionPerformerInterface interface {
	PerformAnySpecialActions(component transport.ComponentID)
}

type ComponentRegistrar struct {
	Logger                 *logging.HoornLogger
	registeredComponents   []transport.ComponentID
	SpecialActionPerformer SpecialActionPerformerInterface
	Listener               message_handling.ListenerInterface
}

func (c *ComponentRegistrar) RegisterComponent(message transport.Message) transport.ComponentID {
	componentID, err := message.GetComponentIDFromRegistrationMessage(c.Logger)

	if err != nil {
		c.Logger.Error(fmt.Sprintf("Error getting component ID from message: %s (on registration)", err.Error()), false, shared.NetworkingComponentName)
		err := c.Listener.SendResponse(shared.ServerUUID, []byte(shared.InvalidRequestResponsePayload), componentID.ComponentUniqueID)
		if err != nil {
			c.Logger.Error(fmt.Sprintf("Failed to send response to '%s@%s': %s", componentID.Title, componentID.Version, err.Error()), false, shared.MainComponentName)
			return transport.ComponentID{}
		}
		return transport.ComponentID{}
	}

	if containsComponentID(c.registeredComponents, *componentID) {
		c.Logger.Info(fmt.Sprintf("Component %s is already registered.", componentID.Title), false, shared.MainComponentName)
		return *componentID
	}

	c.SpecialActionPerformer.PerformAnySpecialActions(*componentID)
	c.registeredComponents = append(c.registeredComponents, *componentID)

	return *componentID
}

func (c *ComponentRegistrar) RemoveRegisteredComponent(id string) {
	for i, component := range c.registeredComponents {
		if component.ComponentUniqueID == id {
			c.Logger.Debug(fmt.Sprintf("Component %s is being unregistered.", component.Title), false, shared.MainComponentName)
			c.registeredComponents = append(c.registeredComponents[:i], c.registeredComponents[i+1:]...)
			break
		}
	}
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
