package messageing

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/protocols"
)

type ApplicationExchangeDispatch struct {
}

func (a *ApplicationExchangeDispatch) SendMessage(
	mgr *transport.SessionManager,
	handle *transport.SessionHandle,
	exchangeId uint16,
	isInitiator bool,
	rmc *ReliableMessageContext,
	isReliableTransmission bool,
	protocol protocols.Id,
	msgType uint8,
	message []byte) error {
	return SendMessage(a, mgr, handle, exchangeId, isInitiator, rmc, isReliableTransmission, protocol, msgType, message)
}

func (a *ApplicationExchangeDispatch) IsReliableTransmissionAllowed() bool {
	return true
}

func (a *ApplicationExchangeDispatch) MessagePermitted(protocol protocols.Id, messageType uint8) bool {
	return !lib.IsSecureChannel(protocol) && !lib.IsSecureMessage(messageType)
}

func (a *ApplicationExchangeDispatch) IsEncryptionRequired() bool {
	return true
}
