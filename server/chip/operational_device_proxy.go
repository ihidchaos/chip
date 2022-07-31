package chip

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/secure_channel"
	"github.com/galenliu/chip/transport"
)

type DeviceProxyInitParams struct {
	SessionManager            transport.SessionManager
	SessionResumptionStorage  secure_channel.SessionResumptionStorage
	CertificateValidityPolicy credentials.CertificateValidityPolicy
	ExchangeMgr               messageing.ExchangeManager
	FabricTable               *credentials.FabricTable
	ClientPool                CASEClientPoolDelegate
	GroupDataProvider         credentials.GroupDataProvider

	MrpLocalConfig *messageing.ReliableMessageProtocolConfig
}
