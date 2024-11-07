package networking

import (
	"fmt"
	"sync"
	"time"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/networking/coding"
	"github.com/component-architecture-foundation/networking/component_registration"
	"github.com/component-architecture-foundation/networking/connectivity"
	"github.com/component-architecture-foundation/networking/connectivity/communication"
	"github.com/component-architecture-foundation/networking/connectivity/peer"
	"github.com/component-architecture-foundation/networking/handlers"
	"github.com/component-architecture-foundation/networking/lifecycle_managing"
	"github.com/component-architecture-foundation/networking/message_handling"
	"github.com/component-architecture-foundation/networking/routing"
	"github.com/component-architecture-foundation/networking/transport"
	"github.com/component-architecture-foundation/shared"
)

type Router struct {
	Logger         *logging.HoornLogger
	MessageChannel chan transport.Message
	Listener       NetworkHandlerInterface

	waitGroup           *sync.WaitGroup
	lifecycleManager    *lifecycle_managing.LifeCycleManager
	componentRegistrar  *component_registration.ComponentRegistrar
	payloadToComponent  *routing.PayloadToComponent
	messageHandlers     map[string]message_handling.MessageProcessorInterface
	setupLoggingPayload []byte
}

func NewRouter(logger *logging.HoornLogger, wg *sync.WaitGroup, messageCoder coding.MessageCoderInterface, setupLoggingPayload []byte) *Router {
	msgChan := make(chan transport.Message)
	shutdownChan := make(chan struct{})

	server := transport.ComponentID{
		Title:        shared.ServerName,
		Version:      shared.ServerVersion,
		Capabilities: nil,
	}

	peerHandler := peer.PeerHandler{Logger: logger}
	messageUtility := message_handling.MessageUtility{Logger: logger, MessageCoder: messageCoder, Server: server}
	communicationHandler := communication.CommunicationHandler{
		Logger:         logger,
		PeerHandler:    &peerHandler,
		MessageUtility: &messageUtility,
		MessageCoder:   messageCoder,
	}

	connectionHandler := connectivity.NewDefaultConnectionHandler(logger, shared.ListeningAddress, messageCoder, &peerHandler, &messageUtility, &communicationHandler, msgChan, shutdownChan)
	listener := handlers.NewTCPHandler(logger, msgChan, connectionHandler, &communicationHandler, shutdownChan)

	specialActionPerformer := component_registration.SpecialActionPerformer{
		Logger:              logger,
		RequesterInterface:  listener,
		SetupLoggingPayload: setupLoggingPayload,
	}
	componentRegistrar := component_registration.ComponentRegistrar{Logger: logger, SpecialActionPerformer: &specialActionPerformer}

	shutdownComponents := lifecycle_managing.ShutdownComponents{
		ComponentRegistrar: &componentRegistrar,
		Logger:             logging.HoornLogger{},
		Listener:           listener,
	}

	shutdownListeners := make([]lifecycle_managing.ShutdownInterface, 0)
	shutdownListeners = append(shutdownListeners, listener)
	shutdownListeners = append(shutdownListeners, &shutdownComponents)

	lifecycleManager := lifecycle_managing.LifeCycleManager{
		Logger:            logger,
		ShutdownListeners: shutdownListeners,
		WaitGroup:         wg,
	}

	router := &Router{
		Logger:              logger,
		Listener:            listener,
		MessageChannel:      msgChan,
		waitGroup:           wg,
		lifecycleManager:    &lifecycleManager,
		payloadToComponent:  &routing.PayloadToComponent{Logger: logger},
		componentRegistrar:  &componentRegistrar,
		setupLoggingPayload: setupLoggingPayload,
	}
	router.messageHandlers = make(map[string]message_handling.MessageProcessorInterface)

	router.RegisterMessageHandler("register", &message_handling.RegisterMessageHandler{
		Logger:            logger,
		Listener:          listener,
		RegisterComponent: componentRegistrar.RegisterComponent,
	})
	router.RegisterMessageHandler("__default__", &message_handling.DefaultMessageHandler{
		Logger:                       logger,
		Listener:                     listener,
		GetTargetComponentForMessage: router.getTargetComponentForMessage,
	})

	return router
}

func (r *Router) RegisterMessageHandler(action string, handler message_handling.MessageProcessorInterface) {
	r.messageHandlers[action] = handler
}

func (r *Router) getTargetComponentForMessage(message transport.Message) (transport.ComponentID, error) {
	return r.payloadToComponent.SearchForComponent(message.Payload, r.componentRegistrar.GetRegisteredComponents())
}

func (r *Router) Start() {
	r.waitGroup.Add(1)
	go func() {
		err := r.Listener.StartListenLoop()
		if err != nil {
			r.Logger.Critical(fmt.Sprintf("Something went wrong while starting the network Listener: '%s'", err.Error()), false, shared.MainComponentName)
			return
		}
	}()

	go r.lifecycleManager.ListenForTermination()
	go r.handleMessages()
}

func (r *Router) handleMessages() {
	for {
		select {
		case msg := <-r.MessageChannel:
			r.Logger.Debug(fmt.Sprintf("Gotten message from '%s@%s' with payload '%s'", msg.Requester.Title, msg.Requester.Version, msg.Payload), false, shared.MainComponentName)
			r.processMessage(msg)
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func (r *Router) processMessage(message transport.Message) {
	for action, processor := range r.messageHandlers {
		if action == message.Payload.Action {
			err := processor.ProcessMessage(message)
			if err != nil {
				r.Logger.Error(fmt.Sprintf("Error processing message '%s' for component '%s': '%s'", message.Payload.Action, message.Requester.Title, err.Error()), false, shared.MainComponentName)
				return
			}
			return
		}
	}

	r.Logger.Debug(fmt.Sprintf("Received message '%s' with normal action. Resorting to default processor", message.Requester.Title), false, shared.MainComponentName)
	defaultProcessor := r.messageHandlers["__default__"]
	err := defaultProcessor.ProcessMessage(message)
	if err != nil {
		r.Logger.Error(fmt.Sprintf("Error processing message '%s' for component '%s': '%s'", message.Payload.Action, message.Requester.Title, err.Error()), false, shared.MainComponentName)
	}
}
