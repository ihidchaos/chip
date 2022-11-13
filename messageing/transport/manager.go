package transport

import (
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/platform/system"
	"net/netip"
)

// MgrDelegate 这个实例为 SessionManagerBase
type MgrDelegate interface {
	OnMessageReceived(peerAddress netip.AddrPort, buf *system.PacketBufferHandle)
}

// MgrBase  this is the delegate for TransportBase,
type MgrBase interface {
	raw.Delegate
	SetSessionManager(sessionManager MgrDelegate)
	SendMessage(port netip.AddrPort, msg []byte) error
	Close()
	Disconnect(addr netip.Addr)
	MulticastGroupJoinLeave(addr netip.Addr, join bool) error
}

// Manager  impl ManagerBase
type Manager struct {
	mTransports     []raw.TransportBase
	mSessionManager MgrDelegate
}

func NewManager(transports ...raw.TransportBase) *Manager {
	impl := &Manager{
		mTransports: transports,
	}
	for _, t := range impl.mTransports {
		t.SetDelegate(impl)
	}
	return impl
}

func (t *Manager) HandleMessageReceived(peerAddress netip.AddrPort, buf *system.PacketBufferHandle) {
	if t.mSessionManager != nil {
		t.mSessionManager.OnMessageReceived(peerAddress, buf)
	}
}

func (t *Manager) GetImplAtIndex(index int) raw.TransportBase {
	return t.mTransports[index]
}

func (t *Manager) GetUpdImpl() raw.UDPTransportBase {
	for _, transport := range t.mTransports {
		udpTransport, ok := transport.(raw.UDPTransportBase)
		if ok {
			return udpTransport
		}
	}
	return nil
}

func (t *Manager) MulticastGroupJoinLeave(addr netip.Addr, join bool) error {
	//TODO implement me
	panic("implement me")
}

func (t *Manager) SendMessage(port netip.AddrPort, msg []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *Manager) Close() {
	//TODO implement me
	panic("implement me")
}

func (t *Manager) Disconnect(addr netip.Addr) {
	//TODO implement me
	panic("implement me")
}

func (t *Manager) SetSessionManager(sessionManager MgrDelegate) {
	t.mSessionManager = sessionManager
}
