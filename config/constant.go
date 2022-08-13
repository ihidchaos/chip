package config

const (
	ChipConfigMaxUnsolicitedMessageHandlers      = 8
	ChipConfigMaxFabrics                    int  = 16
	ChipConfigSecureSessionPoolSize              = ChipConfigMaxFabrics*3 + 2
	ChipConfigMessageCounterWindowSize      uint = 32
)
