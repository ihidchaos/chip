package transport

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/storage"
	"github.com/galenliu/chip/transport/message"
	log "github.com/sirupsen/logrus"
	"net/netip"
)

type SessionManager interface {
	Init(transports Transport, storage storage.StorageDelegate, table *credentials.FabricTable) error
}

type SessionManagerImpl struct {
}

func (s SessionManagerImpl) Init(transports Transport, storage storage.StorageDelegate, table *credentials.FabricTable) error {
	return nil
}

func NewSessionManagerImpl() *SessionManagerImpl {
	return &SessionManagerImpl{}
}

func (s SessionManagerImpl) OnMessageReceived(port netip.AddrPort, data []byte) {
	packetHeadr, err := message.DecodeHeader(data)
	if err != nil {
		log.Printf("failed to decode packet header: %s", err.Error())
		return
	}
	if packetHeadr.IsEncrypted() {
		if packetHeadr.IsGroupSession() {

		}
	}

}
