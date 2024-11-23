package lifecycle_managing

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
)

type ShutdownInterface interface {
	Shutdown() error
}

type LifeCycleManager struct {
	Logger            *logging.HoornLogger
	WaitGroup         *sync.WaitGroup
	ShutdownListeners []ShutdownInterface
}

func (l *LifeCycleManager) ListenForTermination() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGQUIT)

	sig := <-sigs
	l.Logger.Info(fmt.Sprintf("Received termination signal: %v, shutting down server gracefully...", sig), false, shared.MainComponentName)
	l.ShutdownServer()
	l.Logger.Info("Server shutdown complete.", false, shared.MainComponentName)
}

func (l *LifeCycleManager) ShutdownServer() {
	for _, listener := range l.ShutdownListeners {
		err := listener.Shutdown()
		if err != nil {
			l.Logger.Error(fmt.Sprintf("Error shutting down listener: %v", err), false, shared.MainComponentName)
		}
	}

	l.WaitGroup.Done()
}
