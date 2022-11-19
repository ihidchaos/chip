package transport

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/pkg/storage"
	"github.com/galenliu/chip/platform/system"
	log "golang.org/x/exp/slog"
	"net/netip"
)

// SessionMessageDelegate 这里的delegate实例为ExchangeManager
type SessionMessageDelegate interface {
	OnMessageReceived(packetHeader *raw.PacketHeader, payloadHeader *raw.PayloadHeader, session *SessionHandle, isDuplicate bool, buf *system.PacketBufferHandle)
}

var mGroupPeerMsgCounter = NewGroupPeerTable(ConfigMaxFabrice)

const (
	PayloadIsEncrypted uint8 = iota
	PayloadIsUnencrypted
	kNotReady
	kInitialized
)

type EncryptedPacketBufferHandle struct {
	*system.PacketBufferHandle
}

func (e *EncryptedPacketBufferHandle) MarkEncrypted() *system.PacketBufferHandle {
	return e.PacketBufferHandle
}

func (e *EncryptedPacketBufferHandle) MessageCounter() uint32 {
	header := raw.NewPacketHeader()
	err := header.Decode([]byte{})
	if err != nil {
		return 0
	}
	return header.MessageCounter
}

// SessionManagerBase The delegate for TransportManager and FabricTable
// TransportBaseDelegate is the indirect delegate for TransportManager
type SessionManagerBase interface {
	credentials.FabricTableDelegate
	MgrDelegate
	// SecureGroupMessageDispatch  handle the kSecure group messages
	SecureGroupMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *system.PacketBufferHandle)
	// SecureUnicastMessageDispatch  handle the unsecure messages
	SecureUnicastMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *system.PacketBufferHandle)
	// UnauthenticatedMessageDispatch handle the unauthenticated(未经认证的) messages
	UnauthenticatedMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *system.PacketBufferHandle)

	PrepareMessage(session *SessionHandle, payloadHeader *raw.PayloadHeader, msgBuf *system.PacketBufferHandle, encryptedMessage *EncryptedPacketBufferHandle) error
	SendPreparedMessage(session *SessionHandle, preparedMessage *EncryptedPacketBufferHandle) error
	AllocateSession(sessionType session.SecureSessionType, sessionEvictionHint *lib.ScopedNodeId) *SessionHandle

	ExpireAllSessions(node *lib.ScopedNodeId)
	ExpireAllSessionsForFabric(fabricIndex lib.FabricIndex)
	ExpireAllSessionsOnLogicalFabric(node *lib.ScopedNodeId)

	FabricTable() *credentials.FabricTable

	CreateUnauthenticatedSession(peerAddress netip.AddrPort, config *session.ReliableMessageProtocolConfig) *SessionHandle
	FindSecureSessionForNode(nodeId *lib.ScopedNodeId, sessionType session.SecureSessionType) *SessionHandle

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

func NewSessionManager() *SessionManager {
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
		log.Error("packet header failed", err, "Tag", "SessionManager")
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
func (s *SessionManager) UnauthenticatedMessageDispatch(packetHeader *raw.PacketHeader, peerAddress netip.AddrPort, msg *system.PacketBufferHandle) {

	sourceNodeId := packetHeader.SourceNodeId
	destinationNodeId := packetHeader.DestinationNodeId
	if (sourceNodeId.IsSome() && destinationNodeId.IsSome()) || (sourceNodeId.IsNone() && destinationNodeId.IsNone()) {
		log.Info("received malformed unsecure packet", "SourceNodeId", sourceNodeId, "DestinationNodeId", destinationNodeId, "Tag", "SessionManager")
		return
		//ephemeral node id is only assigned to the initiator, there should be one and only one node id exists.
	}
	var optionalSession *SessionHandle = nil
	var err error = nil
	if sourceNodeId.IsSome() {
		// Assume peer is the initiator, we are the responder.
		// 对方是发起人，我们是响应者
		optionalSession, err = s.mUnauthenticatedSessions.FindOrAllocateResponder(sourceNodeId.Unwrap(), session.GetLocalMRPConfig())
		if err != nil {
			log.Error("Unauthenticated exhausted", err)
			return
		}
	} else {
		// Assume peer is the responder, we are the initiator.
		// 对方为响应，我们是发起人
		optionalSession = s.mUnauthenticatedSessions.FindInitiator(destinationNodeId.Unwrap())
		if optionalSession == nil {
			log.Info("Received unknown unsecure packet for initiator", "DestinationNodeId", destinationNodeId, "Tag", "SessionManager")
			return
		}
	}
	var unsecuredSession *session.Unauthenticated
	unsecuredSession = optionalSession.Session.(*session.Unauthenticated)
	unsecuredSession.SetPeerAddress(peerAddress)
	isDuplicate := false
	// 更新Session
	unsecuredSession.MarkActiveRx()

	payloadHeader := raw.NewPayloadHeader()
	err = payloadHeader.DecodeAndConsume(msg)
	if err != nil {
		log.Error("Received invalid packet", err, "Tag", "SessionManager")
		return
	}

	err = unsecuredSession.PeerMessageCounter().VerifyUnencrypted(packetHeader.MessageCounter)
	if err == lib.DuplicateMessageReceived {
		log.Info(
			"Received a duplicate message", "messageCounter", packetHeader.MessageCounter, "payloadHeader", payloadHeader)
		isDuplicate = true
		err = nil
	} else {
		unsecuredSession.PeerMessageCounter().CommitUnencrypted(packetHeader.MessageCounter)
	}
	if s.mCB != nil {
		log.Debug("Message Received", "payloadHeader", payloadHeader, "packetHeader", packetHeader, "unsecuredSession", unsecuredSession, "peerAddress", peerAddress, "message", msg)
		s.mCB.OnMessageReceived(packetHeader, payloadHeader, NewSessionHandle(unsecuredSession), isDuplicate, msg)
	}
}

// SecureUnicastMessageDispatch 处理分支，加密的单播消息
func (s *SessionManager) SecureUnicastMessageDispatch(packetHeader *raw.PacketHeader, peerAddress netip.AddrPort, msg *system.PacketBufferHandle) {

	sessionHandle := s.mSecureSessions.FindSecureSessionByLocalKey(packetHeader.SessionId)
	isDuplicate := false
	if msg.IsNull() {
		log.Info("kSecure transport received unicast NULL packet, discarding")
		return
	}
	if sessionHandle == nil {
		log.Info("Data received on an unknown session, Dropping it!", "LSID", packetHeader.SessionId)
		return
	}

	secureSession := sessionHandle.Session.(*session.Secure)

	if secureSession.IsDefunct() && !secureSession.IsActiveSession() && !secureSession.IsPendingEviction() {
		log.Info("kSecure transport received message on a session in an invalid state", "stata", secureSession.State())
		return
	}

	nonce := session.BuildNonce(packetHeader.SecurityFlags(), packetHeader.MessageCounter, func() lib.NodeId {
		if secureSession.SecureSessionType() == session.SecureSessionTypeCASE {
			return secureSession.PeerNodeId()
		}
		return lib.UndefinedNodeId()
	}())

	payloadHeader, err := Decrypt(secureSession.GetCryptoContext(), nonce, packetHeader, msg)
	if err != nil {
		log.Error("Secure transport received message, but failed to decode/authenticate it", err)
	}

	err = secureSession.SessionMessageCounter().VerifyEncryptedUnicast(packetHeader.MessageCounter)
	if err == lib.DuplicateMessageReceived {
		log.Info("Received a duplicate message on exchange", "MessageCounter", packetHeader.MessageCounter, "PayloadHeader", payloadHeader)
		isDuplicate = true
		err = nil
	}
	if err != nil {
		log.Error("Message counter verify failed", err, "Tag", "Inet")
		return
	}
	secureSession.MarkActiveRx()
	if isDuplicate && !payloadHeader.NeedsAck() {
		//如果这是一个重复的消息，且不需要ACK，则直接返回，节约CPU资源。
		return
	}

	if !isDuplicate {
		secureSession.SessionMessageCounter().PeerMessageCounter().CommitUnencryptedUnicast(packetHeader.MessageCounter)
	}

	if secureSession.PeerAddress() != peerAddress {
		secureSession.SetPeerAddress(peerAddress)
	}

	if s.mCB != nil {
		log.Debug("Message Received", "PayloadHeader", payloadHeader, "packetHeader", packetHeader, "secureSession", secureSession, "PeerAddress", peerAddress, "message", msg)
		s.mCB.OnMessageReceived(packetHeader, payloadHeader, NewSessionHandle(secureSession), isDuplicate, msg)
	}
}

// SecureGroupMessageDispatch 处理加密的组播消息
func (s *SessionManager) SecureGroupMessageDispatch(packetHeader *raw.PacketHeader, peerAddress netip.AddrPort, msg *system.PacketBufferHandle) {
	groups := credentials.GetGroupDataProvider()
	if groups == nil {
		return
	}

	if packetHeader.DestinationGroupId.IsNone() {
		return
	}
	if msg.IsNull() {
		log.Info("Secure transport received Groupcast NULL packet,discarding")
		return
	}
	// Check if Message Header is valid first
	if !(packetHeader.IsValidMCSPMsg() || packetHeader.IsValidGroupMsg()) {
		return
	}

	// Trial decryption with GroupDataProvider
	var err error
	var payloadHeader *raw.PayloadHeader
	var nonce session.NonceStorage
	var groupContext *credentials.GroupSession
	for _, groupContext = range groups.GroupSessions(packetHeader.SessionId) {
		if packetHeader.DestinationGroupId.Unwrap() != groupContext.GroupId {
			continue
		}
		nonce = session.BuildNonce(packetHeader.SecurityFlags(), packetHeader.MessageCounter, packetHeader.SourceNodeId.Unwrap())
		payloadHeader, err = Decrypt(session.NewCryptoContext(groupContext.Key), nonce, packetHeader, msg)
	}
	if err != nil {
		log.Error("Failed to retrieve Key. Discarding everything", err)
		return
	}

	if packetHeader.IsValidMCSPMsg() {
		return
	}
	if payloadHeader.NeedsAck() {
		log.Info("Unexpected ACK requested for group message")
		return
	}
	counter, err := mGroupPeerMsgCounter.FindOrAddPeer(groupContext.FabricIndex, packetHeader.SourceNodeId.Unwrap(), packetHeader.IsSecureSessionControlMsg())
	if err == nil {
		if groupContext.SecurityPolicy == credentials.TrustFirst {
			err = counter.VerifyOrTrustFirstGroup(packetHeader.MessageCounter)
		} else {
			// TODO support cache and sync with MCSP. Issue  #11689
			log.Info("Received Group Msg with key policy Cache and Sync, but MCSP is not implemented")
			return
		}
		if err != nil {
			log.Error("Message counter verify failed", err)
			return
		}

	} else {
		log.Info("Group Counter Tables full or invalid NodeId/FabricIndex after decryption of message, dropping everything")
		return
	}
	counter.CommitGroup(packetHeader.MessageCounter)
	if s.mCB != nil {
		groupSession := session.NewIncomingGroupSession(groupContext.GroupId, groupContext.FabricIndex, packetHeader.SourceNodeId.Unwrap())
		log.Debug("Message Received", "PayloadHeader", payloadHeader, "packetHeader", packetHeader, "GroupSession", groupSession, "PeerAddress", peerAddress, "message", msg)
		s.mCB.OnMessageReceived(packetHeader, payloadHeader, NewSessionHandle(groupSession), false, msg)
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

func (s *SessionManager) Shutdown() {
	if s.mFabricTable != nil {
		s.mFabricTable.RemoveFabricDelegate(s)
		s.mFabricTable = nil
	}
	s.mState = kNotReady
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

func (s *SessionManager) PrepareMessage(session *SessionHandle, payloadHeader *raw.PayloadHeader, msgBuf *system.PacketBufferHandle, encryptedMessage *EncryptedPacketBufferHandle) error {
	return nil
}

func (s *SessionManager) SendPreparedMessage(session *SessionHandle, preparedMessage *EncryptedPacketBufferHandle) error {
	return nil
}

func (s *SessionManager) AllocateSession(sessionType session.SecureSessionType, sessionEvictionHint *lib.ScopedNodeId) *SessionHandle {
	return nil
}

func (s *SessionManager) ExpireAllSessions(node *lib.ScopedNodeId) {

}
func (s *SessionManager) ExpireAllSessionsForFabric(fabricIndex lib.FabricIndex) {

}

func (s *SessionManager) ExpireAllSessionsOnLogicalFabric(node *lib.ScopedNodeId) {

}

func (s *SessionManager) CreateUnauthenticatedSession(peerAddress netip.AddrPort, config *session.ReliableMessageProtocolConfig) *SessionHandle {
	return nil
}

func (s *SessionManager) FindSecureSessionForNode(nodeId *lib.ScopedNodeId, sessionType session.SecureSessionType) *SessionHandle {
	return nil
}

func (s *SessionManager) OnFabricRemoved(table *credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionManager) FabricTable() *credentials.FabricTable {
	return s.mFabricTable
}
