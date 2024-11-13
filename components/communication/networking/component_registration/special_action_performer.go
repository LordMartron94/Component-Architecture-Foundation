package component_registration

import (
	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
)

type RequesterInterface interface {
	SendResponse(id string, response []byte, targetUUID string) (string, error)
}

type SpecialActionPerformer struct {
	Logger             *logging.HoornLogger
	RequesterInterface RequesterInterface
}

func (s *SpecialActionPerformer) PerformAnySpecialActions(component transport.ComponentID) {
	// Does nothing for now because there are no special components
}
