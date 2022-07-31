package transport

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/pkg"
	"github.com/galenliu/chip/storage"
	"github.com/galenliu/chip/transport/message"
	log "github.com/sirupsen/logrus"
	"net/netip"
)

const (
	NotReady = iota
	Initialized

	DuplicateMessageYes uint8 = 0
	DuplicateMessageNo  uint8 = 1
)

type SessionManager interface {
	OnMessageReceived(srcAddr netip.AddrPort, data []byte)
	Init(transports Transport, storage storage.StorageDelegate, table *credentials.FabricTable) error
	SecureGroupMessageDispatch(header *message.Header, addr netip.AddrPort, data []byte)
	SecureUnicastMessageDispatch(header *message.Header, addr netip.AddrPort, data []byte)
	SetMessageDelegate(message.SessionMessageDelegate)
}

type SessionManagerImpl struct {
	mUnauthenticatedSessions UnauthenticatedSessionTable
	mSecureSessions          SecureSessionTable
	mFabricTable             *credentials.FabricTableContainer
	mState                   int
	mCB                      message.SessionMessageDelegate
}

func (s *SessionManagerImpl) SetMessageDelegate(delegate message.SessionMessageDelegate) {
	s.mCB = delegate
}

func (s *SessionManagerImpl) Init(transports Transport, storage storage.StorageDelegate, table *credentials.FabricTable) error {
	return nil
}

func NewSessionManagerImpl() *SessionManagerImpl {
	return &SessionManagerImpl{}
}

func (s *SessionManagerImpl) OnMessageReceived(srcAddr netip.AddrPort, data []byte) {
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
	if (source == lib.KUndefinedNodeId && destination == lib.KUndefinedNodeId) || (source != lib.KUndefinedNodeId && destination != lib.KUndefinedNodeId) {
		log.Infof("received malformed unsecure packet with source %d destination %d", source, destination)
		return // ephemeral node id is only assigned to the initiator, there sho
	}

	var optionalSession SessionHandle
	if source != lib.KUndefinedNodeId {
		optionalSession = s.mUnauthenticatedSessions.FindOrAllocateResponder(source, messageing.GetLocalMRPConfig())
		if optionalSession == nil {
			log.Infof("UnauthenticatedSession exhausted")
			return
		}
	} else {
		optionalSession = s.mUnauthenticatedSessions.FindInitiator(destination)
		if optionalSession == nil {
			log.Infof("Received unknown unsecure packet for initiator %d", destination)
			return
		}
	}

	unsecuredSession := optionalSession.AsUnauthenticatedSession()
	unsecuredSession.SetPeerAddress(addr)

	isDuplicate := DuplicateMessageNo
	unsecuredSession.MarkActiveRx()
	var payloadHeader message.PayloadHeader
	payloadHeader.DecodeAndConsume(data)
	err := unsecuredSession.GetPeerMessageCounter().VerifyUnencrypted(header.GetMessageCounter())
	if err == pkg.ChipErrorDuplicateMessageReceived {
		log.Infof("Received a duplicate message with MessageCounter: %d", header.GetMessageCounter())
		isDuplicate = DuplicateMessageYes
		return
	} else {
		unsecuredSession.GetPeerMessageCounter().CommitUnencrypted(header.GetMessageCounter())
	}
	if s.mCB != nil {
		s.mCB.OnMessageReceived(header, payloadHeader, unsecuredSession, isDuplicate, data)
	}
}
