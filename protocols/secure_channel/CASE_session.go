package secure_channel

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/buffer"
	"github.com/galenliu/chip/lib/tlv"
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

	mFabricsTable *credentials.FabricTable
	mFabricIndex  lib.FabricIndex
	mPeerNodeId   lib.NodeId
	mLocalNodeId  lib.NodeId
	mPeerCATs     lib.CATValues

	mInitiatorRandom []byte

	mResumeResumptionId []byte
	mNewResumptionId    []byte

	mState uint8
}

func NewCASESession() *CASESession {
	return &CASESession{
		PairingSessionImpl: NewPairingSessionImpl(),
	}
}

func (s *CASESession) OnMessageReceived(context *messageing.ExchangeContext, payloadHeader *raw.PayloadHeader, buf *buffer.PacketBuffer) error {
	msgType := messageing.MsgType(payloadHeader.GetMessageType())

	switch s.mState {
	case SInitialized:
		if msgType == messageing.CASESigma1 {
			return s.HandleSigma1AndSendSigma2(buf)
		}
	case SSentSigma1:
	case SSentSigma1Resume:
	case SSentSigma2:
	case SSentSigma3:
	case kSentSigma2Resume:
	default:
		return lib.ChipErrorInvalidMessageType
	}

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

func (s *CASESession) HandleSigma1AndSendSigma2(buf *buffer.PacketBuffer) error {
	return s.HandleSigma1(buf)
}

func (s *CASESession) HandleSigma1(buf *buffer.PacketBuffer) error {

	tlvReader := tlv.NewReader(buf)
	sigma1, err := ParseSigma1(tlvReader)
	if err != nil {
		return err
	}

	log.Infof("Peer assigned session key ID %d", sigma1.initiatorSessionId)
	s.mPeerSessionId = sigma1.initiatorSessionId

	if s.mFabricsTable == nil {
		return lib.ChipErrorIncorrectState
	}

	err = s.FindLocalNodeFromDestinationId(sigma1.destinationId, sigma1.initiatorRandom)
	if err == nil {
		log.Infof("CASE matched destination ID: fabricIndex %x, NodeID 0x", s.mFabricIndex, s.mLocalNodeId)
	} else {
		log.Infof("CASE failed to match destination ID with local fabrics")
	}

	return nil
}

func (s *CASESession) SendSigma2() {

}

func (s *CASESession) FindLocalNodeFromDestinationId(destinationId []byte, initiatorRandom []byte) error {

	for _, fabricInfo := range s.mFabricsTable.GetFabrics() {

		fabriceId := fabricInfo.GetFabricId()
		nodeId := fabricInfo.GetNodeId()
		rootPubKey, err := s.mFabricsTable.FetchRootPubkey(fabricInfo.GetFabricIndex())
		if err != nil {
			return err
		}

		ipKeySet, err := s.mGroupDataProvider.GetIpkKeySet(fabricInfo.GetFabricIndex())
		if err != nil || (ipKeySet.NumKeysUsed == 0 || ipKeySet.NumKeysUsed > credentials.KEpochKeysMax) {
			continue
		}

		// Try every IPK candidate we have for a match
		for keyIdx := 0; keyIdx < ipKeySet.NumKeysUsed; keyIdx++ {
			GenerateCaseDestinationId(ipKeySet.EpochKeys[keyIdx].Key, initiatorRandom, rootPubKey, fabriceId, nodeId)

		}

	}
	return nil
}
