package messageing

import (
	"github.com/galenliu/chip"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/protocols"
	"sync/atomic"
)

type ExchangeMessageDispatch interface {
	IsEncryptionRequired() bool
	SendMessage(mgr *transport.SessionManager, handle *transport.SessionHandle, exchangeId uint16, isInitiator bool, rmc *ReliableMessageContext, isReliableTransmission bool, protocol protocols.Id, msgType uint8, message []byte) error
	IsReliableTransmissionAllowed() bool
	MessagePermitted(id protocols.Id, messageType uint8) bool
}

var defaultEphemeralExchangeDispatch atomic.Value
var defaultApplicationExchangeDispatch atomic.Value

func init() {
	defaultEphemeralExchangeDispatch.Store(&EphemeralExchangeDispatch{nil})
	defaultApplicationExchangeDispatch.Store(&ApplicationExchangeDispatch{})
}

func DefaultEphemeralDispatch() *EphemeralExchangeDispatch {
	return defaultEphemeralExchangeDispatch.Load().(*EphemeralExchangeDispatch)
}

func DefaultApplicationExchangeDispatch() *ApplicationExchangeDispatch {
	return defaultApplicationExchangeDispatch.Load().(*ApplicationExchangeDispatch)
}

func SendMessage(dispatch ExchangeMessageDispatch,
	sessionManager *transport.SessionManager,
	sessionHandle *transport.SessionHandle,
	exchangeId uint16,
	isInitiator bool,
	reliableMessageContext *ReliableMessageContext,
	isReliableTransmission bool,
	protocolsId protocols.Id,
	msgType uint8,
	message []byte) error {

	if !dispatch.MessagePermitted(protocolsId, msgType) {
		return chip.New(chip.ErrorInvalidArgument)
	}
	var payloadHeader = raw.NewPayloadHeader()
	payloadHeader.ExchangeId = exchangeId
	payloadHeader.SetMessageType(protocolsId, msgType).SetInitiator(isInitiator)

	if reliableMessageContext.HasPiggybackAckPending() {
		payloadHeader.SetAckMessageCounter(reliableMessageContext.TakePendingPeerAckMessageCounter())
	}
	if dispatch.IsReliableTransmissionAllowed() &&
		reliableMessageContext.AutoRequestAck() &&
		reliableMessageContext.ReliableMessageMgr() != nil &&
		isReliableTransmission {
		reliableMessageMgr := reliableMessageContext.ReliableMessageMgr()
		payloadHeader.SetNeedsAck(true)

		_, err := reliableMessageMgr.AddToRetransTable(reliableMessageContext)
		if err != nil {
			return err
		}

		//deleter := func(mgr *ReliableMessageMgr, e *RetransTableEntry) {
		//	mgr.ClearRetransTableEntry(e)
		//}
		_, err = sessionManager.PrepareMessage(sessionHandle, payloadHeader, message)
		if err != nil {
			return err
		}

	}
	return nil
}
