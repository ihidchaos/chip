package transport

import (
	log "github.com/sirupsen/logrus"
	"net"
	"net/netip"
)

type UdpTransport interface {
	GetBoundPort() uint16
	Close()
	SendMessage(port netip.AddrPort, msg []byte) error
	MulticastGroupJoinLeave(port netip.AddrPort, joined bool) error
	CanListenMulticast() bool
	CanSendToPeer(port netip.AddrPort) bool
	OnUdpReceive(port netip.AddrPort)
	OnUdpError(endPoint netip.AddrPort)
}

type UdpTransportImpl struct {
}

func NewUdbTransportImpl() *UdpTransportImpl {
	return &UdpTransportImpl{}
}

func (p *UdpTransportImpl) Init(addr netip.AddrPort) error {
	network := "udp6"
	if addr.Addr().Is4() {
		network = "udp4"
	}
	go func() {
		for {
			udpConn, err := net.ListenUDP(network, &net.UDPAddr{
				IP:   addr.Addr().AsSlice(),
				Port: int(addr.Port()),
			})
			if err != nil {
				log.Error("UdpTransport err : %s", err.Error())
				p.Close()
				return
			}
			go p.ReadConnection(udpConn)
		}
	}()
	return nil
}

func (p *UdpTransportImpl) ReadConnection(conn *net.UDPConn) {

}

func (p *UdpTransportImpl) Close() {

}
