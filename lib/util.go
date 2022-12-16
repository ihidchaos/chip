package lib

import (
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/protocols"
	"github.com/galenliu/chip/protocols/secure_channel"
)

func IsControlMessage(payloadHeader *raw.PayloadHeader) bool {
	return payloadHeader.HasMessageType(secure_channel.MsgCounterSyncReq.MessageType()) || payloadHeader.HasMessageType(secure_channel.MsgCounterSyncRsp.MessageType())
}

func IsStandaloneAck(msg uint8) bool {
	return msg == secure_channel.StandaloneAck.MessageType()
}

func IsSecureChannel(id protocols.Id) bool {
	return id == secure_channel.Id
}

func IsSecureMessage(msgType uint8) bool {
	switch secure_channel.MsgType(msgType) {
	case secure_channel.PBKDFParamRequest, secure_channel.PBKDFParamResponse,
		secure_channel.PASEPake1, secure_channel.PASEPake2, secure_channel.PASEPake3,
		secure_channel.CASESigma1, secure_channel.CASESigma2, secure_channel.CASESigma3,
		secure_channel.CASESigma2Resume:
		return true
	default:
		return false
	}
}
