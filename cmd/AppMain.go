package main

import (
	"fmt"
	"github.com/galenliu/chip/app/server"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/core"
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/device"
	"github.com/galenliu/chip/lib/store"
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

	store.DefaultPersistentStorage().SetStorage(store.NewInitStorage(options.KVS))

	err := store.DefaultPersistentStorage().SetKeyValue("Reboot", time.Now().String())
	if err != nil {
		return err
	}

	commissionableDataProvider := device.DefaultCommissionableDateProvider()
	err = commissionableDataProvider.Init(options)
	if err != nil {
		return err
	}
	configProvider := config.DefaultProvider()

	configManager := config.DefaultManager()
	err = configManager.Init(configProvider, options)
	if err != nil {
		return err
	}

	if options.Payload.RendezvousInformation != 0 {
		rendezvousFlags = options.Payload.RendezvousInformation
	}
	err = core.GetPayloadContents(&options.Payload, rendezvousFlags)
	if err != nil {
		return err
	}

	{
		options.Payload.CommissioningFlow = config.CommissioningFlowCustom
		core.PrintOnboardingCodes(options.Payload)

	}

	credentials.SetDeviceAttestationCredentialsProvider(nil)

	deviceInstanceInfo := device.DefaultInstanceInfo()
	err = deviceInstanceInfo.Init(configManager)
	if err != nil {

		return err
	}

	return nil
}

func MainLoop(options *config.DeviceOptions) error {

	serverInitParams := server.NewServerInitParams()
	_, err := serverInitParams.Init(options)
	if err != nil {

		return err
	}

	err = serverInitParams.InitializeStaticResourcesBeforeServerInit()
	if err != nil {

		return err
	}
	chipServer := server.NewServer()
	err = chipServer.Init(serverInitParams)
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
