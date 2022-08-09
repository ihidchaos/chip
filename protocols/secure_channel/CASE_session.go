package secure_channel

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/protocols"
	"github.com/galenliu/gateway/pkg/log"
)

// SessionEstablishmentDelegate : CASEServer implementation
type SessionEstablishmentDelegate interface {
	OnSessionEstablishmentError()
	OnSessionEstablishmentStarted()
	OnSessionEstablished()
}

type CASESessionBase interface {
	messageing.UnsolicitedMessageHandler
	messageing.ExchangeDelegate
	credentials.FabricTableDelegate
	PairingSession
	GetPeerSessionId() uint16
	GetPeer() lib.ScopedNodeId
	GetLocalScopedNodeId() lib.ScopedNodeId
}

// CASESession impl
// UnsolicitedMessageHandler,
// ExchangeDelegate,
// FabricTable::Delegate,
// PairingSession
type CASESession struct {
	*PairingSessionImpl
	mCommissioningHash crypto.HashSha256Stream
	mRemotePubKey      crypto.P256PublicKey
	mEphemeralKey      crypto.P256Keypair
	mSharedSecret      crypto.P256ECDHDerivedSecret
	mValidContext      credentials.ValidationContext
	mGroupDataProvider credentials.GroupDataProvider

	mMessageDigest []byte
	mIPK           []byte

	mSessionResumptionStorage SessionResumptionStorage

	mFabricsTable        *credentials.FabricTable
	mFabricIndex         lib.FabricIndex
	mPeerNodeId          lib.NodeId
	mLocalNodeId         lib.NodeId
	mPeerCATs            lib.CATValues
	mSecureSessionHolder transport.SessionHolderWithDelegate

	mInitiatorRandom []byte

	mResumeResumptionId []byte
	mNewResumptionId    []byte

	mState uint8
}

func NewCASESession() *CASESession {
	return &CASESession{
		PairingSessionImpl:   NewPairingSessionImpl(),
		mSecureSessionHolder: transport.NewSessionHolderWithDelegateImpl(),
	}
}

func (s *CASESession) OnMessageReceived(context *messageing.ExchangeContext, header *raw.PayloadHeader, data []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) GetPeer() lib.ScopedNodeId {
	return lib.ScopedNodeId{
		NodeId:      s.mPeerNodeId,
		FabricIndex: s.GetFabricIndex(),
	}
}

func (s *CASESession) GetLocalScopedNodeId() lib.ScopedNodeId {
	return lib.ScopedNodeId{
		NodeId:      s.mLocalNodeId,
		FabricIndex: s.GetFabricIndex(),
	}
}

func (s *CASESession) OnUnsolicitedMessageReceived(header *raw.PayloadHeader, delegate messageing.ExchangeDelegate) error {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) OnExchangeCreationFailed(delegate messageing.ExchangeDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) OnSessionEstablishmentError() {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) OnSessionEstablishmentStarted() {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) OnSessionEstablished() {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) RegisterUnsolicitedMessageHandlerForProtocol(protocolId *protocols.Id, handler messageing.UnsolicitedMessageHandler) error {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) RegisterUnsolicitedMessageHandlerForType(protocolId *protocols.Id, msgType uint8, handler messageing.UnsolicitedMessageHandler) error {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) Init(
	manger transport.SessionManager,
	policy credentials.CertificateValidityPolicy,
	delegate *CASEServer,
	previouslyEstablishedPeer *lib.ScopedNodeId,
) error {
	s.Clear()

	return nil
}

func (s *CASESession) OnResponseTimeout(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) OnExchangeClosing(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) FabricWillBeRemoved(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) OnFabricRemoved(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) OnFabricCommitted(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) OnFabricUpdated(table credentials.FabricTable, index lib.FabricIndex) {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) SetGroupDataProvider(provider credentials.GroupDataProvider) {
	s.mGroupDataProvider = provider
}

func (s *CASESession) GetFabricIndex() lib.FabricIndex {
	return s.mFabricIndex
}

func (s *CASESession) Clear() {

}

func (s *CASESession) PrepareForSessionEstablishment(
	sessionManger transport.SessionManager,
	fabrics *credentials.FabricTable,
	storage SessionResumptionStorage,
	policy credentials.CertificateValidityPolicy,
	delegate *CASEServer,
	previouslyEstablishedPeer *lib.ScopedNodeId,
	config *transport.ReliableMessageProtocolConfig,
) error {
	err := s.Init(sessionManger, policy, delegate, previouslyEstablishedPeer)
	if err != nil {
		return err
	}
	s.mFabricsTable = fabrics
	s.mRole = transport.KSessionRoleResponder
	s.mSessionResumptionStorage = storage
	s.mLocalMRPConfig = config

	log.Infof("Allocated SecureSessionBase-waiting for Sigma1 msg")
	s.mSecureSessionHolder.Get()

	return nil
}

func (s *CASESession) CopySecureSession() transport.SessionHandle {
	return nil
}
