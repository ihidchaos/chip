package raw

import "net/netip"

type UdpListenParameters struct {
	mAddress netip.AddrPort
}

func (l *UdpListenParameters) SetAddress(address netip.AddrPort) *UdpListenParameters {
	l.mAddress = address
	return l
}

func (l *UdpListenParameters) GetAddress() netip.AddrPort {
	return l.mAddress
}
