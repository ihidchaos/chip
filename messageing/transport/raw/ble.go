package raw

type BLETransportBase interface {
	TransportBase
}

type BLETransport struct {
	mState uint8
	mPort  uint16
}
