package transport

import (
	"github.com/galenliu/chip/transport/message"
	"net/netip"
)

const (
	NotReady uint8 = iota
	Initialized
)

type TransportManager interface {
	message.TransportBase
	message.RawTransportDelegate
	SetSessionManager(sessionManager SessionManager)
	MulticastGroupJoinLeave(addr netip.Addr, join bool) error
	SendMessage(port netip.AddrPort, msg []byte) error
	Close()
	Disconnect(addr netip.Addr)
}

type TransportManagerImpl struct {
	message.TransportBase
	mSessionManager SessionManager
}

func NewTransportManagerImpl(base message.TransportBase) *TransportManagerImpl {
	impl := &TransportManagerImpl{
		TransportBase: base,
	}
	impl.SetDelegate(impl)
	return impl
}

func (t *TransportManagerImpl) HandleMessageReceived(peerAddress netip.AddrPort, msg []byte) {
	if t.mSessionManager != nil {
		t.mSessionManager.OnMessageReceived(peerAddress, msg)
	}
}

func (t *TransportManagerImpl) MulticastGroupJoinLeave(addr netip.Addr, join bool) error {
	//TODO implement me
	panic("implement me")
}

func (t *TransportManagerImpl) SendMessage(port netip.AddrPort, msg []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *TransportManagerImpl) Close() {
	//TODO implement me
	panic("implement me")
}

func (t *TransportManagerImpl) Disconnect(addr netip.Addr) {
	//TODO implement me
	panic("implement me")
}

func (t *TransportManagerImpl) SetSessionManager(sessionManager SessionManager) {
	t.mSessionManager = sessionManager
}
