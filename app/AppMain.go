package app

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials/dac"
	"github.com/galenliu/chip/device"
	core2 "github.com/galenliu/chip/pkg/core"
	"github.com/galenliu/chip/pkg/storage"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"time"
)

func Init(options *config.DeviceOptions) error {

	var rendezvousFlags uint8
	if config.NetworkLayerBle {
		rendezvousFlags = config.RendezvousInformationFlagBLE
	} else {
		rendezvousFlags = config.RendezvousInformationFlagOnNetwork
	}

	if config.RendezvousMode {
		log.Infof("RendezvousMode")
	}

	err := storage.KeyValueStoreMgr().Init(options.KVS)
	if err != nil {
		log.Panic(err.Error())
	}

	err = storage.KeyValueStoreMgr().WriteValueStr("Reboot", time.Now().String())
	if err != nil {
		log.Infof(err.Error())
	}

	commissionableDataProvider := device.NewCommissionableDataImpl()
	commissionableDataProvider, err = commissionableDataProvider.Init(options)
	if err != nil {
		log.Infof(err.Error())
		return err
	}

	configProvider := config.NewConfigProviderImpl()
	configMgr := config.NewConfigurationManagerImpl()
	configMgr, err = configMgr.Init(configProvider, options)
	if err != nil {
		log.Infof(err.Error())
		return err
	}

	if options.Payload.RendezvousInformation != 0 {
		rendezvousFlags = options.Payload.RendezvousInformation
	}
	err = core2.GetPayloadContents(&options.Payload, rendezvousFlags)
	if err != nil {
		return err
	}

	{
		options.Payload.CommissioningFlow = config.CommissioningFlowCustom
		core2.PrintOnboardingCodes(options.Payload)

	}

	dac.SetDeviceAttestationCredentialsProvider(options.DacProvider)

	deviceInstanceInfo := device.NewDeviceInstanceInfo()
	deviceInstanceInfo, err = deviceInstanceInfo.Init(configMgr)
	if err != nil {
		log.Infof(err.Error())
		return err
	}

	return nil
}

func MainLoop(options *config.DeviceOptions) error {

	serverInitParams := core2.NewServerInitParams()
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
	chipServer := core2.NewCHIPServer()
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
