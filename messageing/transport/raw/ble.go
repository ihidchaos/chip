package raw

type BLETransport interface {
	TransportBase
}

type BLETransportImpl struct {
	mState uint8
	mPort  uint16
}
