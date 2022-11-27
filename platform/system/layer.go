package system

import (
	"time"
)

type Layer interface {
	CancelTimer(timeout func(layer Layer, aAppState any), c any)
	//StartTimer(timeout time.Duration, timeout2 func(layer Layer, aAppState any), c *messageing.ExchangeContext)
	StartTimer(timeout time.Duration, timeout2 func(layer Layer, aAppState any), c any)
}
