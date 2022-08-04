package secure_channel

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	transport2 "github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/gateway/pkg/log"
)

// SessionEstablishmentDelegate : CASEServer implementation
type SessionEstablishmentDelegate interface {
	OnSessionEstablishmentError()
	OnSessionEstablishmentStarted()
	OnSessionEstablished()
}

type CASESession struct {
	PairingSession
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
	mSecureSessionHolder transport2.SessionHolderWithDelegate

	mInitiatorRandom []byte

	mResumeResumptionId []byte
	mNewResumptionId    []byte

	mState uint8
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

func (s *CASESession) Init(
	manger transport2.SessionManager,
	policy credentials.CertificateValidityPolicy,
	delegate *CASEServer,
	previouslyEstablishedPeer *lib.ScopedNodeId,
) error {
	s.Clear()

	return nil
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
	sessionManger transport2.SessionManager,
	fabrics *credentials.FabricTable,
	storage SessionResumptionStorage,
	policy credentials.CertificateValidityPolicy,
	delegate *CASEServer,
	previouslyEstablishedPeer *lib.ScopedNodeId,
	config *transport2.ReliableMessageProtocolConfig,
) error {
	err := s.Init(sessionManger, policy, delegate, previouslyEstablishedPeer)
	if err != nil {
		return err
	}
	s.mFabricsTable = fabrics
	s.mRole = transport2.KSessionRoleResponder
	s.mSessionResumptionStorage = storage
	s.mLocalMRPConfig = config

	log.Info("Allocated SecureSession (%s) - waiting for Sigma1 msg")
	s.mSecureSessionHolder.Get()

	return nil
}

func (s *CASESession) CopySecureSession() transport2.SessionHandle {
	return nil
}
