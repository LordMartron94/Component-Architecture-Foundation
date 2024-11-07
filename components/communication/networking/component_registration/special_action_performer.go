package component_registration

import (
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type RequesterInterface interface {
	SendResponse(id transport.ComponentID, response []byte) error
}

type SpecialActionPerformer struct {
	Logger              *logging.HoornLogger
	RequesterInterface  RequesterInterface
	SetupLoggingPayload []byte
}

func (s *SpecialActionPerformer) PerformAnySpecialActions(component transport.ComponentID) {
	if component.HasCapabilityAction("setup_logging") {
		err := s.RequesterInterface.SendResponse(component, s.SetupLoggingPayload)
		if err != nil {
			s.Logger.Error(fmt.Sprintf("Failed to send setup logging payload to '%s@%s': %s", component.Title, component.Version, err.Error()), false, shared.RoutingComponentName)
			return
		}
	}
}
