package core

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/messageing"
	transport2 "github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/secure_channel"
)

type DeviceProxyInitParams struct {
	SessionManager            transport2.SessionManager
	SessionResumptionStorage  secure_channel.SessionResumptionStorage
	CertificateValidityPolicy credentials.CertificateValidityPolicy
	ExchangeMgr               messageing.ExchangeManager
	FabricTable               *credentials.FabricTable
	ClientPool                CASEClientPoolDelegate
	GroupDataProvider         credentials.GroupDataProvider

	MrpLocalConfig *transport2.ReliableMessageProtocolConfig
}
