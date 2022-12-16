package secure_channel

import (
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/protocols"
	"sync/atomic"
)

type SessionEstablishmentExchangeDispatch struct {
}

func (s SessionEstablishmentExchangeDispatch) IsEncryptionRequired() bool {
	//TODO implement me
	panic("implement me")
}

func (s SessionEstablishmentExchangeDispatch) SendMessage(
	mgr *transport.SessionManager,
	handle *transport.SessionHandle,
	exchangeId uint16,
	isInitiator bool,
	rmc *messageing.ReliableMessageContext,
	isReliableTransmission bool,
	protocol protocols.Id,
	msgType uint8,
	message []byte) error {
	return messageing.SendMessage(s, mgr, handle, exchangeId, isInitiator, rmc, isReliableTransmission, protocol, msgType, message)
}

func (s SessionEstablishmentExchangeDispatch) IsReliableTransmissionAllowed() bool {
	//TODO implement me
	panic("implement me")
}

func (s SessionEstablishmentExchangeDispatch) MessagePermitted(id protocols.Id, typ uint8) bool {
	if id.Equal(protocols.New(protocolId, nil)) {
		switch MsgType(typ) {
		case StandaloneAck, PBKDFParamRequest, PBKDFParamResponse, PASEPake1, PASEPake2, PASEPake3, CASESigma1, CASESigma2, CASESigma3, CASESigma2Resume, CASEStatusReport:
			return true
		default:
			return false
		}
	}
	return false
}

func (s SessionEstablishmentExchangeDispatch) isEncryptionRequired() bool {
	return false
}

var defaultSessionEstablishmentExchangeDispatch atomic.Value

func init() {
	defaultSessionEstablishmentExchangeDispatch.Store(&SessionEstablishmentExchangeDispatch{})
}

func SessionEstablishmentExchangeDispatchInstance() *SessionEstablishmentExchangeDispatch {
	return defaultSessionEstablishmentExchangeDispatch.Load().(*SessionEstablishmentExchangeDispatch)
}
