package connectivity

type ConnectionHandlerInterface interface {
	StartListenLoop() error
	CloseConnections()
}
