package messageing

import (
	"fmt"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"sync"
)

type ExchangeContextPool struct {
	container sync.Map
}

func NewExchangeContextContainer() *ExchangeContextPool {
	return &ExchangeContextPool{
		container: sync.Map{},
	}
}

func (c *ExchangeContextPool) add(context *ExchangeContext) {
	c.container.Store(context.mExchangeId, context)
}

func (c *ExchangeContextPool) create(ec *ExchangeManager,
	exchangeId uint16,
	session *transport.SessionHandle,
	initiator bool,
	delegate ExchangeDelegate,
	isEphemeralExchange bool) *ExchangeContext {
	context := NewExchangeContext(ec, exchangeId, session, initiator, delegate, isEphemeralExchange)
	c.add(context)
	return context
}

func (c *ExchangeContextPool) Get(id uint16) *ExchangeContext {
	a, ok := c.container.Load(id)
	if ok {
		context, ok := a.(*ExchangeContext)
		if ok {
			return context
		}
	}
	return nil
}

func (c *ExchangeContextPool) MatchExchange(
	session *transport.SessionHandle,
	packetHeader *raw.PacketHeader,
	payloadHeader *raw.PayloadHeader,
) *ExchangeContext {
	ec := c.Get(payloadHeader.ExchangeId())
	if ec != nil {
		if !ec.MatchExchange(session, packetHeader, payloadHeader) {
			ec = nil
		}
	}
	return ec
}

func (c *ExchangeContextPool) Allocated() int {
	l := 0
	c.container.Range(func(k, v interface{}) bool {
		l++
		fmt.Println(k, v)
		return true
	})
	return l
}

func (c *ExchangeContextPool) Release(ctx *ExchangeContext) {
	c.container.Delete(ctx.mExchangeId)
}

func (c *ExchangeContextPool) CloseContextForDelegate(delegate ExchangeDelegate) {
	c.container.Range(func(key, value any) bool {
		ec, ok := value.(*ExchangeContext)
		if ok {
			if ec.mDelegate == delegate {
				ec.SetDelegate(nil)
				ec.close()
			}
		}
		return false
	})
}
