package dnssd

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/server/dnssd/costants/discovery"
	"math/rand"
	"net"
	"time"
)

const (
	SubtypeServiceNamePart    = "_sub"
	CommissionableServiceName = "_matterc"
	CommissionerServiceName   = "_matterd"
	OperationalServiceName    = "_matter"
	CommissionProtocol        = "_udp"
	LocalDomain               = "local"
	OperationalProtocol       = "_tcp"
)

const kMaxRetryInterval = time.Millisecond * 3600000

type CommissioningModeProvider interface {
	GetCommissioningMode() uint8
}

func (s Server) SetFabricTable(f FabricTable) {
	s.mFabrics = f
}

func (s *Server) SetCommissioningModeProvider(manager CommissioningModeProvider) {
	s.mCommissioningModeProvider = manager
}

func (s *Server) SetSecuredPort(port uint16) {
	s.mSecuredPort = port
}

func (s *Server) SetUnsecuredPort(port uint16) {
	s.mUnsecuredPort = port
}

func (s *Server) GetSecuredPort() uint16 {
	return s.mSecuredPort
}

func (s *Server) GetUnsecuredPort() uint16 {
	return s.mUnsecuredPort
}

func (s *Server) SetInterfaceId(inter net.Interface) {
	s.mInterfaceId = inter
}

func (s *Server) GetInterfaceId() net.Interface {
	return s.mInterfaceId
}

func (s *Server) HandleExtendedDiscoveryExpiration() {

}

func (s *Server) OnExtendedDiscoveryExpiration() {

}

func (s *Server) UpdateCommissionableInstanceName() {
	//0 is not value ,so this must check
	s.mCommissionableInstanceName = rand.Uint64()
	for s.mCurrentCommissioningMode == 0 {
		s.mCommissionableInstanceName = rand.Uint64()
	}
}

func (s *Server) GetCommissionableInstanceName() uint64 {
	if s.mCommissionableInstanceName == 0 {
		s.UpdateCommissionableInstanceName()
	}
	return s.mCommissionableInstanceName
}

func (s *Server) GetExtendedDiscoveryTimeoutSecs() time.Duration {
	if s.mExtendedDiscoveryTimeoutSecs == 0 {
		return time.Duration(config.ChipDeviceConfigExtendedDiscoveryTimeoutSecs) * time.Second
	}
	return s.mExtendedDiscoveryTimeoutSecs
}

func (s Server) HaveOperationalCredentials() bool {
	if s.mFabrics == nil {
		return false
	}
	return s.mFabrics.FabricCount() != 0
}

func (s Server) ScheduleExtendedDiscoveryExpiration() {

}

func MakeServiceSubtype[T discovery.Uint](filter discovery.FilterType, code T) (string, error) {
	switch filter {
	case discovery.ShortDiscriminator:
		//if code >= 1<<4 {
		//	return "", fmt.Errorf("chip error invalid argument")
		//}
		return fmt.Sprintf("_S%d", code), nil
	case discovery.LongDiscriminator:
		//if code >= 1<<12 {
		//	return "", fmt.Errorf("chip error invalid argument")
		//}
		return fmt.Sprintf("_L%d", code), nil
	case discovery.VendorId:
		//if code >= 1<<16 {
		//	return "", fmt.Errorf("chip error invalid argument")
		//}
		return fmt.Sprintf("_V%d", code), nil
	case discovery.DeviceType:
		return fmt.Sprintf("_T%d", code), nil
	case discovery.CommissioningMode:
		return "_CM", nil
	case discovery.Commissioner:
		//if code > 1 {
		//	return "", fmt.Errorf("chip error invalid argument")
		//}
		return fmt.Sprintf("_D%d", code), nil
	case discovery.CompressedFabricId:
		return fmt.Sprintf("_I%016X", code), nil
	case discovery.InstanceName:
		return fmt.Sprintf("%016X", code), nil
	case discovery.None:
		return string([]byte{0}), nil
	}
	return "", nil
}
