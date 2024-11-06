package handlers

type ConnectionHandlerInterface interface {
	StartListenLoop() error
}
