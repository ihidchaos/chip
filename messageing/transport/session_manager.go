package transport

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/pkg/storage"
	log "github.com/sirupsen/logrus"
	"net/netip"
)

type SessionMessageDelegate interface {
	OnMessageReceived(packetHeader *raw.PacketHeader, payloadHeader *raw.PayloadHeader, session SessionHandle, duplicate uint8, buf *raw.PacketBuffer)
}

const (
	kNotReady uint8 = iota
	kInitialized
)

const (
	KDuplicateMessageYes uint8  = 0
	KDuplicateMessageNo  uint8  = 1
	FDuplicateMessage    uint32 = 0x00000001
)

// SessionManager The delegate for TransportManager and FabricTable
// TransportBaseDelegate is the indirect delegate for TransportManager
type SessionManager interface {
	credentials.FabricTableDelegate
	TransportManagerDelegate
	// SecureGroupMessageDispatch  handle the Secure Group messages
	SecureGroupMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *raw.PacketBuffer)
	// SecureUnicastMessageDispatch  handle the unsecure messages
	SecureUnicastMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *raw.PacketBuffer)
	// UnauthenticatedMessageDispatch handle the unauthenticated(未经认证的) messages
	UnauthenticatedMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *raw.PacketBuffer)

	SetMessageDelegate(SessionMessageDelegate)
}

type SessionManagerImpl struct {
	mUnauthenticatedSessions         *UnauthenticatedSessionTable
	mSecureSessions                  *SecureSessionTable
	mFabricTable                     *credentials.FabricTable
	mState                           uint8
	mTransportMgr                    TransportManagerBase
	mGroupClientCounter              *GroupOutgoingCounters
	mCB                              SessionMessageDelegate
	mMessageCounterManager           MessageCounterManagerInterface
	mGlobalUnencryptedMessageCounter *GlobalUnencryptedMessageCounterImpl
}

func NewSessionManagerImpl() *SessionManagerImpl {
	return &SessionManagerImpl{
		mUnauthenticatedSessions:         NewUnauthenticatedSessionTable(),
		mSecureSessions:                  NewSecureSessionTable(),
		mGroupClientCounter:              NewGroupOutgoingCounters(),
		mGlobalUnencryptedMessageCounter: NewGlobalUnencryptedMessageCounterImpl(),
		mFabricTable:                     nil,
		mState:                           0,
		mCB:                              nil,
	}
}

func (s *SessionManagerImpl) SetMessageDelegate(delegate SessionMessageDelegate) {
	s.mCB = delegate
}

func (s *SessionManagerImpl) Init(transportMgr TransportManagerBase, counter MessageCounterManagerInterface, storage storage.KvsPersistentStorageDelegate, table *credentials.FabricTable) error {

	err := s.mFabricTable.AddFabricDelegate(s)
	if err != nil {
		return err
	}

	s.mState = kInitialized

	s.mFabricTable = table

	s.mMessageCounterManager = counter
	s.mSecureSessions.Init()

	s.mGlobalUnencryptedMessageCounter.Init()

	err = s.mGroupClientCounter.Init(storage)
	s.mTransportMgr = transportMgr
	s.mTransportMgr.SetSessionManager(s)
	return err
}

func (s *SessionManagerImpl) OnMessageReceived(srcAddr netip.AddrPort, buf *raw.PacketBuffer) {
	packetHeader := raw.NewPacketHeader()
	err := packetHeader.DecodeAndConsume(buf)
	if err != nil {
		log.Printf("failed to decode packet header: %s", err.Error())
		return
	}
	if packetHeader.IsEncrypted() {
		if packetHeader.IsGroupSession() {
			s.SecureGroupMessageDispatch(packetHeader, srcAddr, buf)
		} else {
			s.SecureUnicastMessageDispatch(packetHeader, srcAddr, buf)
		}
	} else {
		s.UnauthenticatedMessageDispatch(packetHeader, srcAddr, buf)
	}
}

// UnauthenticatedMessageDispatch 处理没有加密码的消息
func (s *SessionManagerImpl) UnauthenticatedMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *raw.PacketBuffer) {

	source := header.GetSourceNodeId()
	destination := header.GetDestinationNodeId()
	if (source == lib.KUndefinedNodeId && destination == lib.KUndefinedNodeId) || (source != lib.KUndefinedNodeId && destination != lib.KUndefinedNodeId) {
		log.Infof("received malformed unsecure packet with source %d destination %d", source, destination)
		return // ephemeral node id is only assigned to the initiator, there sho
	}

	var unsecuredSession *UnauthenticatedSessionImpl
	if source.HasValue() {
		unsecuredSession = s.mUnauthenticatedSessions.FindOrAllocateResponder(source, GetLocalMRPConfig())
		if unsecuredSession == nil {
			log.Infof("UnauthenticatedSessionImpl exhausted")
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
	isDuplicate := KDuplicateMessageNo

	unsecuredSession.MarkActiveRx()

	var payloadHeader = raw.NewPayloadHeader()
	err := payloadHeader.DecodeAndConsume(buf)
	if err != nil {
		return
	}

	err = unsecuredSession.GetPeerMessageCounter().VerifyUnencrypted(header.GetMessageCounter())
	if err != nil {
		isDuplicate = KDuplicateMessageYes
		log.Infof(
			"Received a duplicate message with MessageCounter: %v on exchange %v",
			header.GetMessageCounter(), payloadHeader)
	} else {
		unsecuredSession.GetPeerMessageCounter().CommitUnencrypted(header.GetMessageCounter())
	}
	if s.mCB != nil {
		s.mCB.OnMessageReceived(header, payloadHeader, unsecuredSession, isDuplicate, buf)
	}
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

// SecureGroupMessageDispatch 处理加密的组播消息
func (s *SessionManagerImpl) SecureGroupMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, msg *raw.PacketBuffer) {

}

// SecureUnicastMessageDispatch 处理分支，加密的单播消息
func (s *SessionManagerImpl) SecureUnicastMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, msg *raw.PacketBuffer) {
	secureSession := s.mSecureSessions.FindSecureSessionByLocalKey(header.GetSessionId())
	isDuplicate := KDuplicateMessageNo
	if msg.IsNull() {
		log.Infof("Secure transport received Unicast NULL packet, discarding")
		return
	}
	if secureSession == nil {
		log.Infof("Data received on an unknown session (LSID=%d). Dropping it!", header.GetSessionId())
		return
	}

	if secureSession.IsDefunct() && !secureSession.IsActiveSession() && !secureSession.IsPendingEviction() {
		log.Infof("Secure transport received message on a session in an invalid state (state = '%s')",
			secureSession.GetStateStr())
	}
}
