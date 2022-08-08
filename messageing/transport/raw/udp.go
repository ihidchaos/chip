package raw

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/netip"
)

const (
	NotReady uint8 = iota
	Initialized
)

// UDBTransport Must impl TransportBase
type UDBTransport interface {
	TransportBase
	Init(parameters *UdpListenParameters) error
}

// UDPTransportImpl impl TransportBase and add  init func
type UDPTransportImpl struct {
	mDelegate TransportDelegate
	mState    uint8
	mPort     uint16
}

func NewUdpTransportImpl() *UDPTransportImpl {
	return &UDPTransportImpl{
		mState: 0,
	}
}

func (p *UDPTransportImpl) Disconnect(addr netip.Addr) {
	//TODO implement me
	panic("implement me")
}

func (p *UDPTransportImpl) SendMessage(port netip.AddrPort, msg []byte) error {
	//TODO implement me
	panic("implement me")
}

func (p *UDPTransportImpl) MulticastGroupJoinLeave(addr netip.Addr, join bool) error {
	//TODO implement me
	panic("implement me")
}

func (p *UDPTransportImpl) Init(parameters *UdpListenParameters) error {

	if p.mState != NotReady {
		p.Close()
	}
	p.mPort = parameters.GetAddress().Port()
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
				log.Error("UDPTransport err : %s", err.Error())
				p.Close()
				return
			}
			go p.handelConnection(udpConn, parameters.mAddress.Port())
		}
	}()
	return nil
}

func (p *UDPTransportImpl) GetBoundPort() uint16 {
	return p.mPort
}

func (p *UDPTransportImpl) handelConnection(conn *net.UDPConn, port uint16) {
	var data []byte
	_, err := io.ReadFull(conn, data)
	if err != nil {
		log.Error(err.Error())
		return
	}
	packetBuffer := NewPacketBuffer(data)
	srcAddr, _ := netip.ParseAddr(conn.RemoteAddr().String())
	srcAddrPort := netip.AddrPortFrom(srcAddr, port)
	p.OnUdpReceive(srcAddrPort, packetBuffer)
}

func (p *UDPTransportImpl) OnUdpReceive(srcAddr netip.AddrPort, data *PacketBuffer) {
	p.HandleMessageReceived(srcAddr, data)
}

func (b *UDPTransportImpl) HandleMessageReceived(addrPort netip.AddrPort, data *PacketBuffer) {
	b.mDelegate.HandleMessageReceived(addrPort, data)
}

func (b *UDPTransportImpl) SetDelegate(m TransportDelegate) {
	b.mDelegate = m
}

func (p *UDPTransportImpl) Close() {

}
