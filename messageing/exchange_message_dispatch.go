package messageing

import "github.com/galenliu/chip/messageing/transport"

type ExchangeMessageDispatch interface {
	IsEncryptionRequired() bool
	SendMessage(mgr transport.SessionManager, handle transport.SessionHandle)
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

func (d ExchangeMessageDispatchImpl) SendMessage(mgr transport.SessionManager, handle transport.SessionHandle) {

}

func (d EphemeralExchangeDispatchImpl) SendMessage(mgr transport.SessionManager, handle transport.SessionHandle) {

}
