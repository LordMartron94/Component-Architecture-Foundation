package lifecycle_managing

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/shared"
)

type ShutdownInterface interface {
	Shutdown() error
}

type LifeCycleManager struct {
	Logger            logging.HoornLogger
	WaitGroup         *sync.WaitGroup
	ShutdownListeners []ShutdownInterface
}

func (l *LifeCycleManager) ListenForTermination() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGQUIT)
	<-sigs
	l.Logger.Info("Received termination signal, shutting down server gracefully...", false, shared.MainComponentName)
	l.shutdownServer()
}

func (l *LifeCycleManager) shutdownServer() {
	for _, listener := range l.ShutdownListeners {
		err := listener.Shutdown()
		if err != nil {
			l.Logger.Error(fmt.Sprintf("Error shutting down listener: %v", err), false, shared.MainComponentName)
		}
	}

	l.WaitGroup.Done()
}
