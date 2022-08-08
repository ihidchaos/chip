package raw

type TCPTransport interface {
	TransportBase
}

type TCPTransportImpl struct {
	mState uint8
	mPort  uint16
}
