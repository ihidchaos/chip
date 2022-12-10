package session

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport/raw"
	log "golang.org/x/exp/slog"
	"sync"
	"time"
)

type OptionalFunc func(*Secure)

type SecureDelegate interface {
	ReleaseSession(secure *Secure)
}

type Secure struct {
	*BaseImpl
	mState                 State
	mDelegate              SecureDelegate
	mSecureSessionType     SecureType
	mLocalSessionId        uint16
	mPeerSessionId         uint16
	mLocalNodeId           lib.NodeId
	mPeerNodeId            lib.NodeId
	mRemoteMRPConfig       *messageing.ReliableMessageProtocolConfig
	mCryptoContext         *CryptoContext
	mPeerCATs              *lib.CATValues
	mLastPeerActivityTime  time.Time
	mLastActivityTime      time.Time
	mSessionMessageCounter MessageCounter //
	mPeerAddress           raw.PeerAddress
}

func NewSecure(
	table SecureDelegate,
	secureSessionType SecureType,
	localSessionId uint16,
	options ...OptionalFunc,
) *Secure {
	session := &Secure{
		mDelegate:          table,
		mState:             kEstablishing,
		mSecureSessionType: secureSessionType,
		mLocalSessionId:    localSessionId,
	}
	session.BaseImpl = &BaseImpl{
		locker:           sync.Mutex{},
		mFabricIndex:     lib.UndefinedFabricIndex(),
		mHolders:         nil,
		mSessionType:     kSecure,
		mPeerAddress:     raw.PeerAddress{},
		base:             session,
		ReferenceCounted: lib.NewReferenceCounted(0, session),
	}
	for _, option := range options {
		option(session)
	}
	session.MoveToState(kActive)
	return session
}

func (s *Secure) IsActive() bool {
	return s.mState == kActive
}

func (s *Secure) IsEstablishing() bool {
	return s.mState == kEstablishing
}

func (s *Secure) ComputeRoundTripTimeout(upperlayerProcessingTimeout time.Duration) time.Duration {
	if s.IsGroup() {
		return time.Duration(0)
	}
	return s.AckTimeout() + upperlayerProcessingTimeout
}

func (s *Secure) Released() {
	s.mDelegate.ReleaseSession(s)
}

func (s *Secure) MoveToState(targetState State) {
	if s.mState != targetState {
		log.Info("Moving state", "from", s.mState, "to", targetState, "Tag", s)
		s.mState = targetState
	}
}

func (s *Secure) ClearValue() {
	s.Released()
}

func (s *Secure) LocalSessionId() uint16 {
	return s.mLocalSessionId
}

func (s *Secure) IsDefunct() bool {
	return s.mState == kDefunct
}

func (s *Secure) IsPendingEviction() bool {
	return s.mState == kPendingEviction
}

func (s *Secure) State() State {
	return s.mState
}

func (s *Secure) SecureType() SecureType {
	return s.mSecureSessionType
}

func (s *Secure) GetCryptoContext() *CryptoContext {
	return s.mCryptoContext
}

func (s *Secure) PeerNodeId() lib.NodeId {
	return s.mPeerNodeId
}

func (s *Secure) SessionMessageCounter() MessageCounter {
	return s.mSessionMessageCounter
}

func (s *Secure) MarkActiveRx() {
	s.mLastPeerActivityTime = time.Now()
	s.MarkActive()
	if s.mState == kDefunct {
		s.MoveToState(kActive)
	}
}

func (s *Secure) MarkActive() {
	s.mLastActivityTime = time.Now()
}

func (s *Secure) PeerAddress() raw.PeerAddress {
	return s.mPeerAddress
}

func (s *Secure) SetPeerAddress(address raw.PeerAddress) {
	s.mPeerAddress = address
}

func (s *Secure) GetPeer() lib.ScopedNodeId {
	return lib.NewScopedNodeId(s.mPeerNodeId, s.FabricIndex())
}

func (s *Secure) RemoteMRPConfig() *messageing.ReliableMessageProtocolConfig {
	return s.mRemoteMRPConfig
}

func (s *Secure) AckTimeout() time.Duration {
	switch s.BaseImpl.mPeerAddress.TransportType {
	case raw.Udp:
		return messageing.GetRetransmissionTimeout(s.mRemoteMRPConfig.ActiveRetransTimeout,
			s.mRemoteMRPConfig.IdleRetransTimeout, s.mLastActivityTime, 0)
	case raw.Tcp:
		return 30 * time.Second
	case raw.Ble:
		return 5 * time.Second
	default:
		return time.Duration(0)
	}
}

func (s *Secure) Activate(local, peer lib.ScopedNodeId, ts lib.CATValues, peerSessionId uint16, config *messageing.ReliableMessageProtocolConfig) {

}

func WithLocalNodeId(localNodeId lib.NodeId) OptionalFunc {
	return func(session *Secure) {
		session.mLocalNodeId = localNodeId
	}
}

func WithPeerNodeId(peerNodeId lib.NodeId) OptionalFunc {
	return func(session *Secure) {
		session.mPeerNodeId = peerNodeId
	}
}

func WithPeerCATs(peerCATs *lib.CATValues) OptionalFunc {
	return func(session *Secure) {
		session.mPeerCATs = peerCATs
	}
}

func WithPeerSessionId(peerSessionId uint16) OptionalFunc {
	return func(session *Secure) {
		session.mPeerSessionId = peerSessionId
	}
}

func WithFabricIndex(index lib.FabricIndex) OptionalFunc {
	return func(session *Secure) {
		session.mFabricIndex = index
	}
}

func WithMRPC(config *messageing.ReliableMessageProtocolConfig) OptionalFunc {
	return func(session *Secure) {
		session.mRemoteMRPConfig = config
	}
}
