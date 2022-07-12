package messageing

import (
	"sync"
	"time"
)

var ChipConfigMrpDefaultIdleRetryInterval int64 = 5000
var ChipConfigMrpDefaultActiveRetryInterval int64 = 300

type ReliableMessageProtocolConfig struct {
	IdleRetransTimeout   time.Duration
	ActiveRetransTimeout time.Duration
}

func (c ReliableMessageProtocolConfig) Init() *ReliableMessageProtocolConfig {
	c.ActiveRetransTimeout = time.Duration(0)
	c.ActiveRetransTimeout = time.Duration(0)
	return &c
}

var insRMPC *ReliableMessageProtocolConfig
var rmpcOnce = sync.Once{}

func GetLocalMRPConfig() *ReliableMessageProtocolConfig {
	rmpcOnce.Do(func() {
		insRMPC = newReliableMessageProtocolConfig()
		insRMPC.IdleRetransTimeout = time.Duration(ChipConfigMrpDefaultIdleRetryInterval)
		insRMPC.ActiveRetransTimeout = time.Duration(ChipConfigMrpDefaultActiveRetryInterval)
	})
	return insRMPC
}

func newReliableMessageProtocolConfig() *ReliableMessageProtocolConfig {
	return &ReliableMessageProtocolConfig{}
}
