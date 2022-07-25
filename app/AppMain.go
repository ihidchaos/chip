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

	mgr, err := platform.ConfigurationMgr().Init(platform.GetConfigProviderInstance())
	_ = platform.NewDeviceInstanceInfo(mgr)

	_, err = DeviceLayer.NewCommissionableData(options)
	if err != nil {
		return err
	}
	return nil
}

func AppMainLoop() error {

	serverInitParams := server.NewCommonCaseDeviceServerInitParams()

	err := serverInitParams.InitializeStaticResourcesBeforeServerInit()
	if err != nil {
		return err
	}
	chip, _ := server.NewCHIPServer(&serverInitParams.ServerInitParams)
	RunLoop()
	chip.Shutdown()
	return nil
}

func RunLoop() {
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
