package messageing

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
)

type UnsolicitedMessageHandler interface {
	OnUnsolicitedMessageReceived(header *raw.PayloadHeader, delegate ExchangeDelegate) error
	OnExchangeCreationFailed(delegate ExchangeDelegate)
}

type ExchangeDelegate interface {
	OnMessageReceived(context *ExchangeContext, header *raw.PayloadHeader, data *lib.PacketBuffer) error
	OnResponseTimeout(ec *ExchangeContext)
	OnExchangeClosing(ec *ExchangeContext)
}
