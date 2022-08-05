package transport

import (
	"github.com/galenliu/chip/messageing/transport/raw"
	"net/netip"
)

type ManagerBase interface {
	raw.TransportDelegate
	SetSessionManager(sessionManager SessionManager)
	SendMessage(port netip.AddrPort, msg []byte) error
	Close()
	Disconnect(addr netip.Addr)
	MulticastGroupJoinLeave(addr netip.Addr, join bool) error
}

type ManagerImpl struct {
	mTransports     []raw.TransportBase
	mSessionManager SessionManager
}

func NewTransportManagerImpl(transports ...raw.TransportBase) *ManagerImpl {
	impl := &ManagerImpl{
		mTransports: transports,
	}
	for _, t := range impl.mTransports {
		t.SetDelegate(impl)
	}
	return impl
}

func (t *ManagerImpl) HandleMessageReceived(peerAddress netip.AddrPort, msg []byte) {
	if t.mSessionManager != nil {
		t.mSessionManager.OnMessageReceived(peerAddress, msg)
	}
}

func (t *ManagerImpl) GetImplAtIndex(index int) raw.TransportBase {
	return t.mTransports[index]
}

func (t *ManagerImpl) MulticastGroupJoinLeave(addr netip.Addr, join bool) error {
	//TODO implement me
	panic("implement me")
}

func (t *ManagerImpl) SendMessage(port netip.AddrPort, msg []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *ManagerImpl) Close() {
	//TODO implement me
	panic("implement me")
}

func (t *ManagerImpl) Disconnect(addr netip.Addr) {
	//TODO implement me
	panic("implement me")
}

func (t *ManagerImpl) SetSessionManager(sessionManager SessionManager) {
	t.mSessionManager = sessionManager
}
