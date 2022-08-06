package messageing

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/pkg"
	"github.com/galenliu/chip/protocols"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

const (
	kStateNotInitialized = 0
	kStateInitialized    = 1
)
const KAnyMessageType int16 = -1

type UnsolicitedMessageHandlerSlot struct {
	MessageType int16
	Handler     raw.UnsolicitedMessageHandler
	ProtocolId  *protocols.Id
}

func (receiver UnsolicitedMessageHandlerSlot) Matches(aProtocolId *protocols.Id, aMessageType int16) bool {
	return aProtocolId == receiver.ProtocolId && aMessageType == receiver.MessageType
}

func (receiver UnsolicitedMessageHandlerSlot) Reset() {
	receiver.Handler = nil
}

type ExchangeManager interface {
	transport.SessionMessageDelegate
	RegisterUnsolicitedMessageHandlerForProtocol(protocolId *protocols.Id, handler UnsolicitedMessageHandler) error
	RegisterUnsolicitedMessageHandlerForType(protocolId *protocols.Id, msgType uint8, handler UnsolicitedMessageHandler) error
}

// ExchangeManagerImpl ExchangeManager
type ExchangeManagerImpl struct {
	UMHandlerPool       []*UnsolicitedMessageHandlerSlot
	mContextPool        []*ExchangeContext
	mSessionManager     transport.SessionManager
	mNextExchangeId     uint16
	mNextKeyId          uint16
	mReliableMessageMgr any //TODO
	mInitialized        bool
}

func NewExchangeManagerImpl() *ExchangeManagerImpl {
	impl := &ExchangeManagerImpl{}
	impl.UMHandlerPool = make([]*UnsolicitedMessageHandlerSlot, 0)
	impl.mContextPool = make([]*ExchangeContext, 0)
	return impl
}

func (e *ExchangeManagerImpl) Init(sessionManager transport.SessionManager) error {
	if e.mInitialized {
		return lib.ChipErrorIncorrectState
	}
	e.mSessionManager = sessionManager
	e.mNextExchangeId = uint16(rand.Uint32())
	e.mNextKeyId = 0

	for _, handler := range e.UMHandlerPool {
		handler.Reset()
	}
	sessionManager.SetMessageDelegate(e)
	//e.mReliableMessageMgr.init(sessionManager.SystemLayer())
	e.mInitialized = true
	return nil
}

func (e *ExchangeManagerImpl) OnMessageReceived(packetHeader *raw.PacketHeader, payloadHeader *raw.PayloadHeader, session transport.Session, isDuplicate uint8, data []byte) {
	//var matchingUMH *UnsolicitedMessageHandlerSlot = nil

	log.Infof("Received message of type %d with protocolId %d"+
		"and MessageCounter: %d",
		payloadHeader.GetMessageType(), payloadHeader.GetProtocolID(),
		packetHeader.GetMessageCounter())
	var msgFlags uint32 = 0
	if isDuplicate == transport.DuplicateMessageYes {
		msgFlags = msgFlags | pkg.KDuplicateMessage
	}
	found := false
	if !packetHeader.IsGroupSession() {
		for _, ec := range e.mContextPool {
			if ec.MatchExchange(session, packetHeader, payloadHeader) {
				ec.HandleMessage(packetHeader.GetMessageCounter(), payloadHeader, msgFlags, data)
				found = true
				return
			}
			continue
		}
	} else {
		if found {
			return
		}
	}
}

func (e *ExchangeManagerImpl) RegisterUnsolicitedMessageHandlerForProtocol(
	protocolId *protocols.Id,
	handler UnsolicitedMessageHandler,
) error {

	return e.RegisterUMH(protocolId, KAnyMessageType, handler)

}

func (e *ExchangeManagerImpl) RegisterUnsolicitedMessageHandlerForType(
	protocolId *protocols.Id,
	msgType uint8,
	handler UnsolicitedMessageHandler,
) error {

	return e.RegisterUMH(protocolId, int16(msgType), handler)

}

func (e *ExchangeManagerImpl) RegisterUMH(id *protocols.Id, msgType int16, handle raw.UnsolicitedMessageHandler) error {

	for _, umh := range e.UMHandlerPool {
		if umh.Matches(id, msgType) {
			umh.Handler = handle
			return nil
		}
	}
	e.UMHandlerPool = append(e.UMHandlerPool, &UnsolicitedMessageHandlerSlot{
		MessageType: msgType,
		Handler:     handle,
		ProtocolId:  id,
	})
	return nil
}

func (e *ExchangeManagerImpl) UnregisterUMH(id *protocols.Id, msgType int16) error {
	for i, umh := range e.UMHandlerPool {
		if umh.Matches(id, msgType) {
			e.UMHandlerPool = append(e.UMHandlerPool[:i], e.UMHandlerPool[i+1:]...)
			return nil
		}
	}
	return lib.ChipErrorNoUnsolicitedMessageHandler
}

//func (e *ExchangeManagerImpl) OnMessageReceived(
//	packetHeader message.PacketHeader,
//	payloadHeader message.PayloadHeader,
//	session transport.SessionHandle,
//	isDuplicate uint8,
//	data []byte,
//) error{
//
//}
