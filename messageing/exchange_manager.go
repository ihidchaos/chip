package messageing

import (
	"github.com/galenliu/chip/protocols"
	"github.com/galenliu/chip/secure_channel"
	"github.com/galenliu/chip/transport"
)

type ExchangeManager interface {
	Init(sessions transport.SessionManager) error
	RegisterUnsolicitedMessageHandlerForType(sigma1 uint8, s *secure_channel.CASEServer)
}

type ExchangeManagerImpl struct {
}

func (e *ExchangeManagerImpl) RegisterUnsolicitedMessageHandlerForType(msgType uint8, handler *secure_channel.CASEServer) error {
	return e.registerUnsolicitedMessageHandlerForType(protocols.NotSpecifiedId(), msgType, handler)
}

func NewExchangeManagerImpl() *ExchangeManagerImpl {
	return &ExchangeManagerImpl{}
}

func (e *ExchangeManagerImpl) Init(sessions transport.SessionManager) error {
	return nil
}

func (e *ExchangeManagerImpl) registerUnsolicitedMessageHandlerForType(id protocols.Id, msgType uint8, handler *secure_channel.CASEServer) error {
	return nil
}
