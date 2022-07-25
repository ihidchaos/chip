package server

type ServerFabricDelegate interface {
	Init(s *Server) error
}

type ServerFabricDelegateImpl struct {
}
