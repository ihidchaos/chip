package app

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/server"
	"github.com/galenliu/gateway/pkg/log"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestDnssdServer(t *testing.T) {
	initParams := server.InitializeStaticResourcesBeforeServerInit()
	initParams.OperationalServicePort = config.ChipPort
	initParams.UserDirectedCommissioningPort = config.ChipUdcPort
	chip := server.Server{}.Init(initParams)
	err := chip.StartServer()
	if err != nil {
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGUSR1)

	for {
		s := <-ch
		switch s {
		default:
			log.Info("signal exit")
			return
		}
	}

}
