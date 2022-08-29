package core

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/messageing"
	transport2 "github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/protocols/secure"
)

type DeviceProxyInitParams struct {
	SessionManager            transport2.SessionManager
	SessionResumptionStorage  secure.SessionResumptionStorage
	CertificateValidityPolicy credentials.CertificateValidityPolicy
	ExchangeMgr               messageing.ExchangeManager
	FabricTable               *credentials.FabricTable
	ClientPool                CASEClientPoolDelegate
	GroupDataProvider         credentials.GroupDataProvider

	MrpLocalConfig *transport2.ReliableMessageProtocolConfig
}
