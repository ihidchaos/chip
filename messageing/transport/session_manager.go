package transport

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/pkg/storage"
	"github.com/galenliu/chip/platform/system"
	log "github.com/sirupsen/logrus"
	"net/netip"
)

// SessionMessageDelegate 这里的delegate实例为ExchangeManager
type SessionMessageDelegate interface {
	OnMessageReceived(packetHeader *raw.PacketHeader, payloadHeader *raw.PayloadHeader, session *SessionHandle, duplicate uint8, buf *raw.PacketBuffer)
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

// SessionManagerBase The delegate for TransportManager and FabricTable
// TransportBaseDelegate is the indirect delegate for TransportManager
type SessionManagerBase interface {
	credentials.FabricTableDelegate
	ManagerDelegate
	// SecureGroupMessageDispatch  handle the Secure Group messages
	SecureGroupMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *raw.PacketBuffer)
	// SecureUnicastMessageDispatch  handle the unsecure messages
	SecureUnicastMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *raw.PacketBuffer)
	// UnauthenticatedMessageDispatch handle the unauthenticated(未经认证的) messages
	UnauthenticatedMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *raw.PacketBuffer)

	SetMessageDelegate(SessionMessageDelegate)
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

func (s *SessionManager) OnMessageReceived(srcAddr netip.AddrPort, packet *raw.PacketBuffer) {
	packetHeader, err := packet.PacketHeader()
	if err != nil {
		log.Printf("failed to decode packet header: %s", err.Error())
		return
	}
	if packetHeader.IsEncrypted() {
		if packetHeader.IsGroupSession() {
			s.SecureGroupMessageDispatch(packetHeader, srcAddr, packet)
		} else {
			s.SecureUnicastMessageDispatch(packetHeader, srcAddr, packet)
		}
	} else {
		s.UnauthenticatedMessageDispatch(packetHeader, srcAddr, packet)
	}
}

// UnauthenticatedMessageDispatch 处理没有加密码的消息
func (s *SessionManager) UnauthenticatedMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, buf *raw.PacketBuffer) {

	source := header.SourceNodeId
	destination := header.DestinationNodeId

	if (source.HasValue() && destination.HasValue()) || (!source.HasValue() && !destination.HasValue()) {
		log.Infof("received malformed unsecure packet with source %d destination %d", source, destination)
		return
		//ephemeral node id is only assigned to the initiator, there should be one and only one node id exists.
	}

	var unsecuredSession *UnauthenticatedSession
	if source.HasValue() {
		// Assume peer is the initiator, we are the responder.
		// 对方是发起人，我们是响应者
		unsecuredSession = s.mUnauthenticatedSessions.FindOrAllocateResponder(source, GetLocalMRPConfig())
		if unsecuredSession == nil {
			log.Infof("UnauthenticatedSession exhausted")
			return
		}
	} else {
		// Assume peer is the responder, we are the initiator.
		// 对方为响应，我们是发起人
		unsecuredSession = s.mUnauthenticatedSessions.FindInitiator(destination)
		if unsecuredSession == nil {
			log.Infof("Received unknown unsecure packet for initiator %d", destination)
			return
		}
	}

	unsecuredSession.SetPeerAddress(addr)
	isDuplicate := DuplicateMessageNo
	// 更新Session
	unsecuredSession.MarkActiveRx()

	var payloadHeader = raw.NewPayloadHeader()
	err := payloadHeader.DecodeAndConsume(buf)
	if err != nil {
		log.Warnf("Received invaild packet")
		return
	}

	err = unsecuredSession.GetPeerMessageCounter().VerifyUnencrypted(header.MessageCounter)
	if err != nil && err == lib.MatterErrorDuplicateMessageReceived {
		isDuplicate = DuplicateMessageYes
		log.Infof(
			"Received a duplicate message with MessageCounter: %v on exchange %v",
			header.MessageCounter, payloadHeader)
	} else {
		unsecuredSession.GetPeerMessageCounter().CommitUnencrypted(header.MessageCounter)
	}
	if s.mCB != nil {
		s.mCB.OnMessageReceived(header, payloadHeader, NewSessionHandle(unsecuredSession), isDuplicate, buf)
	}
}

func (s *SessionManager) FabricWillBeRemoved(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionManager) OnFabricRemoved(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionManager) SystemLayer() system.Layer {
	return s.mSystemLayer
}

func (s *SessionManager) OnFabricCommitted(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionManager) OnFabricUpdated(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

// SecureGroupMessageDispatch 处理加密的组播消息
func (s *SessionManager) SecureGroupMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, msg *raw.PacketBuffer) {

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
		log.Errorf("SessionManager.FabricRemoved err:%s", err.Error())
	}
}

// SecureUnicastMessageDispatch 处理分支，加密的单播消息
func (s *SessionManager) SecureUnicastMessageDispatch(header *raw.PacketHeader, addr netip.AddrPort, msg *raw.PacketBuffer) {
	secureSession := s.mSecureSessions.FindSecureSessionByLocalKey(header.SessionId)
	_ = DuplicateMessageNo
	if msg.IsNull() {
		log.Infof("Secure transport received Unicast NULL packet, discarding")
		return
	}
	if secureSession == nil {
		log.Infof("Data received on an unknown session (LSID=%d). Dropping it!", header.SessionId)
		return
	}

	if secureSession.IsDefunct() && !secureSession.IsActiveSession() && !secureSession.IsPendingEviction() {
		log.Infof("Secure transport received message on a session in an invalid state (state = '%s')",
			secureSession.GetStateStr())
	}
	var nodeId = lib.UndefinedNodeId
	if secureSession.GetSecureSessionType() == CASE {
		nodeId = secureSession.GetPeerNodeId()
	}
	nonce, _ := BuildNonce(header.SecFlags, header.MessageCounter, nodeId)
	_ = Decrypt(secureSession.GetCryptoContext(), nonce, header, msg)
}
