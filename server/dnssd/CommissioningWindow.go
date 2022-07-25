package dnssd

import "github.com/galenliu/chip/config"

type AppDelegate interface {
	OnCommissioningSessionStarted()
	OnCommissioningSessionStopped()
	OnCommissioningWindowOpened()
	OnCommissioningWindowClosed()
}

type CommissioningWindowManager interface {
	Init(s ServerDelegate) error
	SetAppDelegate(delegate AppDelegate)
	OpenBasicCommissioningWindow() error
	GetCommissioningMode() int
}

type ServerDelegate interface {
}

type CommissioningWindowManagerImpl struct {
	mServer                      ServerDelegate
	mAppDelegate                 AppDelegate
	mFailedCommissioningAttempts uint8
	mUseECM                      bool
}

func (m CommissioningWindowManagerImpl) Init(s ServerDelegate) error {
	m.mServer = s
	return nil
}

func (m *CommissioningWindowManagerImpl) SetAppDelegate(delegate AppDelegate) {
	m.mAppDelegate = delegate
}

func (m *CommissioningWindowManagerImpl) OpenBasicCommissioningWindow() error {
	if config.ChipConfigNetworkLayerBle {
	}
	m.mFailedCommissioningAttempts = 0
	m.mUseECM = false

	err := m.OpenCommissioningWindow()
	if err != nil {
		m.Cleanup()
	}

	//commissioningTimeout := time.Minute * 15

	return err
}

func (m *CommissioningWindowManagerImpl) GetCommissioningMode() int {
	return CommissioningMode_Disabled
}

func (m *CommissioningWindowManagerImpl) OpenCommissioningWindow() error {
	return nil
}

func (m *CommissioningWindowManagerImpl) Cleanup() {

}
