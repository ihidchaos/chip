package secure_channel

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/transport"
)

type CASEServer struct {
}

func NewCASEServer() *CASEServer {
	return &CASEServer{}
}

func (s CASEServer) ListenForSessionEstablishment(mgr messageing.ExchangeManager, sessions transport.SessionManager, fabrics *credentials.FabricTable, storage any, policy credentials.CertificateValidityPolicy, provider credentials.GroupDataProvider) error {
	return nil
}
