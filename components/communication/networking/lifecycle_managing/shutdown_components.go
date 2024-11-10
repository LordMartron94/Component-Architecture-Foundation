package lifecycle_managing

import (
	"fmt"
	"time"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/component_registration"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type RequestInterface interface {
	SendRequest(id string, payload transport.MessagePayload) error
}

type ShutdownComponents struct {
	ComponentRegistrar *component_registration.ComponentRegistrar
	Logger             *logging.HoornLogger
	Listener           RequestInterface
}

func (s ShutdownComponents) Shutdown() error {
	for _, component := range s.ComponentRegistrar.GetRegisteredComponents() {
		payload, err := transport.MessagePayloadFromBytes([]byte(shared.ShutdownRequestPayload))

		if err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to create shutdown request payload: '%s'", err.Error()), false, shared.NetworkingComponentName)
			continue
		}

		err = s.Listener.SendRequest(component.ComponentUniqueID, payload)
		if err != nil {
			s.Logger.Warn("Unable to shut down listener", false, shared.MainComponentName)
		}
	}

	timePassedInSeconds := 0

	for len(s.ComponentRegistrar.GetRegisteredComponents()) > 0 {
		if timePassedInSeconds >= 30 {
			s.Logger.Error("Unable to shut down all components within 30 seconds", false, shared.MainComponentName)
			return fmt.Errorf("unable to shut down all components within 30 seconds")
		}

		s.Logger.Info("Waiting for all components to be shut down", false, shared.MainComponentName)
		time.Sleep(time.Second * 1)
		timePassedInSeconds++
	}

	return nil
}
