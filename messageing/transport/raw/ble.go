package raw

type BLETransportImpl struct {
	*BaseImpl
	mState uint8
	mPort  uint16
}
