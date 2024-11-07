package lifecycle_managing

import (
	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/component_registration"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type ResponseInterface interface {
	SendResponse(component transport.ComponentID, response []byte) error
}

type ShutdownComponents struct {
	ComponentRegistrar *component_registration.ComponentRegistrar
	Logger             logging.HoornLogger
	Listener           ResponseInterface
}

func (s ShutdownComponents) Shutdown() error {
	for _, component := range s.ComponentRegistrar.GetRegisteredComponents() {
		err := s.Listener.SendResponse(component, []byte(shared.ShutdownRequestPayload))
		if err != nil {
			s.Logger.Warn("Unable to shut down listener", false, shared.MainComponentName)
		}

		s.ComponentRegistrar.RemoveRegisteredComponent(component)
	}

	return nil
}
