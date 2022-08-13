package transport

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"net/netip"
)

type TransportManagerDelegate interface {
	OnMessageReceived(peerAddress netip.AddrPort, buf *lib.PacketBuffer)
}

// TransportManagerBase this is the delegate for TransportBase,
type TransportManagerBase interface {
	raw.TransportDelegate
	SetSessionManager(sessionManager TransportManagerDelegate)
	SendMessage(port netip.AddrPort, msg []byte) error
	Close()
	Disconnect(addr netip.Addr)
	MulticastGroupJoinLeave(addr netip.Addr, join bool) error
}

// TransportManagerImpl  impl TransportManagerBase
type TransportManagerImpl struct {
	mTransports     []raw.TransportBase
	mSessionManager TransportManagerDelegate
}

func NewTransportManagerImpl(transports ...raw.TransportBase) *TransportManagerImpl {
	impl := &TransportManagerImpl{
		mTransports: transports,
	}
	for _, t := range impl.mTransports {
		t.SetDelegate(impl)
	}
	return impl
}

func (t *TransportManagerImpl) HandleMessageReceived(peerAddress netip.AddrPort, buf *lib.PacketBuffer) {
	if t.mSessionManager != nil {
		t.mSessionManager.OnMessageReceived(peerAddress, buf)
	}
}

func (t *TransportManagerImpl) GetImplAtIndex(index int) raw.TransportBase {
	return t.mTransports[index]
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

func (t *TransportManagerImpl) SetSessionManager(sessionManager TransportManagerDelegate) {
	t.mSessionManager = sessionManager
}
