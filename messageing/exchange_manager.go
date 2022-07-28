package messageing

import "github.com/galenliu/chip/transport"

type ExchangeManager interface {
	Init(sessions transport.SessionManager) error
}

type ExchangeManagerImpl struct {
}

func NewExchangeManagerImpl() *ExchangeManagerImpl {
	return &ExchangeManagerImpl{}
}

func (e ExchangeManagerImpl) Init(sessions transport.SessionManager) error {
	return nil
}
