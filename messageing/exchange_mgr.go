package messageing

import "github.com/galenliu/chip/transport"

type ExchangeManager interface {
	Init(sessions transport.SessionManager) error
}

type ExchangeManagerImpl struct {
}
