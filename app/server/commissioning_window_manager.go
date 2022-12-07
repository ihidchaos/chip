package server

import (
	"github.com/galenliu/chip/messageing/transport"
	"time"
)

type CommissioningModeProvider interface {
	GetCommissioningMode() int
}

type CommissioningWindowManager struct {
	server *Server
}

func (c *CommissioningWindowManager) Init(server *Server) {
	c.server = server
}

func (c *CommissioningWindowManager) OpenBasicCommissioningWindow(
	commissioningTimeout time.Duration,
	advertisementMode CommissioningWindowAdvertisement) error {
	if err := c.OpenCommissioningWindow(commissioningTimeout); err != nil {
		return err
	}
	return nil
}

func (c *CommissioningWindowManager) OnSessionReleased() {
	//TODO implement me
	panic("implement me")
}

func (c *CommissioningWindowManager) OnSessionHang() {
	//TODO implement me
	panic("implement me")
}

func (c *CommissioningWindowManager) GetNewSessionHandlingPolicy() transport.NewSessionHandlingPolicy {
	//TODO implement me
	panic("implement me")
}

func (c *CommissioningWindowManager) OnSessionEstablishmentError() {
	//TODO implement me
	panic("implement me")
}

func (c *CommissioningWindowManager) OnSessionEstablishmentStarted() {
	//TODO implement me
	panic("implement me")
}

func (c *CommissioningWindowManager) OnSessionEstablished() {
	//TODO implement me
	panic("implement me")
}

func (c *CommissioningWindowManager) OpenCommissioningWindow(timeout time.Duration) error {
	return nil
}
