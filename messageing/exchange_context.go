package messageing

import (
	"github.com/galenliu/chip/messageing/transport"
	message2 "github.com/galenliu/chip/messageing/transport/raw"
)

type ExchangeContext struct {
}

func (c ExchangeContext) MatchExchange(session transport.SessionHandle, header *message2.PacketHeader, header2 *message2.PayloadHeader) bool {
	return false
}

func (c ExchangeContext) HandleMessage(counter uint32, header *message2.PayloadHeader, flags uint32, data []byte) {

}
