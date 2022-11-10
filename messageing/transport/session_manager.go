package transport

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/pkg/storage"
	"github.com/galenliu/chip/platform/system"
	log "golang.org/x/exp/slog"
	"net/netip"
)

// SessionMessageDelegate 这里的delegate实例为ExchangeManager
type SessionMessageDelegate interface {
	OnMessageReceived(packetHeader *raw.PacketHeader, payloadHeader *raw.PayloadHeader, session *SessionHandle, duplicate uint8, buf *system.PacketBufferHandle)
}

var mGroupPeerMsgCounter = NewGroupPeerTable(ConfigMaxFabrice)

const (
	PayloadIsEncrypted uint8 = iota
	PayloadIsUnencrypted
	DuplicateMessageYes
	DuplicateMessageNo
	NotReady
	kInitialized
	DuplicateMessage uint32 = 0x00000001
)

type EncryptedPacketBufferHandle struct {
	*system.PacketBufferHandle
}

// SessionManagerBase The delegate for TransportManager and FabricTable
// TransportBaseDelegate is the indirect delegate for TransportManager
type SessionManagerBase interface {
	credentials.FabricTableDelegate
	MgrDelegate
	// SecureGroupMessageDispatch  handle the Secure Group messages
	SecureGroupMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *system.PacketBufferHandle)
	// SecureUnicastMessageDispatch  handle the unsecure messages
	SecureUnicastMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *system.PacketBufferHandle)
	// UnauthenticatedMessageDispatch handle the unauthenticated(未经认证的) messages
	UnauthenticatedMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *system.PacketBufferHandle)

	PrepareMessage(session *SessionHandle, payloadHeader *raw.PayloadHeader, msgBuf *system.PacketBufferHandle, encryptedMessage *EncryptedPacketBufferHandle) error
	SendPreparedMessage(session *SessionHandle, preparedMessage *EncryptedPacketBufferHandle) error
	AllocateSession(sessionType SecureSessionType, sessionEvictionHint *lib.ScopedNodeId) *SessionHandle

	ExpireAllSessions(node *lib.ScopedNodeId)
	ExpireAllSessionsForFabric(fabricIndex lib.FabricIndex)
	ExpireAllSessionsOnLogicalFabric(node *lib.ScopedNodeId)

	CreateUnauthenticatedSession(peerAddress netip.AddrPort, config *ReliableMessageProtocolConfig) *SessionHandle
	FindSecureSessionForNode(nodeId *lib.ScopedNodeId, sessionType SecureSessionType) *SessionHandle

	SetMessageDelegate(cb SessionMessageDelegate)
	SystemLayer() system.Layer
}

type SessionManager struct {
	mUnauthenticatedSessions *UnauthenticatedSessionTable
	mSecureSessions          *SecureSessionTable
	mFabricTable             *credentials.FabricTable
	mState                   uint8
	mTransportMgr            MgrBase
	mMessageCounterManager   MessageCounterManagerBase

	mGlobalUnencryptedMessageCounter *GlobalUnencryptedMessageCounter
	mGroupClientCounter              *GroupOutgoingCounters
	mCB                              SessionMessageDelegate
	mSystemLayer                     system.Layer
}

func NewSessionManagerImpl() *SessionManager {
	return &SessionManager{
		mUnauthenticatedSessions:         NewUnauthenticatedSessionTable(),
		mSecureSessions:                  NewSecureSessionTable(),
		mGroupClientCounter:              NewGroupOutgoingCounters(),
		mGlobalUnencryptedMessageCounter: NewGlobalUnencryptedMessageCounterImpl(),
		mFabricTable:                     credentials.NewFabricTable(),
		mState:                           0,
		mCB:                              nil,
	}
}

func (s *SessionManager) SetMessageDelegate(delegate SessionMessageDelegate) {
	s.mCB = delegate
}

func (s *SessionManager) Init(systemLay system.Layer, transportMgr MgrBase, counter MessageCounterManagerBase, storage storage.KvsPersistentStorageDelegate, table *credentials.FabricTable) error {

	s.mState = kInitialized
	s.mSystemLayer = systemLay

	err := s.mFabricTable.AddFabricDelegate(s)
	if err != nil {
		return err
	}

	s.mFabricTable = table

	s.mMessageCounterManager = counter
	s.mSecureSessions.Init()

	s.mGlobalUnencryptedMessageCounter.Init()

	err = s.mGroupClientCounter.Init(storage)
	s.mTransportMgr = transportMgr
	s.mTransportMgr.SetSessionManager(s)
	return err
}

func (s *SessionManager) OnMessageReceived(srcAddr netip.AddrPort, msg *system.PacketBufferHandle) {
	packetHeader := raw.NewPacketHeader()
	err := packetHeader.DecodeAndConsume(msg)
	if err != nil {
		log.Error("failed to decode packet header", err, "Tag", "SessionManager")
		return
	}
	if packetHeader.IsEncrypted() {
		if packetHeader.IsGroupSession() {
			s.SecureGroupMessageDispatch(packetHeader, srcAddr, msg)
		} else {
			s.SecureUnicastMessageDispatch(packetHeader, srcAddr, msg)
		}
	} else {
		s.UnauthenticatedMessageDispatch(packetHeader, srcAddr, msg)
	}
}

// UnauthenticatedMessageDispatch 处理没有加密码的消息
func (s *SessionManager) UnauthenticatedMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *system.PacketBufferHandle) {

	source := header.SourceNodeId
	destination := header.DestinationNodeId

	if (source.HasValue() && destination.HasValue()) || (!source.HasValue() && !destination.HasValue()) {
		log.Info("received malformed unsecure packet with source %d destination %d", source, destination)
		return
		//ephemeral node id is only assigned to the initiator, there should be one and only one node id exists.
	}

	var unsecuredSession *UnauthenticatedSession
	if source.HasValue() {
		// Assume peer is the initiator, we are the responder.
		// 对方是发起人，我们是响应者
		unsecuredSession = s.mUnauthenticatedSessions.FindOrAllocateResponder(source, GetLocalMRPConfig())
		if unsecuredSession == nil {
			log.Info("UnauthenticatedSession exhausted")
			return
		}
	} else {
		// Assume peer is the responder, we are the initiator.
		// 对方为响应，我们是发起人
		unsecuredSession = s.mUnauthenticatedSessions.FindInitiator(destination)
		if unsecuredSession == nil {
			log.Info("Received unknown unsecure packet for initiator", "destination", destination)
			return
		}
	}

	unsecuredSession.SetPeerAddress(addr)
	isDuplicate := DuplicateMessageNo
	// 更新Session
	unsecuredSession.MarkActiveRx()

	payloadHeader := raw.NewPayloadHeader()
	err := payloadHeader.DecodeAndConsume(buf)
	if err != nil {
		log.Info("Received invalid packet")
		return
	}

	err = unsecuredSession.GetPeerMessageCounter().VerifyUnencrypted(header.MessageCounter)
	if err != nil && err == lib.MatterErrorDuplicateMessageReceived {
		isDuplicate = DuplicateMessageYes
		log.Info(
			"Received a duplicate message", "MessageCounter", header.MessageCounter, "payloadHeader", payloadHeader)
	} else {
		unsecuredSession.GetPeerMessageCounter().CommitUnencrypted(header.MessageCounter)
	}
	if s.mCB != nil {
		s.mCB.OnMessageReceived(header, payloadHeader, NewSessionHandle(unsecuredSession), isDuplicate, buf)
	}
}

func (s *SessionManager) FabricWillBeRemoved(table *credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionManager) SystemLayer() system.Layer {
	return s.mSystemLayer
}

func (s *SessionManager) OnFabricCommitted(table *credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionManager) OnFabricUpdated(table *credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

// SecureGroupMessageDispatch 处理加密的组播消息
func (s *SessionManager) SecureGroupMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, msg *system.PacketBufferHandle) {

}

func (s *SessionManager) Shutdown() {
	if s.mFabricTable != nil {
		s.mFabricTable.RemoveFabricDelegate(s)
		s.mFabricTable = nil
	}
	s.mState = NotReady
	//s.mSecureSessions
	s.mMessageCounterManager = nil
	s.mSystemLayer = nil
	s.mTransportMgr = nil
	s.mCB = nil
}

func (s *SessionManager) FabricRemoved(fabricId lib.FabricIndex) {
	err := mGroupPeerMsgCounter.FabricRemoved(fabricId)
	if err != nil {
		log.Error("SessionManager.FabricRemoved", err)
	}
}

// SecureUnicastMessageDispatch 处理分支，加密的单播消息
func (s *SessionManager) SecureUnicastMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, msg *system.PacketBufferHandle) {
	secureSession := s.mSecureSessions.FindSecureSessionByLocalKey(header.SessionId)
	_ = DuplicateMessageNo
	if msg.IsNull() {
		log.Info("Secure transport received Unicast NULL packet, discarding")
		return
	}
	if secureSession == nil {
		log.Info("Data received on an unknown session (LSID=%d). Dropping it!", header.SessionId)
		return
	}

	if secureSession.IsDefunct() && !secureSession.IsActiveSession() && !secureSession.IsPendingEviction() {
		log.Info("Secure transport received message on a session in an invalid state (state = '%s')",
			secureSession.State())
	}
	var nodeId = lib.UndefinedNodeId
	if secureSession.GetSecureSessionType() == CASE {
		nodeId = secureSession.GetPeerNodeId()
	}
	nonce, _ := BuildNonce(header.SecFlags, header.MessageCounter, nodeId)
	_ = Decrypt(secureSession.GetCryptoContext(), nonce, header, msg)
}

func (s *SessionManager) PrepareMessage(session *SessionHandle, payloadHeader *raw.PayloadHeader, msgBuf *system.PacketBufferHandle, encryptedMessage *EncryptedPacketBufferHandle) error {
	return nil
}

func (s *SessionManager) SendPreparedMessage(session *SessionHandle, preparedMessage *EncryptedPacketBufferHandle) error {
	return nil
}

func (s *SessionManager) AllocateSession(sessionType SecureSessionType, sessionEvictionHint *lib.ScopedNodeId) *SessionHandle {
	return nil
}

func (s *SessionManager) ExpireAllSessions(node *lib.ScopedNodeId) {

}
func (s *SessionManager) ExpireAllSessionsForFabric(fabricIndex lib.FabricIndex) {

}

func (s *SessionManager) ExpireAllSessionsOnLogicalFabric(node *lib.ScopedNodeId) {

}

func (s *SessionManager) CreateUnauthenticatedSession(peerAddress netip.AddrPort, config *ReliableMessageProtocolConfig) *SessionHandle {
	return nil
}

func (s *SessionManager) FindSecureSessionForNode(nodeId *lib.ScopedNodeId, sessionType SecureSessionType) *SessionHandle {
	return nil
}

func (s *SessionManager) OnFabricRemoved(table *credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}
