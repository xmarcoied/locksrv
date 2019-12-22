package server

// Server is the main contract for the locking server
type Server interface {
	ListenAndServer() error
	Shutdown()
}

type none struct{}
