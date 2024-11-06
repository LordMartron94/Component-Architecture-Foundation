package networking

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/authentication"
	"github.com/component-architecture-foundation/networking/coding"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type Router struct {
	Logger               logging.HoornLogger
	MessageChannel       chan transport.Message
	Listener             NetworkHandlerInterface
	registeredComponents []transport.ComponentID
	WaitGroup            *sync.WaitGroup
}

func NewRouter(logger logging.HoornLogger, wg *sync.WaitGroup) *Router {
	msgChan := make(chan transport.Message)

	listener := NewTCPHandler(logger, "127.0.0.1:"+shared.ListeningPort, &authentication.WhitelistAuthenticator{Logger: logger}, msgChan, coding.JsonMessageCoder{Logger: logger})

	return &Router{
		Logger:         logger,
		Listener:       listener,
		MessageChannel: msgChan,
		WaitGroup:      wg,
	}
}

func containsComponentID(s []transport.ComponentID, e transport.ComponentID) bool {
	for _, a := range s {
		if a.Equal(e) {
			return true
		}
	}
	return false
}

func (r *Router) registerComponent(message transport.Message) {
	// Checks if the component is already registered, and if so, returns.
	if containsComponentID(r.registeredComponents, message.Requester) {
		r.Logger.Info(fmt.Sprintf("Component %s is already registered.", message.Requester.Title), false, shared.MainComponentName)
		return
	}

	r.registeredComponents = append(r.registeredComponents, message.Requester)
}

func (r *Router) getTargetComponentForMessage(message transport.Message) {

}

func (r *Router) StartListening() {
	r.WaitGroup.Add(1)
	go func() {
		err := r.Listener.StartListenLoop()
		if err != nil {
			r.Logger.Critical(fmt.Sprintf("Something went wrong while starting the network Listener: '%s'", err.Error()), false, shared.MainComponentName)
			return
		}
	}()

	go r.listenForTermination()
	go r.handleMessages()
}

func (r *Router) listenForTermination() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGQUIT)
	<-sigs
	r.Logger.Info("Received termination signal, shutting down server gracefully...", false, shared.MainComponentName)
	r.shutdownServer()
}

func (r *Router) shutdownServer() {
	err := r.Listener.Shutdown()
	if err != nil {
		r.Logger.Error(fmt.Sprintf("Something went wrong while shutting down the network Listener: '%s'", err.Error()), false, shared.MainComponentName)
	}

	r.WaitGroup.Done()
}

func (r *Router) handleMessages() {
	for {
		select {
		case msg := <-r.MessageChannel:
			r.Logger.Debug(fmt.Sprintf("Gotten message from '%s@%s' with payload '%s'", msg.Requester.Title, msg.Requester.Version, msg.Payload), false, shared.MainComponentName)
			r.processMessage(msg)
		default:
			time.Sleep(time.Millisecond * 10) // Or other small duration
		}
	}
}

func (r *Router) processMessage(message transport.Message) {
	// TODO: Implement message processing logic here
	if message.Payload.Action == "register" {
		r.Logger.Info(fmt.Sprintf("Received registration request from '%s@%s'", message.Requester.Title, message.Requester.Version), false, shared.MainComponentName)
		r.registerComponent(message)
	}
}
