package transport

import (
	"github.com/galenliu/chip/lib"
	log "golang.org/x/exp/slog"
	"net/netip"
	"time"
)

type SecureSessionBase interface {
	Session
	SecureSessionType() SecureSessionType
	MarkActiveRx()
	MarkActive()
	PeerAddress() netip.AddrPort
	SetPeerAddress(address netip.AddrPort)
}

type OptionalFunc func(*SecureSession)

type SecureSession struct {
	*SessionBaseImpl
	mState                 SecureSessionState
	mTable                 *SecureSessionTable
	mSecureSessionType     SecureSessionType
	mLocalSessionId        uint16
	mPeerSessionId         uint16
	mLocalNodeId           lib.NodeId
	mPeerNodeId            lib.NodeId
	mRemoteMRPConfig       *ReliableMessageProtocolConfig
	mCryptoContext         *CryptoContext
	mPeerCATs              *lib.CATValues
	mLastPeerActivityTime  time.Time
	mLastActivityTime      time.Time
	mSessionMessageCounter *SessionMessageCounter
	mPeerAddress           netip.AddrPort
}

func NewSecureSession(
	table *SecureSessionTable,
	secureSessionType SecureSessionType,
	localSessionId uint16,
	options ...OptionalFunc,
) *SecureSession {
	session := &SecureSession{
		mTable:             table,
		mState:             kEstablishing,
		mSecureSessionType: secureSessionType,
		mLocalSessionId:    localSessionId,
	}
	session.SessionBaseImpl = NewSessionBaseImpl(0, kSecure, session)
	for _, option := range options {
		option(session)
	}
	session.MoveToState(kActive)
	return session
}

func (s *SecureSession) IsActiveSession() bool {
	return s.mState == kActive
}

func (s *SecureSession) IsGroupSession() bool {
	return s.SessionType() == kSecure
}

func (s *SecureSession) IsEstablishing() bool {
	return s.mState == kEstablishing
}

func (s *SecureSession) DispatchSessionEvent(delegate SessionDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSession) ComputeRoundTripTimeout(duration time.Duration) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s *SecureSession) Released() {
	s.mTable.ReleaseSession(s)
}

func (s *SecureSession) MoveToState(targetState SecureSessionState) {
	if s.mState != targetState {
		log.Info("Moving state", "from", s.mState, "to", targetState, "Tag", s)
		s.mState = targetState
	}
}

func (s *SecureSession) ClearValue() {
	s.Released()
}

func (s *SecureSession) LocalSessionId() uint16 {
	return s.mLocalSessionId
}

func (s *SecureSession) IsDefunct() bool {
	return s.mState == kDefunct
}

func (s *SecureSession) IsPendingEviction() bool {
	return s.mState == kPendingEviction
}

func (s *SecureSession) State() SecureSessionState {
	return s.mState
}

func (s *SecureSession) SecureSessionType() SecureSessionType {
	return s.mSecureSessionType
}

func (s *SecureSession) GetCryptoContext() *CryptoContext {
	return s.mCryptoContext
}

func (s *SecureSession) GetPeerNodeId() lib.NodeId {
	return s.mPeerNodeId
}

func (s *SecureSession) LogValue() log.Value {
	return log.GroupValue(
		log.String("name", "SecureSession"),
	)
}

func (s *SecureSession) SessionMessageCounter() *SessionMessageCounter {
	return s.mSessionMessageCounter
}

func (s *SecureSession) MarkActiveRx() {
	s.mLastPeerActivityTime = time.Now()
	s.MarkActive()
	if s.mState == kDefunct {
		s.MoveToState(kActive)
	}
}

func (s *SecureSession) MarkActive() {
	s.mLastActivityTime = time.Now()
}

func (s *SecureSession) PeerAddress() netip.AddrPort {
	return s.mPeerAddress
}

func (s *SecureSession) SetPeerAddress(address netip.AddrPort) {
	s.mPeerAddress = address
}

func (s *SecureSession) GetPeer() lib.ScopedNodeId {
	return lib.NewScopedNodeId(s.mPeerNodeId, s.FabricIndex())
}

func WithLocalNodeId(localNodeId lib.NodeId) OptionalFunc {
	return func(session *SecureSession) {
		session.mLocalNodeId = localNodeId
	}
}

func WithPeerNodeId(peerNodeId lib.NodeId) OptionalFunc {
	return func(session *SecureSession) {
		session.mPeerNodeId = peerNodeId
	}
}

func WithPeerCATs(peerCATs *lib.CATValues) OptionalFunc {
	return func(session *SecureSession) {
		session.mPeerCATs = peerCATs
	}
}

func WithPeerSessionId(peerSessionId uint16) OptionalFunc {
	return func(session *SecureSession) {
		session.mPeerSessionId = peerSessionId
	}
}

func WithFabricIndex(index lib.FabricIndex) OptionalFunc {
	return func(session *SecureSession) {
		session.mFabricIndex = index
	}
}

func WithMRPC(config *ReliableMessageProtocolConfig) OptionalFunc {
	return func(session *SecureSession) {
		session.mRemoteMRPConfig = config
	}
}
