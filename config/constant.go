package config

const (
	MaxUnsolicitedMessageHandlers      = 8
	MaxFabrics                    int  = 16
	SecureSessionPoolSize              = MaxFabrics*3 + 2
	MessageCounterWindowSize      uint = 32
)
