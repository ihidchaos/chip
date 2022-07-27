package app

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/device_layer"
	"github.com/galenliu/chip/platform"
	"github.com/galenliu/chip/server"
	"github.com/galenliu/chip/storage"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"time"
)

func AppMainInit(options *config.DeviceOptions) error {

	err := storage.KeyValueStoreMgr().Init(options.KVS)
	if err != nil {
		log.Panic(err.Error())
	}

	err = storage.KeyValueStoreMgr().WriteValueStr(fmt.Sprintf("%d", time.Now().Minute()), time.Now().String())

	if err != nil {
		log.Infof(err.Error())
	}

	configProvider := config.NewConfigProviderImpl()
	configMgr := config.NewConfigurationManagerImpl()
	configMgr, err = configMgr.Init(configProvider, options)
	if err != nil {
		log.Infof(err.Error())
		return err
	}

	deviceInstanceInfo := platform.NewDeviceInstanceInfo()
	deviceInstanceInfo, err = deviceInstanceInfo.Init(configMgr)
	if err != nil {
		log.Infof(err.Error())
		return err
	}

	commissionableDataProvider := DeviceLayer.NewCommissionableDataImpl()
	commissionableDataProvider, err = commissionableDataProvider.Init(options)
	if err != nil {
		log.Infof(err.Error())
		return err
	}
	return nil
}

func AppMainLoop(options *config.DeviceOptions) error {

	serverInitParams := server.NewServerInitParams()
	_, err := serverInitParams.Init(options)
	if err != nil {
		log.Infof(err.Error())
		return err
	}

	err = serverInitParams.InitializeStaticResourcesBeforeServerInit()
	if err != nil {
		log.Infof(err.Error())
		return err
	}
	chipServer := server.NewCHIPServer()
	chipServer, err = chipServer.Init(serverInitParams)
	if err != nil {
		return err
	}
	WaitSignal()
	chipServer.Shutdown()
	return nil
}

func WaitSignal() {
	sigs := make(chan os.Signal, 1)

	done := make(chan bool, 1)

	// `signal.Notify` registers the given channel to

	// receive notifications of the specified signals.

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// This goroutine executes a blocking receive for

	// signals. When it gets one it'll print it out

	// and then notify the program that it can finish.

	go func() {

		sig := <-sigs

		fmt.Println()

		fmt.Println(sig)

		done <- true

	}()

	// The program will wait here until it gets the

	// expected signal (as indicated by the goroutine

	// above sending a value on `done`) and then exit.

	fmt.Println("awaiting signal")

	<-done

	fmt.Println("exiting")
}
