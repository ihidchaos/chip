package core

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	log "golang.org/x/exp/slog"
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
	if err := g.mServer.mTransports.MulticastGroupJoinLeave(lib.Multicast(fabric.FabricId(), newGroup.Id), true); err != nil {

	}
}

func (g *GroupDataProviderListenerImpl) OnGroupRemoved(fabricIndex lib.FabricIndex, newGroup *credentials.GroupInfo) {
	fabric := g.mServer.GetFabricTable().FindFabricWithIndex(fabricIndex)
	if fabric == nil {
		log.Info("Group added to nonexistent fabric?")
		return
	}
	_ = g.mServer.mTransports.MulticastGroupJoinLeave(lib.Multicast(fabric.FabricId(), newGroup.Id), false)
}

func (g *GroupDataProviderListenerImpl) Init(server *Server) error {
	g.mServer = server
	return nil
}
