package app

import (
	"github.com/galenliu/chip/platform"
	"github.com/galenliu/chip/platform/device_layer"
	"github.com/galenliu/chip/platform/options"
	"github.com/galenliu/chip/platform/storage"
	"github.com/galenliu/chip/server"
	"github.com/galenliu/gateway/pkg/log"
)

func AppMainInit(info *DeviceLayer.DeviceInfo, options *options.DeviceOptions) error {

	err := storage.KeyValueStoreMgr().Init(options.KVS)
	if err != nil {
		log.Infof("store init err: %s", err.Error())
		return err
	}
	mgr := platform.NewConfigurationManager(options)
	_ = DeviceLayer.NewDeviceInstanceInfo(mgr, info)

	_, err = DeviceLayer.NewCommissionableData(options)
	if err != nil {
		return err
	}
	return nil
}

func AppMainLoop(options *options.DeviceOptions, info *DeviceLayer.DeviceInfo) error {

	initParams := server.NewCommonCaseDeviceServerInitParams(options)

	err := initParams.InitializeStaticResourcesBeforeServerInit()
	if err != nil {
		return err
	}
	chip := server.NewServer(initParams)
	chip.Shutdown()
	return nil
}
