package transport

import (
	log "github.com/sirupsen/logrus"
	"net"
	"net/netip"
)

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

type UdpTransportImpl struct {
}

func (p *UdpTransportImpl) Disconnect(addr netip.Addr) {
	//TODO implement me
	panic("implement me")
}

func (p *UdpTransportImpl) SendMessage(port netip.AddrPort, msg []byte) error {
	//TODO implement me
	panic("implement me")
}

func (p *UdpTransportImpl) HandleMessageReceived(peerAddress netip.AddrPort, msg []byte) {
	//TODO implement me
	panic("implement me")
}

func (p *UdpTransportImpl) MulticastGroupJoinLeave(addr netip.Addr, join bool) error {
	//TODO implement me
	panic("implement me")
}

func NewUdbTransportImpl() *UdpTransportImpl {
	return &UdpTransportImpl{}
}

func (p *UdpTransportImpl) Init(parameters *UdpListenParameters) error {
	network := "udp6"
	if parameters.GetAddress().Addr().Is4() {
		network = "udp4"
	}
	go func() {
		for {
			udpConn, err := net.ListenUDP(network, &net.UDPAddr{
				IP:   parameters.GetAddress().Addr().AsSlice(),
				Port: int(parameters.GetAddress().Port()),
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
