package messageing

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/bitflags"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/platform/system"
	"github.com/galenliu/chip/protocols"
	log "golang.org/x/exp/slog"
	"math/rand"
)

type AnyMessage uint8

const KAnyMessageType AnyMessage = 0xFF

func (a AnyMessage) Name() string {
	return "Any MessageType"
}

func (a AnyMessage) MessageType() uint8 {
	return 0xFF
}

func (a AnyMessage) Matches(m uint8) bool {
	return true
}

type MessageType interface {
	String() string
	MessageType() uint8
	ProtocolId() protocols.Id
}

const (
	kStateNotInitialized = 0
	kStateInitialized    = 1
)

type ExchangeManagerBase interface {
	// SessionMessageDelegate the delegate for transport session manager
	transport.SessionMessageDelegate
	SessionManager() *transport.SessionManager
	RegisterUnsolicitedMessageHandlerForProtocol(protocolId protocols.Id, handler UnsolicitedMessageHandler) error
	RegisterUnsolicitedMessageHandlerForType(protocolId protocols.Id, msgType MessageType, handler UnsolicitedMessageHandler) error
	UnregisterUnsolicitedMessageHandlerForType(id protocols.Id, messageType MessageType) error
	UnregisterUnsolicitedMessageHandlerForProtocol(id protocols.Id) error
	OnResponseTimeout(ec *ExchangeContext)
	OnExchangeClosing(ec *ExchangeContext)
	GetMessageDispatch() ExchangeMessageDispatchBase
	ReleaseContext(ctx *ExchangeContext)
	Shutdown()
}

type UnsolicitedMessageHandlerSlot struct {
	messageType uint8
	handler     UnsolicitedMessageHandler
	protocolId  protocols.Id
}

func (slot *UnsolicitedMessageHandlerSlot) Matches(aProtocolId protocols.Id, aMessageType uint8) bool {
	return aProtocolId == slot.protocolId && aMessageType == slot.messageType
}

func (slot *UnsolicitedMessageHandlerSlot) Reset() {
	slot.handler = nil
}

func (slot *UnsolicitedMessageHandlerSlot) IsInUse() bool {
	return slot.handler != nil
}

type ExchangeManager struct {
	mUMHandlerPool      [config.MaxUnsolicitedMessageHandlers]*UnsolicitedMessageHandlerSlot
	mContextPool        *ExchangeContextPool
	mSessionManager     *transport.SessionManager
	mNextExchangeId     uint16
	mNextKeyId          uint16
	mReliableMessageMgr *ReliableMessageMgr
	mInitialized        bool
}

func NewExchangeManager() *ExchangeManager {
	impl := &ExchangeManager{}
	impl.mUMHandlerPool = [config.MaxUnsolicitedMessageHandlers]*UnsolicitedMessageHandlerSlot{}
	impl.mContextPool = NewExchangeContextContainer()
	return impl
}

func (e *ExchangeManager) Init(sessionManager *transport.SessionManager) error {
	if e.mInitialized {
		return lib.IncorrectState
	}
	e.mSessionManager = sessionManager
	e.mNextExchangeId = uint16(rand.Uint32())
	e.mNextKeyId = 0
	for _, umHandler := range e.mUMHandlerPool {
		if umHandler.handler != nil {
			umHandler.Reset()
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

func (e *ExchangeManager) GetMessageDispatch() ExchangeMessageDispatchBase {
	//TODO implement me
	panic("implement me")
}

func (e *ExchangeManager) SessionManager() *transport.SessionManager {
	return e.mSessionManager
}

func (e *ExchangeManager) OnMessageReceived(
	packetHeader *raw.PacketHeader,
	payloadHeader *raw.PayloadHeader,
	ss *transport.SessionHandle,
	isDuplicate bool,
	msg *system.PacketBufferHandle,
) {

	var matchingUMH *UnsolicitedMessageHandlerSlot = nil

	protocolId := protocols.New(payloadHeader.ProtocolID(), payloadHeader.VendorId())
	{
		//logging
		var compressedFabricId lib.CompressedFabricId = 0
		if ss.IsSecure() && e.mSessionManager.FabricTable() != nil {
			secureSession := ss.Session.(*session.Secure)
			fabricInfo := e.mSessionManager.FabricTable().FindFabricWithIndex(secureSession.FabricIndex())
			if fabricInfo != nil {
				compressedFabricId = fabricInfo.CompressedFabricId()
			}
		}
		log.Debug(">>>",
			"E", fmt.Sprintf("%04X", payloadHeader.ExchangeId()),
			"M", fmt.Sprintf("%08X", packetHeader.MessageCounter),
			"Ack", fmt.Sprintf("%04X", payloadHeader.AckMessageCounter()),
			"S", ss.Type(),
			"From", log.GroupValue(
				log.Any("FabricIndex", ss.FabricIndex()),
				log.Any("NodeId", ss.GetPeer().NodeId()),
				log.Any("FabricId", compressedFabricId),
			),
			"Type", log.GroupValue(
				log.Any("protocolId", payloadHeader.ProtocolID()),
				log.Any("opcode", payloadHeader.MessageType()),
			))
	}

	var msgFlags = bitflags.Some(uint32(0))
	msgFlags.Set(isDuplicate, fDuplicateMessage)

	if !packetHeader.IsGroupSession() {
		ec := e.mContextPool.MatchExchange(ss, packetHeader, payloadHeader)
		if ec != nil {
			log.Info("Found matching",
				"Exchange", ec)
			_ = ec.HandleMessage(packetHeader.MessageCounter, payloadHeader, msgFlags.Value(), msg)
			return
		}
	} else {
		log.Info("Received Group Cast Message", "GroupId", packetHeader.DestinationGroupId)
	}

	if !ss.IsActive() {
		log.Info("Dropping message on inactive session that does not match an existing exchange")
		return
	}

	//如果不是重复的消息，而且如果消息是对方发起
	if msgFlags.Has(fDuplicateMessage) && payloadHeader.IsInitiator() {
		for i, umh := range e.mUMHandlerPool {
			id := protocolId
			if umh.IsInUse() && id.Equal(umh.protocolId) {
				if umh.messageType == payloadHeader.MessageType() {
					matchingUMH = e.mUMHandlerPool[i]
					break
				}
				if umh.messageType == uint8(KAnyMessageType) {
					matchingUMH = umh
				}
			}
		}
	} else if !payloadHeader.NeedsAck() {
		log.Error("ExchangeManager OnMessageReceived failed", lib.MATTER_ERROR_UNSOLICITED_MSG_NO_ORIGINATOR)
		return
	}

	defer e.sendStandaloneAckIfNeeded(packetHeader, payloadHeader, ss, msgFlags.Value(), msg)

	if matchingUMH != nil {
		delegate, err := matchingUMH.handler.OnUnsolicitedMessageReceived(payloadHeader)
		if err != nil {
			log.Error("ExchangeManager OnMessageReceived", err)
			return
		}
		var ec = e.mContextPool.create(e, payloadHeader.ExchangeId(), ss, false, delegate, false)
		if ec == nil {
			if delegate != nil {
				matchingUMH.handler.OnExchangeCreationFailed(delegate)
			}
			log.Error("ExchangeManager OnMessageReceived failed", err)
			return
		}
		log.Info("ExchangeManager Handling", "exchange", ec, "delegate", ec.Delegate())

		if ec.IsEncryptionRequired() != packetHeader.IsEncrypted() {
			log.Info("ExchangeManager OnMessageReceived", lib.MATTER_ERROR_INVALID_MESSAGE_TYPE)
			ec.close()
			return
		}

		err = ec.HandleMessage(packetHeader.MessageCounter, payloadHeader, msgFlags.Value(), msg)
		if err != nil {
			log.Error("ExchangeManager OnMessageReceived failed", err)
		}
		return
	}
}

func (e *ExchangeManager) RegisterUnsolicitedMessageHandlerForProtocol(
	protocolId protocols.Id,
	handler UnsolicitedMessageHandler,
) error {
	return e.registerUMH(protocolId, KAnyMessageType.MessageType(), handler)
}

func (e *ExchangeManager) RegisterUnsolicitedMessageHandlerForType(id protocols.Id, msgType MessageType, handler UnsolicitedMessageHandler) error {
	return e.registerUMH(id, msgType.MessageType(), handler)
}

func (e *ExchangeManager) UnregisterUnsolicitedMessageHandlerForType(id protocols.Id, messageType MessageType) error {
	return e.unregisterUMH(id, messageType.MessageType())
}

func (e *ExchangeManager) UnregisterUnsolicitedMessageHandlerForProtocol(id protocols.Id) error {
	return e.unregisterUMH(id, KAnyMessageType.MessageType())
}

func (e *ExchangeManager) registerUMH(protocolId protocols.Id, msgType uint8, handle UnsolicitedMessageHandler) error {

	var selected *UnsolicitedMessageHandlerSlot
	for _, umh := range e.mUMHandlerPool {
		if !umh.IsInUse() {
			if selected == nil {
				selected = umh
				break
			}
		} else if umh.Matches(protocolId, msgType) {
			umh.handler = handle
			return nil
		}
	}
	if selected == nil {
		return lib.TooManyUnsolicitedMessageHandlers
	}
	selected.handler = handle
	selected.protocolId = protocolId
	selected.messageType = msgType
	return nil
}

func (e *ExchangeManager) unregisterUMH(id protocols.Id, msgType uint8) error {
	for _, umh := range e.mUMHandlerPool {
		if umh.IsInUse() && umh.Matches(id, msgType) {
			umh.Reset()
			return nil
		}
	}
	return lib.NoUnsolicitedMessageHandler
}

func (e *ExchangeManager) sendStandaloneAckIfNeeded(
	packetHeader *raw.PacketHeader,
	payloadHeader *raw.PayloadHeader,
	session *transport.SessionHandle,
	msgFlag uint32,
	buf *system.PacketBufferHandle,
) {
	if !payloadHeader.NeedsAck() {
		return
	}
	ec := e.mContextPool.create(e, payloadHeader.ExchangeId(), session, !payloadHeader.IsInitiator(), nil, true)
	if ec == nil {
		log.Error("ExchangeManager OnMessageReceived failed", lib.MATTER_ERROR_NO_MEMORY)
		return
	}
	log.Debug("ExchangeManager Generating StandaloneAck", "Exchange", ec)
	err := ec.HandleMessage(packetHeader.MessageCounter, payloadHeader, msgFlag, buf)
	if err != nil {
		log.Error("ExchangeManager OnMessageReceived failed", err)
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
	return e.mContextPool.create(e, e.mNextExchangeId, session, isInitiator, delegate, false)
}

func (e *ExchangeManager) ReleaseContext(ctx *ExchangeContext) {
	e.mContextPool.Release(ctx)
}

func (e *ExchangeManager) Shutdown() {
}

func (e *ExchangeManager) ReliableMessageMgr() *ReliableMessageMgr {
	return e.mReliableMessageMgr
}
