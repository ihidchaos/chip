package messageing

import (
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"sync"
)

type ExchangeContextPool struct {
	Container sync.Map
}

func NewExchangeContextContainer() *ExchangeContextPool {
	return &ExchangeContextPool{}
}

func (c *ExchangeContextPool) Add(context *ExchangeContext) {
	c.Container.Store(context.mExchangeId, context)
}

func (c *ExchangeContextPool) Create(ec ExchangeManager,
	exchangeId uint16,
	session transport.SessionHandle,
	initiator bool,
	delegate ExchangeDelegate,
	isEphemeralExchange bool) *ExchangeContext {
	context := NewExchangeContext(ec, exchangeId, session, initiator, delegate, isEphemeralExchange)
	c.Add(context)
	return context
}

func (c *ExchangeContextPool) Delete(id uint16) {
	c.Container.Delete(id)
}

func (c *ExchangeContextPool) Get(id uint16) *ExchangeContext {
	a, ok := c.Container.Load(id)
	if ok {
		context, ok := a.(*ExchangeContext)
		if ok {
			return context
		}
	}
	return nil
}

func (c *ExchangeContextPool) MatchExchange(
	session transport.SessionHandle,
	packetHeader *raw.PacketHeader,
	payloadHeader *raw.PayloadHeader,
) *ExchangeContext {
	ec := c.Get(payloadHeader.GetExchangeID())
	if !ec.MatchExchange(session, packetHeader, payloadHeader) {
		ec = nil
	}
	return ec
}
