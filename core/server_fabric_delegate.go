package core

import (
	"github.com/galenliu/chip/app/server"
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
)

type ServerFabricDelegateImpl struct {
	mServer *server.Server
}

func NewServerFabricDelegateImpl() *ServerFabricDelegateImpl {
	return &ServerFabricDelegateImpl{}
}

func (s ServerFabricDelegateImpl) Init(server *server.Server) error {
	s.mServer = server
	return nil
}

func (s ServerFabricDelegateImpl) FabricWillBeRemoved(table *credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s ServerFabricDelegateImpl) OnFabricRemoved(table *credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s ServerFabricDelegateImpl) OnFabricCommitted(table *credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s ServerFabricDelegateImpl) OnFabricUpdated(table *credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}
