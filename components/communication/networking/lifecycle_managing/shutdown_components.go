package lifecycle_managing

import (
	"fmt"
	"time"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/component_registration"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
)

type RequestInterface interface {
	SendRequest(id string, payload transport.MessagePayload) (string, error)
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

		_, err = s.Listener.SendRequest(component.ComponentUniqueID, payload)
		if err != nil {
			s.Logger.Warn(fmt.Sprintf("Unable to shut down listener because: '%s'", err), false, shared.MainComponentName)
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
