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
	log.Infof("CASE Session HandleSigma1")

	return nil
}

func (s *CASESession) ParseSigma1(buf *buffer.PacketBuffer) (initiatorRandom, destinationId, initiatorEphPubKey []byte, sessionId uint16, err error) {
	var kInitiatorRandomTag uint8 = 1
	var kInitiatorSessionIdTag uint8 = 2
	var kDestinationIdTag uint8 = 3
	var kInitiatorPubKeyTag uint8 = 4
	//var kInitiatorMRPParamsTag uint8 = 5
	//var kResumptionIDTag uint8 = 6
	//var kResume1MICTag uint8 = 7

	// Sigma1，这里应该读取到Structure 0x15

	tlvReader := tlv.NewReader(buf)
	err = tlvReader.NextE(tlv.AnonymousTag(), tlv.TypeStructure)
	if err != nil {
		return
	}
	// Sigma1，Tag = 1 initiatorRandom  20个字节的随机数
	err = tlvReader.NextE(tlv.ContextTag(kInitiatorRandomTag))
	initiatorRandom, err = tlvReader.GetBytesView()
	if err != nil {

		return
	}

	//Sigma1， Tag =2 Session id
	err = tlvReader.NextE(tlv.ContextTag(kInitiatorSessionIdTag))
	if err != nil {

		return
	}
	sessionId, err = tlvReader.GetUint16()

	//Sigma1，Tag=3	destination id 20个字节的认证码
	err = tlvReader.NextE(tlv.ContextTag(kDestinationIdTag))
	destinationId, err = tlvReader.GetBytesView()
	if err != nil {

		return
	}

	//Sigma1，Tag=4	 Initiator PubKey 1个字节的公钥
	err = tlvReader.NextE(tlv.ContextTag(kInitiatorPubKeyTag))
	initiatorEphPubKey, err = tlvReader.GetBytesView()
	if err != nil {

		return
	}

	return
}
