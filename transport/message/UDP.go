package message

import (
	"github.com/galenliu/chip/transport"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/netip"
)

type UdpTransport interface {
	TransportBase
	Init(parameters *UdpListenParameters) error
}

type UdpTransportImpl struct {
	*BaseImpl
	mState uint8
}

func NewUdpTransportImpl() *UdpTransportImpl {
	return &UdpTransportImpl{
		BaseImpl: NewBaseImpl(),
		mState:   0,
	}
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

func (p *UdpTransportImpl) Init(parameters *UdpListenParameters) error {

	if p.mState != transport.NotReady {
		p.Close()
	}
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
			go p.readAll(udpConn, parameters.mAddress.Port())
		}
	}()
	return nil
}

func (p *UdpTransportImpl) readAll(conn *net.UDPConn, port uint16) {
	var data []byte
	_, err := io.ReadFull(conn, data)
	if err != nil {
		log.Error(err.Error())
		return
	}
	srcAddr, _ := netip.ParseAddr(conn.RemoteAddr().String())
	srcAddrPort := netip.AddrPortFrom(srcAddr, port)
	p.OnUdpReceive(srcAddrPort, data)
}

func (p *UdpTransportImpl) OnUdpReceive(srcAddr netip.AddrPort, data []byte) {
	p.HandleMessageReceived(srcAddr, data)
}

func (p *UdpTransportImpl) Close() {

}
