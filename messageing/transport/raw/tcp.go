package raw

type TCPTransportBase interface {
	TransportBase
}

type TCPTransport struct {
	mState uint8
	mPort  uint16
}
