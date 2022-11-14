package secure_channel

import (
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/platform/system"
)

type UnsolicitedStatusHandler interface {
	messageing.ExchangeDelegate
	messageing.UnsolicitedMessageHandler
	Init(mgr messageing.ExchangeManagerBase) error
}

func NewUnsolicitedStatusHandler() *UnsolicitedStatusHandlerImpl {
	return &UnsolicitedStatusHandlerImpl{}
}

type UnsolicitedStatusHandlerImpl struct {
	mExchangeManager messageing.ExchangeManagerBase
}

func (h UnsolicitedStatusHandlerImpl) OnMessageReceived(context *messageing.ExchangeContext,
	header *raw.PayloadHeader, data *system.PacketBufferHandle) error {
	//TODO implement me
	panic("implement me")
}

func (h UnsolicitedStatusHandlerImpl) OnResponseTimeout(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (h UnsolicitedStatusHandlerImpl) OnExchangeClosing(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (h UnsolicitedStatusHandlerImpl) OnUnsolicitedMessageReceived(header *raw.PayloadHeader, delegate messageing.ExchangeDelegate) error {
	//TODO implement me
	panic("implement me")
}

func (h UnsolicitedStatusHandlerImpl) OnExchangeCreationFailed(delegate messageing.ExchangeDelegate) {
	//TODO implement me
	panic("implement me")
}

func (h UnsolicitedStatusHandlerImpl) Init(mgr messageing.ExchangeManagerBase) error {
	h.mExchangeManager = mgr
	return nil
}
