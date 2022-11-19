package system

import (
	"github.com/galenliu/chip/messageing"
	"time"
)

type Layer interface {
	CancelTimer(timeout func(layer Layer, aAppState any), c any)
	StartTimer(timeout time.Duration, timeout2 func(layer Layer, aAppState any), c *messageing.ExchangeContext)
}
