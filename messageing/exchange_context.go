package messageing

import (
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
)

type ExchangeSessionHolder interface {
	transport.SessionHolderWithDelegate
	transport.SessionHolder
}

type ExchangeContext struct {
	ReliableMessageContext
	mExchangeId uint16
	mDispatch   ExchangeMessageDispatch
	mSession    ExchangeSessionHolder
	mDelegate   ExchangeDelegate
}

func NewExchangeContext(
	ec ExchangeManager,
	exchangeId uint16,
	session transport.SessionHandle,
	initiator bool,
	delegate ExchangeDelegate,
	isEphemeralExchange bool,
) *ExchangeContext {
	return &ExchangeContext{}
}

func (c *ExchangeContext) MatchExchange(session transport.SessionHandle, packetHeader *raw.PacketHeader, payloadHeader *raw.PayloadHeader) bool {
	return (c.mExchangeId == payloadHeader.GetExchangeID()) &&
		(c.mSession.Contains(session)) &&
		(c.IsEncryptionRequired() == packetHeader.IsEncrypted()) &&
		(payloadHeader.IsInitiator() != c.IsInitiator())
}

func (c *ExchangeContext) HandleMessage(counter uint32, header *raw.PayloadHeader, flags uint32, buf *raw.PacketBuffer) error {
	return nil
}

func (c *ExchangeContext) IsEncryptionRequired() bool {
	return c.mDispatch.IsEncryptionRequired()
}

func (c *ExchangeContext) SetDelegate(delegate ExchangeDelegate) {

}

func (c *ExchangeContext) IsInitiator() bool {
	return c.mFlags&kFlagInitiator != 0
}

func (c *ExchangeContext) GetDelegate() ExchangeDelegate {
	return c.mDelegate
}

func (c *ExchangeContext) Close() {

}
