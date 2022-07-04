package server

import (
	"github.com/galenliu/chip/config"
)

type CommissioningWindowManager struct {
	mServer      *Server
	mAppDelegate any
}

func NewCommissioningWindowManager(s *Server) *CommissioningWindowManager {
	return &CommissioningWindowManager{
		mServer: s,
	}
}

func (m *CommissioningWindowManager) SetAppDelegate(delegate any) {
	m.mAppDelegate = delegate
}

func (m *CommissioningWindowManager) OpenBasicCommissioningWindow() error {
	return nil
}

func (m *CommissioningWindowManager) GetCommissioningMode() config.CommissioningMode {
	return config.KDisabled
}
