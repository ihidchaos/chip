package messageing

import (
	"github.com/galenliu/chip/transport"
	"github.com/galenliu/chip/transport/message"
)

type ExchangeContext struct {
}

func (c ExchangeContext) MatchExchange(session transport.SessionHandle, header message.PacketHeader, header2 message.PayloadHeader) bool {
	return false
}

func (c ExchangeContext) HandleMessage(counter uint32, header message.PayloadHeader, flags uint32, data []byte) {

}
