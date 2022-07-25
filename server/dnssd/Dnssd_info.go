package dnssd

import (
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

//
//func (mdns *MdnsServerImpl) SetFabricTable(f FabricTable) {
//	mdns.mFabricTable = f
//}
//
//func (mdns *MdnsServerImpl) SetCommissioningModeProvider(manager CommissioningModeProvider) {
//	mdns.mCommissioningModeProvider = manager
//}
//
//func (mdns *MdnsServerImpl) SetSecuredPort(port uint16) {
//	mdns.mSecuredPort = port
//}
//
//func (mdns *MdnsServerImpl) SetUnsecuredPort(port uint16) {
//	mdns.mUnsecuredPort = port
//}
//
//func (mdns *MdnsServerImpl) GetSecuredPort() uint16 {
//	return mdns.mSecuredPort
//}
//
//func (mdns *MdnsServerImpl) GetUnsecuredPort() uint16 {
//	return mdns.mUnsecuredPort
//}
//
//func (mdns *MdnsServerImpl) SetInterfaceId(inter net.Interface) {
//	mdns.mInterfaceId = inter
//}
//
//func (mdns *MdnsServerImpl) GetInterfaceId() net.Interface {
//	return mdns.mInterfaceId
//}
//
//func (mdns *MdnsServerImpl) UpdateCommissionableInstanceName() {
//	//0 is not value ,so this must check
//	mdns.mCommissionableInstanceName = fmt.Sprintf("%016X", rand.Uint64())
//
//}
//
//func (mdns *MdnsServerImpl) GetCommissionableInstanceName() (string, error) {
//	if mdns.mCommissionableInstanceName == "" {
//		mdns.UpdateCommissionableInstanceName()
//	}
//	return mdns.mCommissionableInstanceName, nil
//}
//
//func (mdns *MdnsServerImpl) GetExtendedDiscoveryTimeoutSecs() time.Duration {
//	if mdns.mExtendedDiscoveryTimeoutSecs == 0 {
//		return time.Duration(config.ChipDeviceConfigExtendedDiscoveryTimeoutSecs) * time.Second
//	}
//	return mdns.mExtendedDiscoveryTimeoutSecs
//}
//
//func (mdns MdnsServerImpl) haveOperationalCredentials() bool {
//	if mdns.mFabricTable == nil {
//		return false
//	}
//	return mdns.mFabricTable.FabricCount() != 0
//}
//
//func (mdns MdnsServerImpl) ScheduleExtendedDiscoveryExpiration() {
//
//}
