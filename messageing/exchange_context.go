package messageing

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/buffer"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/protocols"
)

type ExchangeSessionHolder interface {
	transport.SessionHolderWithDelegate
	transport.SessionHolder
}

type ExchangeSessionHolderImpl struct {
	*transport.SessionHolderWithDelegateImpl
}

func NewExchangeSessionHolderImpl(delegate *ExchangeContext) *ExchangeSessionHolderImpl {
	return &ExchangeSessionHolderImpl{
		SessionHolderWithDelegateImpl: transport.NewSessionHolderWithDelegateImpl(delegate),
	}
}

type ExchangeContext struct {
	ReliableMessageContext
	mExchangeId  uint16
	mExchangeMgr ExchangeManager
	mDispatch    ExchangeMessageDispatch
	mSession     *ExchangeSessionHolderImpl
	mDelegate    ExchangeDelegate
	mFlags       uint16
}

func NewExchangeContext(
	em ExchangeManager,
	exchangeId uint16,
	session transport.SessionHandle,
	initiator bool,
	delegate ExchangeDelegate,
	isEphemeralExchange bool,
) *ExchangeContext {
	var flags uint16 = 0
	flags = lib.SetFlag(initiator, flags, kFlagInitiator)
	flags = lib.SetFlag(isEphemeralExchange, flags, kFlagEphemeralExchange)
	ec := &ExchangeContext{
		ReliableMessageContext: ReliableMessageContext{},
		mExchangeId:            exchangeId,
		mExchangeMgr:           em,
		mDispatch:              GetMessageDispatch(isEphemeralExchange, delegate),
		mDelegate:              delegate,
		mFlags:                 flags,
	}
	ec.mSession = NewExchangeSessionHolderImpl(ec)
	ec.mSession.SessionHolderWithDelegateImpl.Grad(session)

	return ec
}

func (c *ExchangeContext) MatchExchange(session transport.SessionHandle, packetHeader *raw.PacketHeader, payloadHeader *raw.PayloadHeader) bool {
	return (c.mExchangeId == payloadHeader.GetExchangeID()) &&
		(c.mSession.Contains(session)) &&
		(c.IsEncryptionRequired() == packetHeader.IsEncrypted()) &&
		(payloadHeader.IsInitiator() != c.IsInitiator())
}

func (c *ExchangeContext) HandleMessage(counter uint32, payloadHeader *raw.PayloadHeader, flags uint32, buf *buffer.PacketBuffer) error {

	//isStandaloneAck := payloadHeader.HasMessageType(uint8(StandaloneAck))
	//isDuplicate := lib.HasFlags(flags, transport.FDuplicateMessage)
	if c.mDelegate != nil {
		err := c.mDelegate.OnMessageReceived(c, payloadHeader, buf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ExchangeContext) IsEncryptionRequired() bool {
	return c.mDispatch.IsEncryptionRequired()
}

func (c *ExchangeContext) SetDelegate(delegate ExchangeDelegate) {

}

func GetMessageDispatch(isEphemeralExchange bool, delegate ExchangeDelegate) ExchangeMessageDispatch {
	if isEphemeralExchange {
		return EphemeralExchangeDispatchImpl{delegate: delegate}
	}
	return ExchangeMessageDispatchImpl{delegate: delegate}
}

func (c *ExchangeContext) IsInitiator() bool {
	return c.mFlags&kFlagInitiator != 0
}

func (c *ExchangeContext) GetDelegate() ExchangeDelegate {
	return c.mDelegate
}

func (c *ExchangeContext) Close() {

}

func (c *ExchangeContext) SendMessage(protocolId protocols.Id, msgType MsgType, r2 []byte, response uint16) error {

	//isStandaloneAck := protocolId == protocols.StandardSecureChannelProtocolId && msgType == StandaloneAck

	return nil
}

func (c *ExchangeContext) sendMessage(id protocols.Id) error {

	return nil
}
