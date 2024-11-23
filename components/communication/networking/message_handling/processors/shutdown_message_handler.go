package processors

import (
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/lifecycle_managing"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
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
