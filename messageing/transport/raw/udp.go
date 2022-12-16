package raw

import (
	"github.com/galenliu/chip/platform/system"
	log "golang.org/x/exp/slog"
	"io"
	"net"
	"net/netip"
)

const (
	NotReady uint8 = iota
	Initialized
)

// UDPTransportBase Must impl TransportBase
type UDPTransportBase interface {
	TransportBase
	BoundPort() uint16
	Init(parameters *UdpListenParameters) error
}

// UDPTransport impl TransportBase and add  init func
type UDPTransport struct {
	mDelegate Delegate
	mState    uint8
	mPort     uint16
}

func NewUdpTransport() *UDPTransport {
	return &UDPTransport{
		mState: 0,
	}
}

func (u *UDPTransport) Init(parameters *UdpListenParameters) error {

	if u.mState != NotReady {
		u.Close()
	}
	u.mPort = parameters.Address().Port()
	network := "udp6"
	if parameters.Address().Addr().Is4() {
		network = "udp4"
	}
	go func() {
		for {
			udpConn, err := net.ListenUDP(network, &net.UDPAddr{
				IP:   parameters.Address().Addr().AsSlice(),
				Port: int(parameters.Address().Port()),
			})
			if err != nil {
				log.Error("UDPTransportBase Connection", err, "NextTag", "UDPTransportBase")
				u.Close()
				return
			}
			go u.onConnection(udpConn, parameters.mAddress.Port())
		}
	}()
	return nil
}

func (u *UDPTransport) Disconnect(addr netip.AddrPort) {
	//TODO implement me
	panic("implement me")
}

func (u *UDPTransport) SendMessage(port netip.AddrPort, msg *system.PacketBufferHandle) error {
	//TODO implement me
	panic("implement me")
}

func (u *UDPTransport) MulticastGroupJoinLeave(addr netip.Addr, join bool) error {
	//TODO implement me
	panic("implement me")
}

func (u *UDPTransport) BoundPort() uint16 {
	return u.mPort
}

func (u *UDPTransport) SetDelegate(m Delegate) {
	u.mDelegate = m
}

func (u *UDPTransport) Close() {

}

func (u *UDPTransport) onConnection(conn *net.UDPConn, port uint16) {
	var data []byte
	_, err := io.ReadFull(conn, data)
	if err != nil {
		log.Error("read data", err, "NextTag", "UDPTransportBase")
		return
	}
	packetBuffer := system.NewPacketBufferHandle(data)
	if err := packetBuffer.IsValid(); err != nil {
		log.Error("invalid message", err, "NextTag", "UDPTransportBase")
		return
	}
	srcAddr, _ := netip.ParseAddr(conn.RemoteAddr().String())
	srcAddrPort := netip.AddrPortFrom(srcAddr, port)
	u.onUdpReceive(srcAddrPort, packetBuffer)
}

func (u *UDPTransport) onUdpReceive(srcAddr netip.AddrPort, data *system.PacketBufferHandle) {
	if u.mDelegate == nil {
		log.Warn("not delegate", "NextTag", "UDPTransportBase")
		return
	}
	u.mDelegate.HandleMessageReceived(srcAddr, data)
}
