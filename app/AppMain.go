package app

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/device_layer"
	"github.com/galenliu/chip/platform"
	"github.com/galenliu/chip/server"
	"github.com/galenliu/chip/storage"
	log "github.com/sirupsen/logrus"

	"time"
)

func AppMainInit() error {

	err := storage.KeyValueStoreMgr().Init(config.GetDeviceOptionsInstance().KVS)
	if err != nil {
		log.Panic(err.Error())
	}

	for i := 0; i < 100; i++ {
		err = storage.KeyValueStoreMgr().WriteValueStr(fmt.Sprintf("%d", i), time.Now().String())
		if err != nil {
			log.Infof(err.Error())
		}
	}

	mgr, err := platform.ConfigurationMgr().Init(platform.GetConfigProviderInstance())
	_ = platform.NewDeviceInstanceInfo(mgr)

	_, err = DeviceLayer.NewCommissionableData(config.GetDeviceOptionsInstance())
	if err != nil {
		return err
	}
	return nil
}

func AppMainLoop() error {

	initParams := server.NewCommonCaseDeviceServerInitParams()

	err := initParams.InitializeStaticResourcesBeforeServerInit()
	if err != nil {
		return err
	}
	chip := server.NewServer(initParams)
	chip.Shutdown()
	return nil
}
