package dnssd

import (
	"fmt"
	log "golang.org/x/exp/slog"
	"net"
	"strconv"
	"testing"
)

func TestBase(t *testing.T) {
	base := BaseAdvertisingParams{}
	ic, _ := net.Interfaces()
	add, _ := net.InterfaceAddrs()
	for _, it := range ic {
		log.Info("----------")
		log.Info(strconv.Itoa(it.Index))
		log.Info(it.HardwareAddr.String())
		if len(it.HardwareAddr) > 5 {
			log.Info("%d", it.HardwareAddr[0])
			log.Info("%d", it.HardwareAddr[1])
			log.Info("%d", it.HardwareAddr[2])
			log.Info("%d", it.HardwareAddr[3])
			log.Info("%d", it.HardwareAddr[4])
			log.Info("%d", it.HardwareAddr[5])
		}

		log.Info("len= %d", len(it.HardwareAddr))

	}
	for _, a := range add {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				log.Info("+++++++++")
				fmt.Println(ipnet.IP.String())
				fmt.Println(ipnet.Mask)
				fmt.Println(ipnet.Network())

			}
		}

		log.Info("----------")
		log.Info(a.Network())
		log.Info(a.String())
	}
	for i := 0; i < 20; i++ {
		t.Log(base.GetMac())
	}
}
