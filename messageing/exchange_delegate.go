package messageing

import (
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/platform/system"
)

type UnsolicitedMessageHandler interface {
	OnUnsolicitedMessageReceived(header *raw.PayloadHeader) (ExchangeDelegate, error)
	OnExchangeCreationFailed(delegate ExchangeDelegate)
}

type ExchangeDelegate interface {
	OnMessageReceived(context *ExchangeContext, header *raw.PayloadHeader, data *system.PacketBufferHandle) error
	OnResponseTimeout(ec *ExchangeContext)
	OnExchangeClosing(ec *ExchangeContext)
	GetMessageDispatch() ExchangeMessageDispatchBase
}
