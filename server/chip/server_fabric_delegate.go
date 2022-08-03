package chip

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
)

type ServerFabricDelegateImpl struct {
	mServer *Server
}

func NewServerFabricDelegateImpl() *ServerFabricDelegateImpl {
	return &ServerFabricDelegateImpl{}
}

func (s ServerFabricDelegateImpl) Init(server *Server) error {
	s.mServer = server
	return nil
}

func (s ServerFabricDelegateImpl) FabricWillBeRemoved(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s ServerFabricDelegateImpl) OnFabricRemoved(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s ServerFabricDelegateImpl) OnFabricCommitted(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s ServerFabricDelegateImpl) OnFabricUpdated(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}
