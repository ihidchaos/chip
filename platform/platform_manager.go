package platform

import "sync"

type EventHandlerFunct func(*ChipDeviceEvent, uint64)

type ManagerDelegate interface {
}

type PlatformManager interface {
}

type PlatformManagerImpl struct {
}

var _instance *PlatformManagerImpl
var once sync.Once

func PlatformMgr() *PlatformManagerImpl {
	once.Do(func() {
		_instance = &PlatformManagerImpl{}
	})
	return _instance
}

func (m *PlatformManagerImpl) AddEventHandler(funct EventHandlerFunct, uint642 uint64) {

}

func (m *PlatformManagerImpl) RunEventLoop() {

}
