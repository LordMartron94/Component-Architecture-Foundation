package networking

import (
	"fmt"
	"sync"
	"time"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/coding"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/component_registration"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/connectivity"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/connectivity/communication"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/connectivity/peer"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/handlers"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/lifecycle_managing"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/message_handling"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/message_handling/processors"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/routing"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
)

type Router struct {
	Logger         *logging.HoornLogger
	MessageChannel chan transport.Message
	Listener       NetworkHandlerInterface

	waitGroup          *sync.WaitGroup
	lifecycleManager   *lifecycle_managing.LifeCycleManager
	componentRegistrar *component_registration.ComponentRegistrar
	payloadToComponent *routing.PayloadToComponent
	messageHandlers    map[string]processors.MessageProcessorInterface
}

func NewRouter(logger *logging.HoornLogger, wg *sync.WaitGroup, messageCoder coding.MessageCoderInterface, shutdownChan chan struct{}) *Router {
	msgChan := make(chan transport.Message, 15)

	server := transport.ComponentID{
		Title:             shared.ServerName,
		Version:           shared.ServerVersion,
		Capabilities:      nil,
		ComponentUniqueID: shared.ServerUUID,
	}

	peerHandler := peer.NewPeerHandler(logger)
	messageUtility := message_handling.MessageUtility{Logger: logger, MessageCoder: messageCoder, Server: server}
	communicationHandler := communication.CommunicationHandler{
		Logger:         logger,
		PeerHandler:    peerHandler,
		MessageUtility: &messageUtility,
		MessageCoder:   messageCoder,
	}

	connectionHandler := connectivity.NewDefaultConnectionHandler(logger, shared.ListeningAddress, messageCoder, peerHandler, &messageUtility, &communicationHandler, msgChan, shutdownChan)
	listener := handlers.NewTCPHandler(logger, connectionHandler, &communicationHandler, shutdownChan, wg)

	componentRegistrar := component_registration.ComponentRegistrar{
		Logger:   logger,
		Listener: listener,
	}

	shutdownComponents := lifecycle_managing.ShutdownComponents{
		ComponentRegistrar: &componentRegistrar,
		Logger:             logger,
		Listener:           listener,
	}

	shutdownListeners := make([]lifecycle_managing.ShutdownInterface, 0)
	shutdownListeners = append(shutdownListeners, &shutdownComponents, listener)

	lifecycleManager := lifecycle_managing.LifeCycleManager{
		Logger:            logger,
		ShutdownListeners: shutdownListeners,
		WaitGroup:         wg,
	}

	router := &Router{
		Logger:             logger,
		Listener:           listener,
		MessageChannel:     msgChan,
		waitGroup:          wg,
		lifecycleManager:   &lifecycleManager,
		payloadToComponent: &routing.PayloadToComponent{Logger: logger},
		componentRegistrar: &componentRegistrar,
	}
	router.messageHandlers = make(map[string]processors.MessageProcessorInterface)

	router.RegisterMessageHandler("register", &processors.RegisterMessageHandler{
		Logger:            logger,
		Listener:          listener,
		RegisterComponent: componentRegistrar.RegisterComponent,
	})
	router.RegisterMessageHandler("shutdown", &processors.ShutdownMessageHandler{
		Logger:           logger,
		LifeCycleManager: lifecycleManager,
	})
	router.RegisterMessageHandler("__default__", &message_handling.DefaultMessageHandler{
		Logger:                       logger,
		Listener:                     listener,
		GetTargetComponentForMessage: router.getTargetComponentForMessage,
		GetComponentByID:             componentRegistrar.GetComponentByID,
		GetExpectedClientResponses:   router.getExpectedClientResponses,
	})
	router.RegisterMessageHandler("unregister", &processors.UnregisterMessageHandler{
		Logger:             logger,
		ComponentRegistrar: &componentRegistrar,
		GetComponentByID:   componentRegistrar.GetComponentByID,
	})
	router.RegisterMessageHandler("keep_alive", &processors.KeepAliveMessageHandler{
		Logger:      logger,
		PeerHandler: peerHandler,
	})

	peerHandler.SetStopConnection(router.StopConnection)

	return router
}

func (r *Router) RegisterMessageHandler(action string, handler processors.MessageProcessorInterface) {
	r.messageHandlers[action] = handler
}

func (r *Router) getTargetComponentForMessage(message transport.Message) (transport.ComponentID, error) {
	return r.payloadToComponent.SearchForComponent(*message.Payload, r.componentRegistrar.GetRegisteredComponents())
}

func (r *Router) Start() {
	go func() {
		err := r.Listener.StartListenLoop()
		if err != nil {
			r.Logger.Critical(fmt.Sprintf("Something went wrong while starting the network Listener: '%s'", err.Error()), false, shared.MainComponentName)
			return
		}
	}()

	go r.lifecycleManager.ListenForTermination()
	go r.handleMessages()
	go r.info()
}

func (r *Router) handleMessages() {
	for {
		select {
		case msg := <-r.MessageChannel:
			componentID, err := r.componentRegistrar.GetComponentByID(*msg.SenderID)

			if err != nil {
				r.Logger.Warn(fmt.Sprintf("Failed to get component ID from message: '%s' (this can be ignored pre-registration)", err.Error()), false, shared.MainComponentName)
			}

			go r.ProcessMessage(msg, componentID, err)
		}
	}
}

func (r *Router) ProcessMessage(message transport.Message, componentID transport.ComponentID, err1 error) {
	for action, processor := range r.messageHandlers {
		if action == message.Payload.Action {
			err := processor.ProcessMessage(message)
			if err != nil {
				if err1 != nil {
					r.Logger.Error(fmt.Sprintf("Error processing message '%s' for component '%s': '%s'", message.Payload.Action, *message.SenderID, err.Error()), false, shared.MainComponentName)
				} else {
					r.Logger.Error(fmt.Sprintf("Error processing message '%s' for component '%s': '%s'", message.Payload.Action, componentID.Title, err.Error()), false, shared.MainComponentName)
				}
				return
			}
			return
		}
	}

	//r.Logger.Debug(fmt.Sprintf("Received message '%s' with normal action. Resorting to default processor", componentID.Title), false, shared.MainComponentName)
	defaultProcessor := r.messageHandlers["__default__"]
	err := defaultProcessor.ProcessMessage(message)
	if err != nil {
		if err1 != nil {
			r.Logger.Error(fmt.Sprintf("Error processing message '%s' for component '%s': '%s'", message.Payload.Action, *message.SenderID, err.Error()), false, shared.MainComponentName)
		} else {
			r.Logger.Error(fmt.Sprintf("Error processing message '%s' for component '%s': '%s'", message.Payload.Action, componentID.Title, err.Error()), false, shared.MainComponentName)
		}
	}
}

func (r *Router) info() {
	for {
		r.Logger.Debug(fmt.Sprintf("Router is running. Components registered: %d | Active peers: %d", len(r.componentRegistrar.GetRegisteredComponents()), r.Listener.GetActiveConnectionsNumber()), false, shared.InfoComponentName)
		time.Sleep(time.Minute * 10)
	}
}

func (r *Router) getExpectedClientResponses(message transport.Message) int {
	return r.payloadToComponent.GetExpectedClientResponses(*message.Payload, r.componentRegistrar.GetRegisteredComponents())
}

func (r *Router) StopConnection(p *peer.Peer) error {
	return r.Listener.StopConnection(p)
}

// IsComponentRegistered is used for Benchmarks only since it uses a componentUUID hardcoded string.
func (r *Router) IsComponentRegistered(componentUUID string) bool {
	registered := r.componentRegistrar.GetRegisteredComponents()

	for _, component := range registered {
		if component.ComponentUniqueID == componentUUID {
			return true
		}
	}

	return false
}

func (r *Router) Stop() {
	r.lifecycleManager.ShutdownServer()
}

func (r *Router) WaitUntilShutdown() {
	r.waitGroup.Wait()
}
