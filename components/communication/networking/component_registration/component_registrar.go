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

func (c *ComponentRegistrar) GetRegisteredComponents() []transport.ComponentID {
	return c.registeredComponents
}

func containsComponentID(s []transport.ComponentID, e transport.ComponentID) bool {
	for _, a := range s {
		if a.Equal(e) {
			return true
		}
	}
	return false
}
