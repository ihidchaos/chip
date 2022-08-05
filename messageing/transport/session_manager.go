package transport

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/pkg/storage"
	log "github.com/sirupsen/logrus"
	"net/netip"
)

const (
	DuplicateMessageYes uint8 = 0
	DuplicateMessageNo  uint8 = 1
)

type SessionManager interface {
	credentials.FabricTableDelegate
	OnMessageReceived(srcAddr netip.AddrPort, data []byte)
	SecureGroupMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, data []byte)
	SecureUnicastMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, data []byte)
	UnauthenticatedMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, data []byte)
	SetMessageDelegate(SessionMessageDelegate)
}

type SessionManagerImpl struct {
	mUnauthenticatedSessions UnauthenticatedSessionTable
	mSecureSessions          SecureSessionTable
	mFabricTable             *credentials.FabricTableContainer
	mState                   int
	mCB                      SessionMessageDelegate
}

func (s *SessionManagerImpl) FabricWillBeRemoved(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionManagerImpl) OnFabricRemoved(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionManagerImpl) OnFabricCommitted(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionManagerImpl) OnFabricUpdated(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionManagerImpl) SetMessageDelegate(delegate SessionMessageDelegate) {
	s.mCB = delegate
}

func (s *SessionManagerImpl) Init(transports ManagerBase, storage storage.StorageDelegate, table *credentials.FabricTable) error {
	transports.SetSessionManager(s)
	return nil
}

func NewSessionManagerImpl() *SessionManagerImpl {
	return &SessionManagerImpl{}
}

func (s *SessionManagerImpl) OnMessageReceived(srcAddr netip.AddrPort, data []byte) {
	packetHeader := raw.NewPacketHeader()
	err := packetHeader.DecodeAndConsume(data)
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

// UnauthenticatedMessageDispatch 处理没有加密码的消息
func (s *SessionManagerImpl) UnauthenticatedMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, data []byte) {

	source := header.GetSourceNodeId()
	destination := header.GetDestinationNodeId()
	if (source == lib.KUndefinedNodeId && destination == lib.KUndefinedNodeId) || (source != lib.KUndefinedNodeId && destination != lib.KUndefinedNodeId) {
		log.Infof("received malformed unsecure packet with source %d destination %d", source, destination)
		return // ephemeral node id is only assigned to the initiator, there sho
	}

	var unsecuredSession *UnauthenticatedSession
	if source != lib.KUndefinedNodeId {
		unsecuredSession = s.mUnauthenticatedSessions.FindOrAllocateResponder(source, GetLocalMRPConfig())
		if unsecuredSession == nil {
			log.Infof("UnauthenticatedSession exhausted")
			return
		}
	} else {
		unsecuredSession = s.mUnauthenticatedSessions.FindInitiator(destination)
		if unsecuredSession == nil {
			log.Infof("Received unknown unsecure packet for initiator %d", destination)
			return
		}
	}
	unsecuredSession.SetPeerAddress(addr)
	isDuplicate := DuplicateMessageNo

	unsecuredSession.MarkActiveRx()

	var payloadHeader = raw.NewPayloadHeader()
	err := payloadHeader.DecodeAndConsume(data[header.EncodeSizeBytes():])
	if err != nil {
		return
	}

	err = unsecuredSession.GetPeerMessageCounter().VerifyUnencrypted(header.GetMessageCounter())
	if err == lib.ChipErrorDuplicateMessageReceived {
		isDuplicate = DuplicateMessageNo
		log.Infof(
			"Received a duplicate message with MessageCounter: %v on exchange %v",
			header.GetMessageCounter(), payloadHeader)
	} else {
		unsecuredSession.GetPeerMessageCounter().CommitUnencrypted(header.GetMessageCounter())
	}
	if s.mCB != nil {
		s.mCB.OnMessageReceived(header, payloadHeader, unsecuredSession, isDuplicate, data)
	}
}

// SecureGroupMessageDispatch 处理加密的组播消息
func (s *SessionManagerImpl) SecureGroupMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, data []byte) {

}

// SecureUnicastMessageDispatch 处理分支，加密的单播消息
func (s *SessionManagerImpl) SecureUnicastMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, data []byte) {

}
