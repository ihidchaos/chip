package transport

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/storage"
	"github.com/galenliu/chip/transport/message"
	log "github.com/sirupsen/logrus"
	"net/netip"
)

type SessionManager interface {
	Init(transports Transport, storage storage.StorageDelegate, table *credentials.FabricTable) error
}

type SessionManagerImpl struct {
	mUnauthenticatedSessions UnauthenticatedSessionTable
}

func (s SessionManagerImpl) Init(transports Transport, storage storage.StorageDelegate, table *credentials.FabricTable) error {
	return nil
}

func NewSessionManagerImpl() *SessionManagerImpl {
	return &SessionManagerImpl{}
}

func (s SessionManagerImpl) OnMessageReceived(srcAddr netip.AddrPort, data []byte) {
	packetHeader, err := message.DecodeHeader(data)
	if err != nil {
		log.Printf("failed to decode packet header: %s", err.Error())
		return
	}
	if packetHeader.IsEncrypted() {
		if packetHeader.IsGroupSession() {
			s.SecureGroupMessageDispatch(packetHeader, srcAddr, data)
		} else {
			s.SecureUnicastMessageDispatch(packetHeader, srcAddr, data)
		}
	} else {
		s.UnauthenticatedMessageDispatch(packetHeader, srcAddr, data)
	}

}

func (s *SessionManagerImpl) SecureGroupMessageDispatch(header *message.Header, addr netip.AddrPort, data []byte) {

}

func (s *SessionManagerImpl) SecureUnicastMessageDispatch(header *message.Header, addr netip.AddrPort, data []byte) {

}

func (s *SessionManagerImpl) UnauthenticatedMessageDispatch(header *message.Header, addr netip.AddrPort, data []byte) {
	source := header.GetSourceNodeId()
	destination := header.GetDestinationNodeId()
	if (source == 0 && destination == 0) || (source != 0 && destination != 0) {
		log.Infof("received malformed unsecure packet with source %d destination %d", source, destination)
		return // ephemeral node id is only assigned to the initiator, there sho
	}
	if source != 0 {
		optionalSession := s.mUnauthenticatedSessions.FindOrAllocateResponder(source, messageing.GetLocalMRPConfig())
		if optionalSession == nil {
			return
		}
	}
}
