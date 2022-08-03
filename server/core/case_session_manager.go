package core

type CASESessionManagerConfig struct {
	SessionInitParams DeviceProxyInitParams
	DevicePool        OperationalDeviceProxyPoolDelegate
}

type CASESessionManager struct {
}

func NewCASESessionManager() *CASESessionManager {
	return &CASESessionManager{}
}

func (c *CASESessionManager) Init(layer Layer, config *CASESessionManagerConfig) error {
	return nil
}
