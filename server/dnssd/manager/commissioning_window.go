package manager

import (
	"github.com/galenliu/chip/server/dnssd/costants/commissioning_mode"
)

type AppDelegate interface {
	OnCommissioningSessionStarted()
	OnCommissioningSessionStopped()
	OnCommissioningWindowOpened()
	OnCommissioningWindowClosed()
}

type ServerDelegate interface {
}

type CommissioningWindowManager struct {
	mServer      ServerDelegate
	mAppDelegate AppDelegate
}

func (m CommissioningWindowManager) Init(s ServerDelegate) *CommissioningWindowManager {
	m.mServer = s
	return &m
}

func (m *CommissioningWindowManager) SetAppDelegate(delegate AppDelegate) {
	m.mAppDelegate = delegate
}

func (m *CommissioningWindowManager) OpenBasicCommissioningWindow() error {
	return nil
}

func (m *CommissioningWindowManager) GetCommissioningMode() uint8 {
	return CommissioningMode.Disabled
}
