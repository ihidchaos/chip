package messageing

type ExchangeMessageDispatch interface {
	IsEncryptionRequired() bool
}

type ExchangeMessageDispatchImpl struct {
	delegate ExchangeDelegate
}

type EphemeralExchangeDispatchImpl struct {
	delegate ExchangeDelegate
}

func (d ExchangeMessageDispatchImpl) IsEncryptionRequired() bool {
	return true
}

func (d EphemeralExchangeDispatchImpl) IsEncryptionRequired() bool {
	return false
}
