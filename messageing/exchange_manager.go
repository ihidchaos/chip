package messageing

import (
	"errors"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/platform/system"
	"github.com/galenliu/chip/protocols"
	log "golang.org/x/exp/slog"
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
	UMHandlerPool       [config.MaxUnsolicitedMessageHandlers]UnsolicitedMessageHandlerSlot
	mContextPool        *ExchangeContextPool
	mSessionManager     transport.SessionManagerBase
	mNextExchangeId     uint16
	mNextKeyId          uint16
	mReliableMessageMgr *ReliableMessageMgr
	mInitialized        bool
}

func NewExchangeManager() *ExchangeManager {
	impl := &ExchangeManager{}
	impl.UMHandlerPool = [config.MaxUnsolicitedMessageHandlers]UnsolicitedMessageHandlerSlot{}
	impl.mContextPool = NewExchangeContextContainer()
	return impl
}

func (e *ExchangeManager) Init(sessionManager transport.SessionManagerBase) error {
	if e.mInitialized {
		return lib.IncorrectState
	}
	e.mSessionManager = sessionManager
	e.mNextExchangeId = uint16(rand.Uint32())
	e.mNextKeyId = 0
	e.UMHandlerPool = [config.MaxUnsolicitedMessageHandlers]UnsolicitedMessageHandlerSlot{}

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
	buf *system.PacketBufferHandle,
) {
	var matchingUMH *UnsolicitedMessageHandlerSlot = nil
	log.Info("Received message",
		"type", payloadHeader.GetMessageType(), "protocolId", payloadHeader.GetProtocolID(),
		"messageCounter", packetHeader.MessageCounter, "Tag", "ExchangeManager")

	var msgFlags uint32 = 0
	msgFlags = lib.SetFlag(isDuplicate == transport.DuplicateMessageYes, msgFlags, transport.DuplicateMessage)

	if !packetHeader.IsGroupSession() {
		ec := e.mContextPool.MatchExchange(session, packetHeader, payloadHeader)
		if ec != nil {
			log.Info("Found matching exchange",
				"exchange", ec.Marshall(), "Tag", "ExchangeManager")
			_ = ec.HandleMessage(packetHeader.MessageCounter, payloadHeader, msgFlags, buf)
			return
		}
	}
	log.Info("Received Groupcast Message with GroupId of %d", packetHeader.DestinationGroupId)

	if !session.IsActiveSession() {
		log.Info("Dropping message on inactive session that does not match an existing exchange")
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
		log.Info("OnMessageReceived failed", "Tag", "ExchangeManager")
		return
	}
	if matchingUMH != nil {
		var delegate ExchangeDelegate = nil
		err := matchingUMH.handler.OnUnsolicitedMessageReceived(payloadHeader, delegate)
		if err != nil {
			log.Error("OnMessageReceived", err, "Tag", "ExchangeManager")
			e.SendStandaloneAckIfNeeded(packetHeader, payloadHeader, session, msgFlags, buf)
			return
		}
		var ec = e.mContextPool.Create(e, payloadHeader.GetExchangeID(), session, false, delegate, false)
		if ec == nil {
			if delegate == nil {
				matchingUMH.handler.OnExchangeCreationFailed(delegate)
			}
			log.Error("OnMessageReceived failed", err, "Tag", "ExchangeManager")
			return
		}
		log.Info("Handling", "exchange", ec.Marshall(), "delegate", ec.GetDelegate(), "Tag", "ExchangeManager")

		if ec.IsEncryptionRequired() != packetHeader.IsEncrypted() {
			log.Info("OnMessageReceived", errors.New("invalid message type"), "Tag", "ExchangeManager")
			ec.Close()
			e.SendStandaloneAckIfNeeded(packetHeader, payloadHeader, session, msgFlags, buf)
		}

		err = ec.HandleMessage(packetHeader.MessageCounter, payloadHeader, msgFlags, buf)
		if err != nil {
			log.Info("OnMessageReceived failed", err, "Tag", "ExchangeManager")
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
		return lib.TooManyUnsolicitedMessageHandlers
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
	return lib.NoUnsolicitedMessageHandler
}

func (e *ExchangeManager) SendStandaloneAckIfNeeded(
	packetHeader *raw.PacketHeader,
	payloadHeader *raw.PayloadHeader,
	session *transport.SessionHandle,
	msgFlag uint32,
	buf *system.PacketBufferHandle,
) {
	if !payloadHeader.NeedsAck() {
		return
	}
	ec := e.mContextPool.Create(e, payloadHeader.GetExchangeID(), session, !payloadHeader.IsInitiator(), nil, true)
	if ec == nil {
		log.Error("ExchangeManager OnMessageReceived ", errors.New("no memory"))
		return
	}
	log.Info("ExchangeManager Generating StandaloneAck", "ExchangeContext", ec.Marshall(), "Tag", "ExchangeManager")
	err := ec.HandleMessage(packetHeader.GetMessageCounter(), payloadHeader, msgFlag, buf)
	if err != nil {
		log.Error("OnMessageReceived", err, "Tag", "ExchangeManager")
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
