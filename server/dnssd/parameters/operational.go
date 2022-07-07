package parameters

import "github.com/galenliu/chip/core"

type OperationalAdvertisingParameters struct {
	*BaseAdvertisingParams
	mPeerId *core.PeerId
}

func (o OperationalAdvertisingParameters) Init() *OperationalAdvertisingParameters {
	o.BaseAdvertisingParams = BaseAdvertisingParams{}.Init()
	return &o
}

func (o *OperationalAdvertisingParameters) SetPeerId(peerId *core.PeerId) {
	o.mPeerId = peerId
}

func (o *OperationalAdvertisingParameters) GetCompressedFabricId() core.CompressedFabricId {
	return o.mPeerId.GetCompressedFabricId()
}

func (o *OperationalAdvertisingParameters) GetPeerId() *core.PeerId {
	return o.mPeerId
}
