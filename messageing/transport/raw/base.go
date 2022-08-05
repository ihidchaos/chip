package raw

import "net/netip"

type TransportDelegate interface {
	HandleMessageReceived(addrPort netip.AddrPort, msg []byte)
}

type TransportBase interface {
	GetBoundPort() uint16
	HandleMessageReceived(addrPort netip.AddrPort, msg []byte)
	SetDelegate(m TransportDelegate)
}

type BaseImpl struct {
	mDelegate TransportDelegate
}

func (b *BaseImpl) HandleMessageReceived(addrPort netip.AddrPort, msg []byte) {
	b.mDelegate.HandleMessageReceived(addrPort, msg)
}

func NewBaseImpl() *BaseImpl {
	return &BaseImpl{}
}

func (b *BaseImpl) SetDelegate(m TransportDelegate) {
	b.mDelegate = m
}
