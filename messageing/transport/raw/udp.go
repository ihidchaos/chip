package raw

import (
	"github.com/galenliu/chip/lib/buffer"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/netip"
)

const (
	NotReady uint8 = iota
	Initialized
)

// UDPTransport Must impl TransportBase
type UDPTransport interface {
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

func (u *UDPTransportImpl) Disconnect(addr netip.Addr) {
	//TODO implement me
	panic("implement me")
}

func (u *UDPTransportImpl) SendMessage(port netip.AddrPort, msg []byte) error {
	//TODO implement me
	panic("implement me")
}

func (u *UDPTransportImpl) MulticastGroupJoinLeave(addr netip.Addr, join bool) error {
	//TODO implement me
	panic("implement me")
}

func (u *UDPTransportImpl) Init(parameters *UdpListenParameters) error {

	if u.mState != NotReady {
		u.Close()
	}
	u.mPort = parameters.GetAddress().Port()
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
				log.Infof("UDPTransport listen err : %s", err.Error())
				u.Close()
				return
			}
			go u.handelConnection(udpConn, parameters.mAddress.Port())
		}
	}()
	return nil
}

func (u *UDPTransportImpl) BoundPort() uint16 {
	return u.mPort
}

func (u *UDPTransportImpl) handelConnection(conn *net.UDPConn, port uint16) {
	var data []byte
	_, err := io.ReadFull(conn, data)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if data == nil {
		return
	}
	packetBuffer := buffer.NewPacketBuffer(data)
	srcAddr, _ := netip.ParseAddr(conn.RemoteAddr().String())
	srcAddrPort := netip.AddrPortFrom(srcAddr, port)
	u.OnUdpReceive(srcAddrPort, packetBuffer)
}

func (u *UDPTransportImpl) OnUdpReceive(srcAddr netip.AddrPort, data *buffer.PacketBuffer) {
	u.mDelegate.HandleMessageReceived(srcAddr, data)
}

func (u *UDPTransportImpl) HandleMessageReceived(addrPort netip.AddrPort, data *buffer.PacketBuffer) {
	u.mDelegate.HandleMessageReceived(addrPort, data)
}

func (u *UDPTransportImpl) SetDelegate(m TransportDelegate) {
	u.mDelegate = m
}

func (u *UDPTransportImpl) Close() {

}
