package secure_channel

import "github.com/galenliu/chip/messageing"

type MessageCounterManager struct {
}

func (m MessageCounterManager) Init(mgr messageing.ExchangeManager) error {
	return nil
}

func NewMessageCounterManager() *MessageCounterManager {
	return &MessageCounterManager{}
}
