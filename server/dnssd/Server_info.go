package DnssdServer

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"math/rand"
	"net"
	"time"
)

const (
	KSubtypeServiceNamePart    = "_sub"
	KCommissionableServiceName = "_matterc"
	KCommissionerServiceName   = "_matterd"
	KOperationalServiceName    = "_matter"
	KCommissionProtocol        = "_udp"
	KLocalDomain               = "local"
	KOperationalProtocol       = "_tcp"
)

const kMaxRetryInterval = time.Millisecond * 3600000

type CommissioningModeProvider interface {
	GetCommissioningMode() int
}

func (s *ServerImpl) SetFabricTable(f FabricTable) {
	s.mFabricTable = f
}

func (s *ServerImpl) SetCommissioningModeProvider(manager CommissioningModeProvider) {
	s.mCommissioningModeProvider = manager
}

func (s *ServerImpl) SetSecuredPort(port uint16) {
	s.mSecuredPort = port
}

func (s *ServerImpl) SetUnsecuredPort(port uint16) {
	s.mUnsecuredPort = port
}

func (s *ServerImpl) GetSecuredPort() uint16 {
	return s.mSecuredPort
}

func (s *ServerImpl) GetUnsecuredPort() uint16 {
	return s.mUnsecuredPort
}

func (s *ServerImpl) SetInterfaceId(inter net.Interface) {
	s.mInterfaceId = inter
}

func (s *ServerImpl) GetInterfaceId() net.Interface {
	return s.mInterfaceId
}

func (s *ServerImpl) UpdateCommissionableInstanceName() {
	//0 is not value ,so this must check
	s.mCommissionableInstanceName = fmt.Sprintf("%016X", rand.Uint64())

}

func (s *ServerImpl) GetCommissionableInstanceName() (string, error) {
	if s.mCommissionableInstanceName == "" {
		s.UpdateCommissionableInstanceName()
	}
	return s.mCommissionableInstanceName, nil
}

func (s *ServerImpl) GetExtendedDiscoveryTimeoutSecs() time.Duration {
	if s.mExtendedDiscoveryTimeoutSecs == 0 {
		return time.Duration(config.ChipDeviceConfigExtendedDiscoveryTimeoutSecs) * time.Second
	}
	return s.mExtendedDiscoveryTimeoutSecs
}

func (s ServerImpl) haveOperationalCredentials() bool {
	if s.mFabricTable == nil {
		return false
	}
	return s.mFabricTable.FabricCount() != 0
}

func (s ServerImpl) ScheduleExtendedDiscoveryExpiration() {

}
