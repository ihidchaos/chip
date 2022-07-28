package device

import "sync"

type EventHandlerFunct func(*ChipDeviceEvent, uint64)

type ManagerDelegate interface {
}

type PlatformManager interface {
}

type PlatformManagerImpl struct {
}

var __instance *PlatformManagerImpl
var once sync.Once

func PlatformMgr() *PlatformManagerImpl {
	once.Do(func() {
		__instance = &PlatformManagerImpl{}
	})
	return __instance
}

func (m *PlatformManagerImpl) AddEventHandler(funct EventHandlerFunct, uint642 uint64) {

}

func (m *PlatformManagerImpl) RunEventLoop() {

}
