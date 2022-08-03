package messageing

import "github.com/galenliu/chip/transport/message"

type ExchangeDelegate interface {
	OnMessageReceived(context *ExchangeContext, header message.PayloadHeader, data []byte) error
	OnResponseTimeout(ec *ExchangeContext)
	OnExchangeClosing(ec *ExchangeContext)
}

type UnsolicitedMessageHandler interface {
	OnUnsolicitedMessageReceived(header message.PayloadHeader, delegate ExchangeDelegate) error
	OnExchangeCreationFailed(delegate ExchangeDelegate)
}
