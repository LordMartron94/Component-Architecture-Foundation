package component_registration

import (
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type ComponentRegistrar struct {
	Logger               logging.HoornLogger
	registeredComponents []transport.ComponentID
}

func (c *ComponentRegistrar) RegisterComponent(message transport.Message) {
	if containsComponentID(c.registeredComponents, message.Requester) {
		c.Logger.Info(fmt.Sprintf("Component %s is already registered.", message.Requester.Title), false, shared.MainComponentName)
		return
	}

	c.registeredComponents = append(c.registeredComponents, message.Requester)
}

func (c *ComponentRegistrar) RemoveRegisteredComponent(id transport.ComponentID) {
	for i, component := range c.registeredComponents {
		if component.Equal(id) {
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
