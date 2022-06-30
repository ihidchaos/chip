package server

import (
	"github.com/galenliu/chip/access"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/inet/Interface"
	"github.com/galenliu/chip/lib"
	"net"
)

type InitParams struct {
	OperationalServicePort        int
	UserDirectedCommissioningPort int
	InterfaceId                   Interface.Id
	AppDelegate                   any //unknown
	PersistentStorageDelegate     lib.PersistentStorageDelegate
	SessionResumptionStorage      any
	AccessDelegate                access.Delegate
	AclStorage                    AclStorage
	EndpointNativeParams          func()
}

func DefaultServerInitParams() *InitParams {
	return &InitParams{
		OperationalServicePort:        config.ChipPort,
		UserDirectedCommissioningPort: config.ChipUdcPort,
	}
}

func InitializeStaticResourcesBeforeServerInit() (initParams InitParams) {
	initParams = InitParams{
		OperationalServicePort:        0,
		UserDirectedCommissioningPort: 0,
	}
	list, _ := net.Interfaces()
	for _, inter := range list {
		adders, _ := inter.Addrs()
		if len(adders) > 1 {
			initParams.InterfaceId = Interface.Id{Interface: inter}
		}
	}
	return
}
