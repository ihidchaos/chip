package message

import "net/netip"

type RawTransportDelegate interface {
	HandleMessageReceived(addrPort netip.AddrPort, msg []byte)
}

type TransportBase interface {
	HandleMessageReceived(addrPort netip.AddrPort, msg []byte)
	SetDelegate(m RawTransportDelegate)
}

type BaseImpl struct {
	mDelegate RawTransportDelegate
}

func (b *BaseImpl) HandleMessageReceived(addrPort netip.AddrPort, msg []byte) {
	b.mDelegate.HandleMessageReceived(addrPort, msg)
}

func NewBaseImpl() *BaseImpl {
	return &BaseImpl{}
}

func (b *BaseImpl) SetDelegate(m RawTransportDelegate) {
	b.mDelegate = m
}
