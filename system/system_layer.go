package system

type Layer interface {
	CancelTimer(timeout func(layer Layer, aAppState any), c any)
}
