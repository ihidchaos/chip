package secure_channel

import "github.com/galenliu/chip/messageing"

type UnsolicitedStatusHandler struct {
}

func (h UnsolicitedStatusHandler) Init(mgr messageing.ExchangeManager) error {
	return nil
}

func NewUnsolicitedStatusHandler() *UnsolicitedStatusHandler {
	return &UnsolicitedStatusHandler{}
}
