package platform

import "sync"

type EventHandlerFunct func(*ChipDeviceEvent, uint64)

type ManagerDelegate interface {
}

type PlatformManager struct {
}

var _instance *PlatformManager
var once sync.Once

func PlatformMgr() *PlatformManager {
	once.Do(func() {
		_instance = &PlatformManager{}
	})
	return _instance
}

func (m PlatformManager) AddEventHandler(funct EventHandlerFunct, uint642 uint64) {

}

func (m PlatformManager) RunEventLoop() {

}
