package dnssd

import (
	"github.com/galenliu/chip/lib"
)

type OperationalAdvertisingParameters struct {
	BaseAdvertisingParams
	mPeerId PeerId
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

func (o *OperationalAdvertisingParameters) SetPeerId(peerId PeerId) *OperationalAdvertisingParameters {
	o.mPeerId = peerId
	return o
}

func (o *OperationalAdvertisingParameters) GetCompressedFabricId() lib.CompressedFabricId {
	return o.mPeerId.GetCompressedFabricId()
}

func (o *OperationalAdvertisingParameters) GetPeerId() PeerId {
	return o.mPeerId
}
