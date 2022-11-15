package messageing

import (
	"errors"
	"fmt"
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

type ExchangeManagerBase interface {
	// SessionMessageDelegate the delegate for transport session manager
	transport.SessionMessageDelegate
	SessionManager() *transport.SessionManager
	RegisterUnsolicitedMessageHandlerForProtocol(protocolId *protocols.Id, handler UnsolicitedMessageHandler) error
	RegisterUnsolicitedMessageHandlerForType(protocolId *protocols.Id, msgType uint8, handler UnsolicitedMessageHandler) error
	UnregisterUnsolicitedMessageHandlerForType(id *protocols.Id, messageType uint8) error
	UnregisterUnsolicitedMessageHandlerForProtocol(id *protocols.Id) error
	OnResponseTimeout(ec *ExchangeContext)
	OnExchangeClosing(ec *ExchangeContext)
	GetMessageDispatch() ExchangeMessageDispatchBase
	ReleaseContext(ctx *ExchangeContext)
	Shutdown()
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
	session *transport.SessionHandle,
	isDuplicate bool,
	msg *system.PacketBufferHandle,
) {

	var matchingUMH *UnsolicitedMessageHandlerSlot = nil

	{
		//logging
		var compressedFabricId lib.CompressedFabricId = 0
		if session.IsSecureSession() && e.mSessionManager.FabricTable() != nil {
			secureSession := session.Session.(*transport.SecureSession)
			fabricInfo := e.mSessionManager.FabricTable().FindFabricWithIndex(secureSession.FabricIndex())
			if fabricInfo != nil {
				compressedFabricId = fabricInfo.CompressedFabricId()
			}
		}
		log.Debug(">>>",
			"E", fmt.Sprintf("%04X", payloadHeader.ExchangeId()),
			"M", fmt.Sprintf("%08X", packetHeader.MessageCounter),
			"Ack", fmt.Sprintf("%04X", payloadHeader.AckMessageCounter()),
			"S", session.SessionType(),
			"From", log.GroupValue(
				log.Any("FabricIndex", session.FabricIndex()),
				log.Any("NodeId", session.GetPeer().NodeId()),
				log.Any("FabricId", compressedFabricId),
			),
			"Type", log.GroupValue(
				log.Any("protocolId", payloadHeader.ProtocolID()),
				log.Any("opcode", payloadHeader.MessageType()),
				log.Any("protocolName", payloadHeader.ProtocolID().ProtocolName()),
				log.Any("message type", payloadHeader.ProtocolID().MessageTypeName(payloadHeader.MessageType())),
			))
	}

	var msgFlags uint32 = 0
	msgFlags = lib.SetFlag(isDuplicate, msgFlags, fDuplicateMessage)

	if !packetHeader.IsGroupSession() {
		ec := e.mContextPool.MatchExchange(session, packetHeader, payloadHeader)
		if ec != nil {
			log.Info("Found matching",
				"Exchange", ec)
			_ = ec.HandleMessage(packetHeader.MessageCounter, payloadHeader, msgFlags, msg)
			return
		}
	} else {
		log.Info("Received Group Cast Message", "GroupId", packetHeader.DestinationGroupId)
	}

	if !session.IsActiveSession() {
		log.Info("Dropping message on inactive session that does not match an existing exchange")
		return
	}

	//如果不是重复的消息，而且如果消息是对方发起
	if !lib.HasFlags(msgFlags, fDuplicateMessage) && payloadHeader.IsInitiator() {
		matchingUMH = nil
		for _, umh := range e.mUMHandlerPool {
			if umh.IsInUse() && payloadHeader.ProtocolID().Equal(umh.protocolId) {
				if umh.messageType == int16(payloadHeader.MessageType()) {
					matchingUMH = umh
					break
				}
				if umh.messageType == KAnyMessageType {
					matchingUMH = umh
				}
			}
		}
	} else if !payloadHeader.NeedsAck() {
		log.Error("ExchangeManager OnMessageReceived failed", lib.MATTER_ERROR_UNSOLICITED_MSG_NO_ORIGINATOR)
		return
	}

	if matchingUMH != nil {

		var delegate ExchangeDelegate = nil
		err := matchingUMH.handler.OnUnsolicitedMessageReceived(payloadHeader, delegate)
		if err != nil {
			log.Error("ExchangeManager OnMessageReceived", err)
			e.sendStandaloneAckIfNeeded(packetHeader, payloadHeader, session, msgFlags, msg)
			return
		}
		var ec = e.mContextPool.Create(e, payloadHeader.ExchangeId(), session, false, delegate, false)
		if ec == nil {
			if delegate == nil {
				matchingUMH.handler.OnExchangeCreationFailed(delegate)
			}
			log.Error("ExchangeManager OnMessageReceived failed", err)
			return
		}
		log.Info("ExchangeManager Handling", "exchange", ec, "delegate", ec.Delegate(), "Tag", "ExchangeManager")

		if ec.IsEncryptionRequired() != packetHeader.IsEncrypted() {
			log.Info("OnMessageReceived", errors.New("invalid message type"), "Tag", "ExchangeManager")
			ec.Close()
			e.sendStandaloneAckIfNeeded(packetHeader, payloadHeader, session, msgFlags, msg)
		}

		err = ec.HandleMessage(packetHeader.MessageCounter, payloadHeader, msgFlags, msg)
		if err != nil {
			log.Error("ExchangeManager OnMessageReceived failed", err)
		}
		return
	}
	e.sendStandaloneAckIfNeeded(packetHeader, payloadHeader, session, msgFlags, msg)
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

func (e *ExchangeManager) registerUMH(protocolId *protocols.Id, msgType int16, handle UnsolicitedMessageHandler) error {

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

func (e *ExchangeManager) unregisterUMH(id *protocols.Id, msgType int16) error {
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
	ec := e.mContextPool.Create(e, payloadHeader.ExchangeId(), session, !payloadHeader.IsInitiator(), nil, true)
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
	return e.mContextPool.Create(e, e.mNextExchangeId, session, isInitiator, delegate, false)
}

func (e *ExchangeManager) ReleaseContext(ctx *ExchangeContext) {
	e.mContextPool.Release(ctx)
}

func (e *ExchangeManager) Shutdown() {
}

func (e *ExchangeManager) ReliableMessageMgr() *ReliableMessageMgr {
	return e.mReliableMessageMgr
}
