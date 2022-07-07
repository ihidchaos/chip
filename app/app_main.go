package app

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/server"
)

type Config struct {
	ConfigNetworkLayerBle bool
	mSecuredServicePort   int
	mUnsecuredServicePort int
}

func AppMainInit(con Config) {

}

func AppMainLoop() {
	initParams := server.InitializeStaticResourcesBeforeServerInit()
	initParams.OperationalServicePort = config.ChipPort
	initParams.UserDirectedCommissioningPort = config.ChipUdcPort
	chip := new(server.Server).Init(initParams)
	chip.Shutdown()
}
