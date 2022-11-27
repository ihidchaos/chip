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
	MessagePermitted(id protocols.Id, messageType uint8) bool
}

var defaultEphemeralExchangeDispatch atomic.Value

func init() {
	defaultEphemeralExchangeDispatch.Store(&EphemeraExchangeDispatch{nil})
}

func DefaultEphemeraExchangeDispatch() *EphemeraExchangeDispatch {
	return defaultEphemeralExchangeDispatch.Load().(*EphemeraExchangeDispatch)
}

type EphemeraExchangeDispatch struct {
	delegate ExchangeDelegate
}

func (d EphemeraExchangeDispatch) SetAckPending(b bool) {
	//TODO implement me
	panic("implement me")
}

func (d EphemeraExchangeDispatch) MessagePermitted(id protocols.Id, messageType uint8) bool {
	//TODO implement me
	panic("implement me")
}

func (d EphemeraExchangeDispatch) IsReliableTransmissionAllowed() bool {
	//TODO implement me
	panic("implement me")
}

func (d EphemeraExchangeDispatch) IsEncryptionRequired() bool {
	return true
}

func (d EphemeraExchangeDispatch) SendMessage(mgr transport.SessionManagerBase, handle *transport.SessionHandle) {

}
