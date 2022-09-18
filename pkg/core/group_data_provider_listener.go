package core

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/gateway/pkg/log"
)

type GroupDataProviderListener interface {
	credentials.GroupListener
	Init(server *Server) error
}

type GroupDataProviderListenerImpl struct {
	mServer *Server
}

func (g *GroupDataProviderListenerImpl) OnGroupAdded(fabricIndex lib.FabricIndex, newGroup *credentials.GroupInfo) {
	fabric := g.mServer.GetFabricTable().FindFabricWithIndex(fabricIndex)
	if fabric == nil {
		log.Info("Group added to nonexistent fabric?")
		return
	}
	if err := g.mServer.GetTransportManager().MulticastGroupJoinLeave(lib.Multicast(fabric.GetFabricId(), newGroup.Id), true); err != nil {

	}
}

func (g *GroupDataProviderListenerImpl) OnGroupRemoved(fabricIndex lib.FabricIndex, newGroup *credentials.GroupInfo) {
	fabric := g.mServer.GetFabricTable().FindFabricWithIndex(fabricIndex)
	if fabric == nil {
		log.Info("Group added to nonexistent fabric?")
		return
	}
	_ = g.mServer.GetTransportManager().MulticastGroupJoinLeave(lib.Multicast(fabric.GetFabricId(), newGroup.Id), false)
}

func (g *GroupDataProviderListenerImpl) Init(server *Server) error {
	g.mServer = server
	return nil
}
