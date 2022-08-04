package raw

type TCPTransportImpl struct {
	*BaseImpl
	mState uint8
	mPort  uint16
}
