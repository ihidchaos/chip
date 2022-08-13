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
	MessageType int16
	Handler     UnsolicitedMessageHandler
	ProtocolId  *protocols.Id
}

func (receiver *UnsolicitedMessageHandlerSlot) Matches(aProtocolId *protocols.Id, aMessageType int16) bool {
	return aProtocolId == receiver.ProtocolId && aMessageType == receiver.MessageType
}

func (receiver *UnsolicitedMessageHandlerSlot) Reset() {
	receiver.Handler = nil
}

func (receiver *UnsolicitedMessageHandlerSlot) IsInUse() bool {
	return receiver.Handler != nil
}

// ExchangeManager
// impl transport.SessionMessageDelegate
type ExchangeManager interface {
	// SessionMessageDelegate the delegate for transport session manager
	transport.SessionMessageDelegate
	RegisterUnsolicitedMessageHandlerForProtocol(protocolId *protocols.Id, handler UnsolicitedMessageHandler) error
	RegisterUnsolicitedMessageHandlerForType(protocolId *protocols.Id, msgType uint8, handler UnsolicitedMessageHandler) error
	OnResponseTimeout(ec *ExchangeContext)
	OnExchangeClosing(ec *ExchangeContext)
	GetMessageDispatch() ExchangeMessageDispatch
}

// ExchangeManagerImpl ExchangeManager
type ExchangeManagerImpl struct {
	UMHandlerPool       [config.ChipConfigMaxUnsolicitedMessageHandlers]UnsolicitedMessageHandlerSlot
	mContextPool        *ExchangeContextPool
	mSessionManager     transport.SessionManager
	mNextExchangeId     uint16
	mNextKeyId          uint16
	mReliableMessageMgr any //TODO
	mInitialized        bool
}

func NewExchangeManagerImpl() *ExchangeManagerImpl {
	impl := &ExchangeManagerImpl{}
	impl.UMHandlerPool = [config.ChipConfigMaxUnsolicitedMessageHandlers]UnsolicitedMessageHandlerSlot{}
	impl.mContextPool = NewExchangeContextContainer()
	return impl
}

func (e *ExchangeManagerImpl) Init(sessionManager transport.SessionManager) error {
	if e.mInitialized {
		return lib.ChipErrorIncorrectState
	}
	e.mSessionManager = sessionManager
	e.mNextExchangeId = uint16(rand.Uint32())
	e.mNextKeyId = 0
	e.UMHandlerPool = [config.ChipConfigMaxUnsolicitedMessageHandlers]UnsolicitedMessageHandlerSlot{}

	for _, handler := range e.UMHandlerPool {
		if handler.Handler != nil {
			handler.Reset()
		}
	}
	sessionManager.SetMessageDelegate(e)
	//e.mReliableMessageMgr.init(sessionManager.SystemLayer())
	e.mInitialized = true
	return nil
}

func (e *ExchangeManagerImpl) OnResponseTimeout(ec *ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (e *ExchangeManagerImpl) OnExchangeClosing(ec *ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (e *ExchangeManagerImpl) GetMessageDispatch() ExchangeMessageDispatch {
	//TODO implement me
	panic("implement me")
}

func (e *ExchangeManagerImpl) OnMessageReceived(
	packetHeader *raw.PacketHeader,
	payloadHeader *raw.PayloadHeader,
	session transport.SessionHandle,
	isDuplicate uint8,
	buf *lib.PacketBuffer,
) {
	var matchingUMH *UnsolicitedMessageHandlerSlot = nil
	log.Infof("Received message of type %d with protocolId %d"+
		"and MessageCounter: %d",
		payloadHeader.GetMessageType(), payloadHeader.GetProtocolID(),
		packetHeader.GetMessageCounter())

	var msgFlags uint32 = 0
	msgFlags = lib.SetFlag(isDuplicate == transport.KDuplicateMessageYes, msgFlags, transport.FDuplicateMessage)

	if !packetHeader.IsGroupSession() {
		ec := e.mContextPool.MatchExchange(session, packetHeader, payloadHeader)
		if ec != nil {
			log.Infof("Found matching exchange: %d Delegate: %p",
				ec.mExchangeId, ec.GetDelegate())
			_ = ec.HandleMessage(packetHeader.GetMessageCounter(), payloadHeader, msgFlags, buf)
			return
		}
	}
	log.Infof("Received Groupcast Message with GroupId of %d", packetHeader.GetDestinationGroupId())

	if !session.IsActiveSession() {
		log.Infof("Dropping message on inactive session that does not match an existing exchange")
		return
	}

	if lib.HasFlags(msgFlags, transport.FDuplicateMessage) && payloadHeader.IsInitiator() {
		matchingUMH = nil
		for _, umh := range e.UMHandlerPool {
			if umh.IsInUse() && payloadHeader.HasProtocol(umh.ProtocolId) {
				matchingUMH = &umh
				break
			}
			if umh.MessageType == KAnyMessageType {
				matchingUMH = &umh
			}
		}
	} else if !payloadHeader.NeedsAck() {
		log.Infof("OnMessageReceived failed")

		return
	}
	if matchingUMH != nil {
		var delegate ExchangeDelegate = nil
		err := matchingUMH.Handler.OnUnsolicitedMessageReceived(payloadHeader, delegate)
		if err != nil {
			log.Infof("OnMessageReceived failed, err = %s", err.Error())
			e.SendStandaloneAckIfNeeded(packetHeader, payloadHeader, session, msgFlags, buf)
			return
		}
		var ec = e.mContextPool.Create(e, payloadHeader.GetExchangeID(), session, false, delegate, false)
		if ec == nil {
			if delegate == nil {
				matchingUMH.Handler.OnExchangeCreationFailed(delegate)
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

		err = ec.HandleMessage(packetHeader.GetMessageCounter(), payloadHeader, msgFlags, buf)
		if err != nil {
			log.Infof("OnMessageReceived failed,err : %s", err.Error())
		}
		return
	}
	e.SendStandaloneAckIfNeeded(packetHeader, payloadHeader, session, msgFlags, buf)
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

func (e *ExchangeManagerImpl) RegisterUMH(id *protocols.Id, msgType int16, handle UnsolicitedMessageHandler) error {

	var selected *UnsolicitedMessageHandlerSlot
	for _, umh := range e.UMHandlerPool {
		if !umh.IsInUse() {
			if selected == nil {
				selected = &umh
			}
		} else if umh.Matches(id, msgType) {
			umh.Handler = handle
			return nil
		}
	}
	if selected == nil {
		return lib.ChipErrorTooManyUnsolicitedMessageHandlers
	}
	selected.Handler = handle
	selected.ProtocolId = id
	selected.MessageType = msgType
	return nil
}

func (e *ExchangeManagerImpl) UnregisterUMH(id *protocols.Id, msgType int16) error {
	for _, umh := range e.UMHandlerPool {
		if umh.IsInUse() && umh.Matches(id, msgType) {
			umh.Reset()
			return nil
		}
	}
	return lib.ChipErrorNoUnsolicitedMessageHandler
}

func (e *ExchangeManagerImpl) SendStandaloneAckIfNeeded(
	header *raw.PacketHeader,
	payloadHeader *raw.PayloadHeader,
	session transport.SessionHandle,
	msgFlag uint32,
	buf *lib.PacketBuffer,
) {

}
