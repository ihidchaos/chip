package messageing

import (
	"github.com/galenliu/chip/pkg"
	"github.com/galenliu/chip/protocols"
	"github.com/galenliu/chip/secure_channel"
	"github.com/galenliu/chip/transport"
	"github.com/galenliu/chip/transport/message"
	log "github.com/sirupsen/logrus"
)

type UnsolicitedMessageHandlerSlot struct {
	MessageType uint8
	Handler     message.UnsolicitedMessageHandler
	ProtocolId  protocols.Id
}

func (receiver UnsolicitedMessageHandlerSlot) Matches(aProtocolId protocols.Id, aMessageType uint8) bool {
	return aProtocolId == receiver.ProtocolId && aMessageType == receiver.MessageType
}

type ExchangeManager interface {
	Init(sessions transport.SessionManager) error
	RegisterUnsolicitedMessageHandlerForType(sigma1 uint8, s *secure_channel.CASEServer, id ...protocols.Id) error
}

type ExchangeManagerImpl struct {
	UMHandlerPool []*UnsolicitedMessageHandlerSlot
	mContextPool  []ExchangeContext
}

func (e *ExchangeManagerImpl) RegisterUnsolicitedMessageHandlerForType(msgType uint8, handler *secure_channel.CASEServer, id ...protocols.Id) error {
	var _id = protocols.NotSpecifiedId
	if len(id) != 0 {
		_id = id[0]
	}
	err := e.RegisterUMH(_id, msgType, handler)
	if err != nil {
		return err
	}
	return nil
}

func NewExchangeManagerImpl() *ExchangeManagerImpl {
	return &ExchangeManagerImpl{}
}

func (e *ExchangeManagerImpl) Init(sessions transport.SessionManager) error {
	return nil
}

func (e *ExchangeManagerImpl) RegisterUMH(id protocols.Id, msgType uint8, handle message.UnsolicitedMessageHandler) error {

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

func (e *ExchangeManagerImpl) UnregisterUMH(id protocols.Id, msgType uint8) error {
	for i, umh := range e.UMHandlerPool {
		if umh.Matches(id, msgType) {
			e.UMHandlerPool = append(e.UMHandlerPool[:i], e.UMHandlerPool[i+1:]...)
			return nil
		}
	}
	return pkg.ChipErrorNoUnsolicitedMessageHandler
}

func (e *ExchangeManagerImpl) OnMessageReceived(
	packetHeader message.PacketHeader,
	payloadHeader message.PayloadHeader,
	session transport.SessionHandle,
	isDuplicate uint8,
	data []byte,
) {
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
