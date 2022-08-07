package secure_channel

import (
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport/raw"
)

type UnsolicitedStatusHandler interface {
	OnMessageReceived(ec *messageing.ExchangeContext, header *raw.PayloadHeader, data []byte) error
	OnResponseTimeout(ec *messageing.ExchangeContext)
	OnUnsolicitedMessageReceived(header *raw.PayloadHeader, delegate messageing.ExchangeDelegate) error
}

type UnsolicitedStatusHandlerImpl struct {
}

func (h UnsolicitedStatusHandlerImpl) Init(mgr messageing.ExchangeManager) error {
	return nil
}

func NewUnsolicitedStatusHandler() *UnsolicitedStatusHandlerImpl {
	return &UnsolicitedStatusHandlerImpl{}
}
