package params

import (
	"github.com/galenliu/chip/device"
)

type OperationalAdvertisingParameters struct {
	BaseAdvertisingParams
	mPeerId device.PeerId
}

func NewOperationalAdvertisingParameters() *OperationalAdvertisingParameters {
	return &OperationalAdvertisingParameters{
		BaseAdvertisingParams: BaseAdvertisingParams{},
	}
}

func (o OperationalAdvertisingParameters) Init() *OperationalAdvertisingParameters {
	o.BaseAdvertisingParams = BaseAdvertisingParams{}.Init()
	return &o
}

func (o *OperationalAdvertisingParameters) SetPeerId(peerId device.PeerId) *OperationalAdvertisingParameters {
	o.mPeerId = peerId
	return o
}

func (o *OperationalAdvertisingParameters) GetCompressedFabricId() device.CompressedFabricId {
	return o.mPeerId.GetCompressedFabricId()
}

func (o *OperationalAdvertisingParameters) GetPeerId() device.PeerId {
	return o.mPeerId
}
