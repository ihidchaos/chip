package session

import (
	"sync"
	"time"
)

var ConfigMrpDefaultIdleRetryInterval int64 = 5000
var ConfigMrpDefaultActiveRetryInterval int64 = 300

type ReliableMessageProtocolConfig struct {
	IdleRetransTimeout   time.Duration
	ActiveRetransTimeout time.Duration
}

func DefaultMRPConfig() *ReliableMessageProtocolConfig {
	return &ReliableMessageProtocolConfig{
		IdleRetransTimeout:   300 * time.Millisecond,
		ActiveRetransTimeout: 300 * time.Millisecond,
	}
}

func (c ReliableMessageProtocolConfig) Init() *ReliableMessageProtocolConfig {
	c.ActiveRetransTimeout = time.Duration(0)
	c.ActiveRetransTimeout = time.Duration(0)
	return &c
}

var _rmpc *ReliableMessageProtocolConfig
var rmpcOnce = sync.Once{}

func GetLocalMRPConfig() *ReliableMessageProtocolConfig {
	rmpcOnce.Do(func() {
		_rmpc = newReliableMessageProtocolConfig()

	})
	return _rmpc
}

func newReliableMessageProtocolConfig() *ReliableMessageProtocolConfig {
	rmpc := &ReliableMessageProtocolConfig{}
	rmpc.IdleRetransTimeout = time.Duration(ConfigMrpDefaultIdleRetryInterval)
	rmpc.ActiveRetransTimeout = time.Duration(ConfigMrpDefaultActiveRetryInterval)
	return rmpc
}

func GetRetransmissionTimeout(activeRetransTimeout time.Duration, idleRetransTimeout time.Duration, lastActivityTime time.Time, activeTime time.Duration) time.Duration {
	return time.Duration(0)
}
