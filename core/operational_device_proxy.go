package core

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/messageing"
	transport2 "github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/protocols/secure_channel"
)

type DeviceProxyInitParams struct {
	SessionManager            transport2.SessionManagerBase
	SessionResumptionStorage  secure_channel.SessionResumptionStorage
	CertificateValidityPolicy credentials.CertificateValidityPolicy
	ExchangeMgr               messageing.ExchangeManagerBase
	FabricTable               *credentials.FabricTable
	ClientPool                CASEClientPoolDelegate
	GroupDataProvider         *credentials.GroupDataProvider

	MrpLocalConfig *session.ReliableMessageProtocolConfig
}
