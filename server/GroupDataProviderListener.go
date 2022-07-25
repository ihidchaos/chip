package server

type GroupDataProviderListener interface {
	Init(s *Server) error
}

type GroupDataProviderListenerImpl struct {
}
