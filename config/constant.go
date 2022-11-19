package config

const (
	MaxUnsolicitedMessageHandlers      = 8
	MaxFabrics                    int  = 16
	MessageCounterWindowSize      uint = 32

	UnauthenticatedConnectionPoolSize int = 4
	SecureSessionPoolSize                 = MaxFabrics*3 + 2

	MaxExchangeContexts = 16

	RMPRetransTableSize = 16

	Sha256ContextSize = (4 * (8 + 2 + 16 + 2)) + 8
)
