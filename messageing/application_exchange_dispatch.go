package messageing

import "sync/atomic"

type ApplicationExchangeDispatch struct {
	ExchangeMessageDispatchBase
}

var defaultApplicationExchangeDispatch atomic.Value

func init() {
	defaultEphemeralExchangeDispatch.Store(&ApplicationExchangeDispatch{})
}

func DefaultApplicationExchangeDispatch() *ApplicationExchangeDispatch {
	return defaultApplicationExchangeDispatch.Load().(*ApplicationExchangeDispatch)
}
