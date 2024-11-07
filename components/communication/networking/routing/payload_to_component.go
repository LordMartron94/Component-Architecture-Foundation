package routing

import (
	"fmt"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type PayloadToComponent struct {
	Logger *logging.HoornLogger
}

// SearchForComponent finds the associated component for a message payload. Returns an error if no match is found.
func (p *PayloadToComponent) SearchForComponent(payload transport.MessagePayload, components []transport.ComponentID) (transport.ComponentID, error) {
	for _, component := range components {
		p.Logger.Debug(fmt.Sprintf("Searching for component with ID '%s'", component.Title), false, shared.RoutingComponentName)
		for _, capability := range component.Capabilities {
			p.Logger.Debug(fmt.Sprintf("Checking capability '%s' for component '%s'", capability.Name, component.Title), false, shared.RoutingComponentName)
			if capability.Name != payload.Action {
				continue
			}

			if p.signatureMatchesArgs(capability.Signature, payload.Args) {
				p.Logger.Debug(fmt.Sprintf("Found component with ID '%s' for payload", component.Title), false, shared.RoutingComponentName)
				return component, nil
			}
		}
	}

	p.Logger.Warn(fmt.Sprintf("No matching component found for payload with action request '%s'", payload.Action), false, shared.RoutingComponentName)
	return transport.ComponentID{}, fmt.Errorf("no match found for payload")
}

func (p *PayloadToComponent) signatureMatchesArgs(signature transport.Signature, args []transport.Argument) bool {
	if signature.NumOfArgs != len(args) {
		return false
	}

	for i, arg := range args {
		if signature.TypeOfArgs[i] != arg.Type {
			return false
		}
	}

	return true
}
