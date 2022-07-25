package transport

import (
	"github.com/galenliu/chip/inet/udp_endpoint"
	"net/netip"
)

type UdpListenParameters struct {
	mAddr         netip.Addr
	mPort         uint16
	mNativeParams func()
}

func (p *UdpListenParameters) SetListenPort(port uint16) {
	p.mPort = port
}

func (p *UdpListenParameters) SetNativeParams(params func()) {
	p.mNativeParams = params
}

type UdpTransport struct {
}

func NewUdpTransport(mgr udp_endpoint.UDPEndpoint, params UdpListenParameters) (*UdpTransport, error) {
	return nil, nil
}
