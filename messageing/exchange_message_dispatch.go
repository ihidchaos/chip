package messageing

import (
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/protocols"
	"sync/atomic"
)

type ExchangeMessageDispatchBase interface {
	IsEncryptionRequired() bool
	SendMessage(mgr transport.SessionManagerBase, handle *transport.SessionHandle)
	IsReliableTransmissionAllowed() bool
	MessagePermitted(id *protocols.Id, messageType uint8) bool
}

var defaultEphemeralExchangeDispatch atomic.Value

func init() {
	defaultEphemeralExchangeDispatch.Store(&EphemeraExchangeDispatch{})
}

func DefaultEphemeraExchangeDispatch() *EphemeralExchangeDispatch {
	return defaultEphemeralExchangeDispatch.Load().(*EphemeralExchangeDispatch)
}

type EphemeraExchangeDispatch struct {
	delegate ExchangeDelegate
}

func (d EphemeraExchangeDispatch) MessagePermitted(id *protocols.Id, messageType uint8) bool {
	//TODO implement me
	panic("implement me")
}

func (d EphemeraExchangeDispatch) IsReliableTransmissionAllowed() bool {
	//TODO implement me
	panic("implement me")
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

func (d EphemeraExchangeDispatch) IsEncryptionRequired() bool {
	return true
}

func (d EphemeralExchangeDispatch) IsEncryptionRequired() bool {
	return false
}

func (d EphemeraExchangeDispatch) SendMessage(mgr transport.SessionManagerBase, handle *transport.SessionHandle) {

}

func (d EphemeralExchangeDispatch) SendMessage(mgr transport.SessionManagerBase, handle *transport.SessionHandle) {

}
