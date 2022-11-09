package messageing

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
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
	messageType int16
	handler     UnsolicitedMessageHandler
	protocolId  *protocols.Id
}

func (slot *UnsolicitedMessageHandlerSlot) Matches(aProtocolId *protocols.Id, aMessageType int16) bool {
	return aProtocolId == slot.protocolId && aMessageType == slot.messageType
}

func (slot *UnsolicitedMessageHandlerSlot) Reset() {
	slot.handler = nil
}

func (slot *UnsolicitedMessageHandlerSlot) IsInUse() bool {
	return slot.handler != nil
}

// ExchangeManagerBase
// impl transport.SessionMessageDelegate
type ExchangeManagerBase interface {
	// SessionMessageDelegate the delegate for transport session manager
	transport.SessionMessageDelegate
	SessionManager() transport.SessionManagerBase
	RegisterUnsolicitedMessageHandlerForProtocol(protocolId *protocols.Id, handler UnsolicitedMessageHandler) error
	RegisterUnsolicitedMessageHandlerForType(protocolId *protocols.Id, msgType uint8, handler UnsolicitedMessageHandler) error
	UnregisterUnsolicitedMessageHandlerForType(id *protocols.Id, messageType uint8) error
	UnregisterUnsolicitedMessageHandlerForProtocol(id *protocols.Id) error
	OnResponseTimeout(ec *ExchangeContext)
	OnExchangeClosing(ec *ExchangeContext)
	GetMessageDispatch() ExchangeMessageDispatch
	ReleaseContext(ctx *ExchangeContext)
	Shutdown()
}

// ExchangeManager ExchangeManager
type ExchangeManager struct {
	UMHandlerPool       [config.ChipConfigMaxUnsolicitedMessageHandlers]UnsolicitedMessageHandlerSlot
	mContextPool        *ExchangeContextPool
	mSessionManager     transport.SessionManagerBase
	mNextExchangeId     uint16
	mNextKeyId          uint16
	mReliableMessageMgr *ReliableMessageMgr
	mInitialized        bool
}

func NewExchangeManager() *ExchangeManager {
	impl := &ExchangeManager{}
	impl.UMHandlerPool = [config.ChipConfigMaxUnsolicitedMessageHandlers]UnsolicitedMessageHandlerSlot{}
	impl.mContextPool = NewExchangeContextContainer()
	return impl
}

func (e *ExchangeManager) Init(sessionManager transport.SessionManagerBase) error {
	if e.mInitialized {
		return lib.MatterErrorIncorrectState
	}
	e.mSessionManager = sessionManager
	e.mNextExchangeId = uint16(rand.Uint32())
	e.mNextKeyId = 0
	e.UMHandlerPool = [config.ChipConfigMaxUnsolicitedMessageHandlers]UnsolicitedMessageHandlerSlot{}

	for _, handler := range e.UMHandlerPool {
		if handler.handler != nil {
			handler.Reset()
		}
	}
	sessionManager.SetMessageDelegate(e)
	//e.mReliableMessageMgr.init(sessionManager.SystemLayer())
	e.mInitialized = true
	return nil
}

func (e *ExchangeManager) OnResponseTimeout(ec *ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (e *ExchangeManager) OnExchangeClosing(ec *ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (e *ExchangeManager) GetMessageDispatch() ExchangeMessageDispatch {
	//TODO implement me
	panic("implement me")
}

func (e *ExchangeManager) SessionManager() transport.SessionManagerBase {
	return e.mSessionManager
}

func (e *ExchangeManager) OnMessageReceived(
	packetHeader *raw.PacketHeader,
	payloadHeader *raw.PayloadHeader,
	session *transport.SessionHandle,
	isDuplicate uint8,
	buf *raw.PacketBuffer,
) {
	var matchingUMH *UnsolicitedMessageHandlerSlot = nil
	log.Infof("Received message of type %d with protocolId %d"+
		"and MessageCounter: %d",
		payloadHeader.GetMessageType(), payloadHeader.GetProtocolID(),
		packetHeader.MessageCounter)

	var msgFlags uint32 = 0
	msgFlags = lib.SetFlag(isDuplicate == transport.DuplicateMessageYes, msgFlags, transport.DuplicateMessage)

	if !packetHeader.IsGroupSession() {
		ec := e.mContextPool.MatchExchange(session, packetHeader, payloadHeader)
		if ec != nil {
			log.Infof("Found matching exchange: %d Delegate: %p",
				ec.mExchangeId, ec.GetDelegate())
			_ = ec.HandleMessage(packetHeader.MessageCounter, payloadHeader, msgFlags, buf)
			return
		}
	}
	log.Infof("Received Groupcast Message with GroupId of %d", packetHeader.DestinationGroupId)

	if !session.IsActiveSession() {
		log.Infof("Dropping message on inactive session that does not match an existing exchange")
		return
	}

	//如果不是重复的消息，而且如果消息是对方发起
	if !lib.HasFlags(msgFlags, transport.DuplicateMessage) && payloadHeader.IsInitiator() {
		matchingUMH = nil

		for _, umh := range e.UMHandlerPool {
			if umh.IsInUse() && payloadHeader.HasProtocol(umh.protocolId) {
				matchingUMH = &umh
				break
			}
			if umh.messageType == KAnyMessageType {
				matchingUMH = &umh
			}
		}
	} else if !payloadHeader.NeedsAck() {
		log.Infof("OnMessageReceived failed")
		return
	}
	if matchingUMH != nil {
		var delegate ExchangeDelegate = nil
		err := matchingUMH.handler.OnUnsolicitedMessageReceived(payloadHeader, delegate)
		if err != nil {
			log.Infof("OnMessageReceived failed, err = %s", err.Error())
			e.SendStandaloneAckIfNeeded(packetHeader, payloadHeader, session, msgFlags, buf)
			return
		}
		var ec = e.mContextPool.Create(e, payloadHeader.GetExchangeID(), session, false, delegate, false)
		if ec == nil {
			if delegate == nil {
				matchingUMH.handler.OnExchangeCreationFailed(delegate)
			}
			log.Infof("OnMessageReceived failed, err = %s", err.Error())
			return
		}
		log.Infof("Handling via exchange: %d Delegate: %p", ec.mExchangeId, ec.GetDelegate())

		if ec.IsEncryptionRequired() != packetHeader.IsEncrypted() {
			log.Infof("OnMessageReceived failed,err:= ERROR_INVALID_MESSAGE_TYPE ")
			ec.Close()
			e.SendStandaloneAckIfNeeded(packetHeader, payloadHeader, session, msgFlags, buf)
		}

		err = ec.HandleMessage(packetHeader.MessageCounter, payloadHeader, msgFlags, buf)
		if err != nil {
			log.Infof("OnMessageReceived failed,err : %s", err.Error())
		}
		return
	}
	e.SendStandaloneAckIfNeeded(packetHeader, payloadHeader, session, msgFlags, buf)
}

func (e *ExchangeManager) RegisterUnsolicitedMessageHandlerForProtocol(
	protocolId *protocols.Id,
	handler UnsolicitedMessageHandler,
) error {
	return e.registerUMH(protocolId, KAnyMessageType, handler)
}

func (e *ExchangeManager) RegisterUnsolicitedMessageHandlerForType(
	protocolId *protocols.Id,
	msgType uint8,
	handler UnsolicitedMessageHandler,
) error {
	return e.registerUMH(protocolId, int16(msgType), handler)
}

func (e *ExchangeManager) UnregisterUnsolicitedMessageHandlerForType(id *protocols.Id, messageType uint8) error {
	return e.unregisterUMH(id, int16(messageType))
}

func (e *ExchangeManager) UnregisterUnsolicitedMessageHandlerForProtocol(id *protocols.Id) error {
	return e.unregisterUMH(id, KAnyMessageType)
}

func (e *ExchangeManager) registerUMH(id *protocols.Id, msgType int16, handle UnsolicitedMessageHandler) error {

	var selected *UnsolicitedMessageHandlerSlot
	for _, umh := range e.UMHandlerPool {
		if !umh.IsInUse() {
			if selected == nil {
				selected = &umh
				break
			}
		} else if umh.Matches(id, msgType) {
			umh.handler = handle
			return nil
		}
	}
	if selected == nil {
		return lib.MatterErrorTooManyUnsolicitedMessageHandlers
	}
	selected.handler = handle
	selected.protocolId = id
	selected.messageType = msgType
	return nil
}

func (e *ExchangeManager) unregisterUMH(id *protocols.Id, msgType int16) error {
	for _, umh := range e.UMHandlerPool {
		if umh.IsInUse() && umh.Matches(id, msgType) {
			umh.Reset()
			return nil
		}
	}
	return lib.MatterErrorNoUnsolicitedMessageHandler
}

func (e *ExchangeManager) SendStandaloneAckIfNeeded(
	packetHeader *raw.PacketHeader,
	payloadHeader *raw.PayloadHeader,
	session *transport.SessionHandle,
	msgFlag uint32,
	buf *raw.PacketBuffer,
) {
	if !payloadHeader.NeedsAck() {
		return
	}
	ec := e.mContextPool.Create(e, payloadHeader.GetExchangeID(), session, !payloadHeader.IsInitiator(), nil, true)
	if ec == nil {
		log.Errorf("ExchangeManager OnMessageReceived faild,error= %s", "no memory")
		return
	}
	log.Infof("ExchangeManager Generating StandaloneAck via exchange: %s", ec.Marshall())
	err := ec.HandleMessage(packetHeader.GetMessageCounter(), payloadHeader, msgFlag, buf)
	if err != nil {
		log.Errorf("ExchangeManager OnMessageReceived faild,error= %s", err.Error())
	}
}

func (e *ExchangeManager) CloseAllContextForDelegate(delegate ExchangeDelegate) {
	e.mContextPool.CloseContextForDelegate(delegate)
}

func (e *ExchangeManager) GetNumActiveExchanges() int {
	return e.mContextPool.Allocated()
}

func (e *ExchangeManager) NewContext(session *transport.SessionHandle, delegate ExchangeDelegate, isInitiator bool) *ExchangeContext {
	e.mNextExchangeId = e.mNextExchangeId + 1
	return e.mContextPool.Create(e, e.mNextExchangeId, session, isInitiator, delegate, false)
}

func (e *ExchangeManager) ReleaseContext(ctx *ExchangeContext) {
	e.mContextPool.Release(ctx)
}

func (e *ExchangeManager) Shutdown() {
}

func (e *ExchangeManager) GetReliableMessageMgr() *ReliableMessageMgr {
	return e.mReliableMessageMgr
}

func (e *ExchangeManager) GetSessionManager() transport.SessionManagerBase {
	return e.mSessionManager
}
