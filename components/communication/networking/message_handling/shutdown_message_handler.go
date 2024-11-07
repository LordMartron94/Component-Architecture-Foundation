package message_handling

import (
	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/lifecycle_managing"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type ShutdownMessageHandler struct {
	Logger           *logging.HoornLogger
	LifeCycleManager lifecycle_managing.LifeCycleManager
}

func (s ShutdownMessageHandler) ProcessMessage(message transport.Message) error {
	s.Logger.Info("Received Shutdown Request, will shutdown.", false, shared.RoutingComponentName)
	s.LifeCycleManager.ShutdownServer()
	return nil
}
