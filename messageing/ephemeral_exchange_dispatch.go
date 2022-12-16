package messageing

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/protocols"
)

type EphemeralExchangeDispatch struct {
	delegate ExchangeDelegate
}

func (e *EphemeralExchangeDispatch) SendMessage(
	mgr *transport.SessionManager,
	handle *transport.SessionHandle,
	exchangeId uint16,
	isInitiator bool,
	rmc *ReliableMessageContext,
	isReliableTransmission bool,
	protocol protocols.Id,
	msgType uint8,
	message []byte) error {
	return SendMessage(e, mgr, handle, exchangeId, isInitiator, rmc, isReliableTransmission, protocol, msgType, message)
}

func (e *EphemeralExchangeDispatch) IsEncryptionRequired() bool {
	return false
}

func (e *EphemeralExchangeDispatch) IsReliableTransmissionAllowed() bool {
	return true
}

func (e *EphemeralExchangeDispatch) MessagePermitted(protocol protocols.Id, messageType uint8) bool {
	return lib.IsSecureChannel(protocol) && lib.IsStandaloneAck(messageType)
}
