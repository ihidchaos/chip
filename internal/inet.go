package internal

import (
	"net"
	"net/netip"
)

func GetAddress(i net.Interface) (ipAdders []netip.Addr) {
	adders, err := i.Addrs()
	if err != nil {
		return
	}
	for _, adder := range adders {

		ipAddr, _, err := net.ParseCIDR(adder.String())
		if err != nil {
			continue
		}
		a, e := netip.ParseAddr(ipAddr.String())
		if e != nil {
			continue
		}
		if !a.Is6() && !a.Is4() {
			continue
		}
		if !a.IsGlobalUnicast() {
			continue
		}
		ipAdders = append(ipAdders, a)
	}
	return
}
