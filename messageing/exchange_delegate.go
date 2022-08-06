package messageing

import (
	"github.com/galenliu/chip/messageing/transport/raw"
)

type ExchangeDelegate interface {
	OnMessageReceived(context *ExchangeContext, header raw.PayloadHeader, data []byte) error
	OnResponseTimeout(ec *ExchangeContext)
	OnExchangeClosing(ec *ExchangeContext)
}

type UnsolicitedMessageHandler interface {
	OnUnsolicitedMessageReceived(header raw.PayloadHeader, delegate ExchangeDelegate) error
	OnExchangeCreationFailed(delegate ExchangeDelegate)
}