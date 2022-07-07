package manager

import (
	"github.com/galenliu/chip/server/dnssd/costants/commissioning_mode"
)

type ServerDelegate interface {
}

type CommissioningWindowManager struct {
	mServer      ServerDelegate
	mAppDelegate any
}

func (m CommissioningWindowManager) Init(s ServerDelegate) *CommissioningWindowManager {
	m.mServer = s
	return &m
}

func (m *CommissioningWindowManager) SetAppDelegate(delegate any) {
	m.mAppDelegate = delegate
}

func (m *CommissioningWindowManager) OpenBasicCommissioningWindow() error {
	return nil
}

func (m *CommissioningWindowManager) GetCommissioningMode() CommissioningMode.T {
	return CommissioningMode.Disabled
}
