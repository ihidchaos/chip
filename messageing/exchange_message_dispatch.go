package messageing

import (
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/protocols"
)

type ExchangeMessageDispatchBase interface {
	IsEncryptionRequired() bool
	SendMessage(mgr transport.SessionManagerBase, handle *transport.SessionHandle)
	IsReliableTransmissionAllowed() bool
	MessagePermitted(id *protocols.Id, messageType uint8) bool
}

type ExchangeMessageDispatchImpl struct {
	delegate ExchangeDelegate
}

func (d ExchangeMessageDispatchImpl) MessagePermitted(id *protocols.Id, messageType uint8) bool {
	//TODO implement me
	panic("implement me")
}

func (d ExchangeMessageDispatchImpl) IsReliableTransmissionAllowed() bool {
	//TODO implement me
	panic("implement me")
}

func EphemeralExchangeDispatchInstance() *EphemeralExchangeDispatch {
	return nil
}

type EphemeralExchangeDispatch struct {
	delegate ExchangeDelegate
}

func (d EphemeralExchangeDispatch) SetAckPending(b bool) {
	//TODO implement me
	panic("implement me")
}

func (d EphemeralExchangeDispatch) MessagePermitted(id *protocols.Id, messageType uint8) bool {
	//TODO implement me
	panic("implement me")
}

func (d EphemeralExchangeDispatch) IsReliableTransmissionAllowed() bool {
	//TODO implement me
	panic("implement me")
}

func (d ExchangeMessageDispatchImpl) IsEncryptionRequired() bool {
	return true
}

func (d EphemeralExchangeDispatch) IsEncryptionRequired() bool {
	return false
}

func (d ExchangeMessageDispatchImpl) SendMessage(mgr transport.SessionManagerBase, handle *transport.SessionHandle) {

}

func (d EphemeralExchangeDispatch) SendMessage(mgr transport.SessionManagerBase, handle *transport.SessionHandle) {

}
