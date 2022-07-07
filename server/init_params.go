package server

import (
	"github.com/galenliu/chip/access"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/lib"
	"net"
)

type InitParams struct {
	OperationalServicePort        uint16
	UserDirectedCommissioningPort uint16
	InterfaceId                   net.Interface
	AppDelegate                   any //unknown
	PersistentStorageDelegate     lib.PersistentStorageDelegate
	SessionResumptionStorage      any
	AccessDelegate                access.Delegate
	AclStorage                    AclStorage
	EndpointNativeParams          func()
}

func (i InitParams) Default() *InitParams {
	i.OperationalServicePort = config.ChipPort
	i.UserDirectedCommissioningPort = config.ChipUdcPort
	return &i
}

func InitializeStaticResourcesBeforeServerInit() (initParams *InitParams) {
	return InitParams{}.Default()
}
