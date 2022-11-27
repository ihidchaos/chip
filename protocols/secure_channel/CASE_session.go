package secure_channel

import (
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/pkg/tlv"
	"github.com/galenliu/chip/platform/system"
	"golang.org/x/exp/rand"
	log "golang.org/x/exp/slog"
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
	PairingSessionBase
	GetPeerSessionId() uint16
	GetPeer() lib.ScopedNodeId
	GetLocalScopedNodeId() lib.ScopedNodeId
}

// CASESession
// UnsolicitedMessageHandler,
// ExchangeDelegate,
// FabricTable::Delegate,
// PairingSessionBase
type CASESession struct {
	*PairingSession
	mCommissioningHash *crypto.HashSha256Stream

	mRemotePubKey      *crypto.P256PublicKey
	mEphemeralKey      *crypto.P256Keypair
	mSharedSecret      []byte //crypto.P256ECDHDerivedSecret
	mValidContext      credentials.ValidationContext
	mGroupDataProvider *credentials.GroupDataProvider

	mMessageDigest []byte
	mIPK           []byte

	mSessionResumptionStorage SessionResumptionStorage

	mFabricsTable *credentials.FabricTable
	mFabricIndex  lib.FabricIndex
	mPeerNodeId   lib.NodeId
	mLocalNodeId  lib.NodeId
	mPeerCATs     lib.CATValues

	mInitiatorRandom []byte

	mResumeResumptionId []byte //会话恢复ID 16个字节
	mNewResumptionId    []byte //会话恢复ID 16个字节

	mState state
}

func NewCASESession() *CASESession {
	return &CASESession{
		PairingSession: NewPairingSessionImpl(),
	}
}

func (s *CASESession) GetMessageDispatch() messageing.ExchangeMessageDispatchBase {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) OnMessageReceived(ec *messageing.ExchangeContext, payloadHeader *raw.PayloadHeader, msg *system.PacketBufferHandle) error {

	err := s.ValidateReceivedMessage(ec, payloadHeader, msg)
	if err != nil {
		return err
	}
	msgType := MsgType(payloadHeader.MessageType())

	sha256.New()

	switch s.mState {
	case initialized:
		if msgType == CASE_Sigma1 {
			return s.HandleSigma1AndSendSigma2(msg)
		}
	case sentSigma1:
	case sentSigma1Resume:
	case sentSigma2:
	case sentSigma3:
	case sentSigma2Resume:
	default:
		return lib.InvalidMessageType
	}
	return nil
}

func (s *CASESession) HandleSigma1(msg *system.PacketBufferHandle) error {

	log.Debug("CASESession HandleSigma1")

	s.mCommissioningHash.AddData(msg.Bytes())

	var sessionResumptionRequested = false
	tlvReader := tlv.NewReader(msg)
	sigma1, err := ParseSigma1(tlvReader, sessionResumptionRequested)
	if err != nil {
		return err
	}

	log.Info("Peer assigned session key ID %d", sigma1.initiatorSessionId)
	s.mPeerSessionId = sigma1.initiatorSessionId

	if s.mFabricsTable == nil {
		return lib.IncorrectState
	}

	if sigma1.sessionResumptionRequested && len(sigma1.resumptionId) == 16 {
	}

	err = s.FindLocalNodeFromDestinationId(sigma1.destinationId, sigma1.initiatorRandom)
	if err == nil {
		log.Info("CASE matched destination", "FabricIndex", s.mFabricIndex, "NodeID", s.mLocalNodeId)
	} else {
		log.Info("CASE failed to match destination ID with local fabrics")
	}

	pubKey, err := crypto.UnmarshalP256PublicKey(sigma1.initiatorEphPubKey)
	if err != nil {
		return err
	}
	s.mRemotePubKey = &pubKey

	err = s.SendSigma2()
	if err != nil {
		//s.SendStatusReport()
		s.mState = initialized
		return err
	}
	s.mDelegate.OnSessionEstablishmentStarted()
	return nil
}

func (s *CASESession) SendSigma2() error {
	sessionId, err := s.LocalSessionId()
	if err != nil {
		return err
	}
	//获取ICA证书
	icaCert, err := s.mFabricsTable.FetchICACert(s.mFabricIndex)
	if err != nil && len(icaCert) != credentials.MaxCHIPCertLength {
		return fmt.Errorf("sigma2 icaCert err")
	}
	//获取节点操作证书(NOC)
	nocCert, err := s.mFabricsTable.FetchNOCCert(s.mFabricIndex)
	if err != nil && len(nocCert) != credentials.MaxCHIPCertLength {
		return fmt.Errorf("sigma2 nocCert err")
	}

	//生成一个32字节的随机数
	msgRand := make([]byte, sigmaParamRandomNumberSize)
	_, err = rand.Read(msgRand)

	//生成临时密钥对
	s.mEphemeralKey = s.mFabricsTable.AllocateEphemeralKeypairForCASE()
	if s.mEphemeralKey == nil {
		return fmt.Errorf("sigma2 ephemeral key error")
	}

	// 生成共享密钥
	s.mSharedSecret = s.mEphemeralKey.ECDHDeriveSecret(*s.mRemotePubKey)

	//生成一个盐值
	salt, err := s.ConstructSaltSigma2(msgRand, elliptic.Marshal(elliptic.P256(), s.mEphemeralKey.PublicKey.X, s.mEphemeralKey.Y), s.mIPK)
	if err != nil || salt == nil {
		return fmt.Errorf("sigma2 salt error")
	}

	// 使用HKDF函数派生出密钥sr2k,sr2k为16个字节
	sr2k := crypto.HKDFSha256(s.mSharedSecret, salt, KDFSR2Info)
	if sr2k == nil || len(sr2k) != crypto.SymmetricKeyLengthBytes {
		return fmt.Errorf("sigma2 hkdfSha256 error")
	}

	//msgR2SignedLen := tlv.EstimateStructOverhead(credentials.MaxCHIPCertLength, credentials.MaxCHIPCertLength, crypto.P256PublicKeyLength, crypto.P256PublicKeyLength)

	//msgR2Signed := s.ConstructTBSData(nocCert, icaCert, s.mEphemeralKey.MarshalPublicKey(), s.mRemotePubKey.Marshal())

	tbsData2Signature := s.mFabricsTable.SignWithOpKeypair(s.mFabricIndex).Bytes()

	tlvWriterMsg1 := tlv.NewWriter()
	_, err = tlvWriterMsg1.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure)
	if err != nil {
		return err
	}

	err = tlvWriterMsg1.PutBytes(tlv.ContextTag(TagTBEDataSenderNOC), nocCert)
	if err != nil {
		return err
	}
	if len(icaCert) > 0 {
		err = tlvWriterMsg1.PutBytes(tlv.ContextTag(TagTBEDataSenderICAC), icaCert)
		if err != nil {
			return err
		}
	}
	err = tlvWriterMsg1.PutBytes(tlv.ContextTag(TagTBEDataSignature), tbsData2Signature)
	if err != nil {
		return err
	}

	s.mNewResumptionId = make([]byte, resumptionIdSize)
	err = crypto.DRBGBytes(s.mNewResumptionId)
	if err != nil {
		return err
	}
	err = tlvWriterMsg1.PutBytes(tlv.ContextTag(TagTBEDataResumptionID), s.mNewResumptionId)
	if err != nil {
		return err
	}

	err = tlvWriterMsg1.EndContainer(tlv.TypeStructure)
	if err != nil {
		return err
	}

	// 使用对称密钥sr2k 对 Sigma2数量进行加密
	msgR2Encrypted, err := crypto.AesCcmEncrypt(tlvWriterMsg1.Bytes(), sr2k, kTBEData2Nonce, crypto.AEADMicLengthBytes)

	tlvWriterMsg2 := tlv.NewWriter()
	_, err = tlvWriterMsg2.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure)
	if err != nil {
		return err
	}
	err = tlvWriterMsg2.PutBytes(tlv.ContextTag(1), msgRand)
	if err != nil {
		return err
	}
	err = tlvWriterMsg2.PutUint(tlv.ContextTag(2), uint64(sessionId))
	if err != nil {
		return err
	}
	err = tlvWriterMsg2.PutBytes(tlv.ContextTag(3), s.mEphemeralKey.PubBytes())
	if err != nil {
		return err
	}
	err = tlvWriterMsg2.PutBytes(tlv.ContextTag(4), msgR2Encrypted)
	if err != nil {
		return err
	}

	if s.mLocalMRPConfig != nil {
		_, err = tlvWriterMsg2.StartContainer(tlv.ContextTag(5), tlv.TypeStructure)
		if err != nil {
			return err
		}
		err = tlvWriterMsg2.PutUint(tlv.ContextTag(1), uint64(s.mLocalMRPConfig.IdleRetransTimeout.Milliseconds()))
		if err != nil {
			return err
		}
		err = tlvWriterMsg2.PutUint(tlv.ContextTag(2), uint64(s.mLocalMRPConfig.ActiveRetransTimeout.Milliseconds()))
		if err != nil {
			return err
		}
		err = tlvWriterMsg2.EndContainer(tlv.TypeStructure)
		if err != nil {
			return err
		}
	}
	err = tlvWriterMsg2.EndContainer(tlv.TypeStructure)
	if err != nil {
		return err
	}

	//记录下Hash值
	s.mCommissioningHash.AddData(tlvWriterMsg2.Bytes())

	err = s.mExchangeCtxt.SendMessage(ProtocolId, uint8(CASE_Sigma2), tlvWriterMsg2.Bytes(), messageing.ExpectResponse)
	if err != nil {
		return err
	}

	s.mState = sentSigma2

	log.Info("Sent sigma2 message")

	return nil
}

func (s *CASESession) GetPeer() lib.ScopedNodeId {
	return lib.NewScopedNodeId(
		s.mPeerNodeId,
		s.FabricIndex(),
	)
}

func (s *CASESession) LocalScopedNodeId() lib.ScopedNodeId {
	return lib.NewScopedNodeId(s.mLocalNodeId, s.FabricIndex())
}

func (s *CASESession) OnUnsolicitedMessageReceived(header *raw.PayloadHeader) (messageing.ExchangeDelegate, error) {
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

func (s *CASESession) Init(
	manger transport.SessionManagerBase,
	policy credentials.CertificateValidityPolicy,
	delegate *CASEServer,
	previouslyEstablishedPeer *lib.ScopedNodeId,
) error {
	s.Clear()
	s.mCommissioningHash.Begin()

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

func (s *CASESession) SetGroupDataProvider(provider *credentials.GroupDataProvider) {
	s.mGroupDataProvider = provider
}

func (s *CASESession) FabricIndex() lib.FabricIndex {
	return s.mFabricIndex
}

func (s *CASESession) Clear() {

}

func (s *CASESession) PrepareForSessionEstablishment(
	sessionManger transport.SessionManagerBase,
	fabrics *credentials.FabricTable,
	storage SessionResumptionStorage,
	policy credentials.CertificateValidityPolicy,
	delegate *CASEServer,
	previouslyEstablishedPeer *lib.ScopedNodeId,
	config *session.ReliableMessageProtocolConfig,
) error {
	err := s.Init(sessionManger, policy, delegate, previouslyEstablishedPeer)
	if err != nil {
		return err
	}
	s.mFabricsTable = fabrics
	s.mRole = session.KSessionRoleResponder
	s.mSessionResumptionStorage = storage
	s.mLocalMRPConfig = config

	log.Info("Allocated SecureSession-waiting for Sigma1 msg")

	return nil
}

func (s *CASESession) CopySecureSession() *transport.SessionHandle {
	return nil
}

func (s *CASESession) HandleSigma1AndSendSigma2(buf *system.PacketBufferHandle) error {
	log.Debug("CASESession HandleSigma1_and_SendSigma2")
	return s.HandleSigma1(buf)
}

func (s *CASESession) FindLocalNodeFromDestinationId(destinationId []byte, initiatorRandom []byte) error {
	for _, fabricInfo := range s.mFabricsTable.Fabrics() {
		fabriceId := fabricInfo.FabricId()
		nodeId := fabricInfo.GetNodeId()
		rootPubKey, err := s.mFabricsTable.FetchRootPubkey(fabricInfo.FabricIndex())
		if err != nil {
			return err
		}
		ipkKeySet, err := s.mGroupDataProvider.GetIpkKeySet(fabricInfo.FabricIndex())
		if err != nil || (ipkKeySet.NumKeysUsed == 0 || ipkKeySet.NumKeysUsed > credentials.KEpochKeysMax) {
			continue
		}

		// Try every IPK candidate we have for a match
		for keyIdx := 0; uint8(keyIdx) < ipkKeySet.NumKeysUsed; keyIdx++ {
			candidateDestinationId, err := GenerateCASEDestinationId(ipkKeySet.EpochKeys[keyIdx].Key[:], initiatorRandom, rootPubKey, fabriceId, nodeId)
			if err != nil && len(candidateDestinationId) == len(destinationId) {
				s.mFabricIndex = fabricInfo.FabricIndex()
				s.mLocalNodeId = nodeId
			}
		}
	}
	return nil
}

func (s *CASESession) ConstructSaltSigma2(rand []byte, publicKey []byte, ipk []byte) (saltSpan []byte, err error) {
	//md := make([]byte, crypto.Sha256HashLength)
	saltSpan = make([]byte, ipkSize+sigmaParamRandomNumberSize+crypto.P256PublicKeyLength+crypto.Sha256HashLength)
	buf := bytes.NewBuffer(saltSpan)
	_, err = buf.Write(ipk)
	_, err = buf.Write(rand)
	_, err = buf.Write(publicKey)
	_, err = buf.Write(s.mCommissioningHash.Bytes())
	return saltSpan, nil
}

func (s *CASESession) ValidateReceivedMessage(ec *messageing.ExchangeContext, header *raw.PayloadHeader, msg *system.PacketBufferHandle) error {
	if s.mExchangeCtxt != nil {
		if s.mExchangeCtxt != ec {
			return lib.MATTER_ERROR_INVALID_ARGUMENT
		}
	} else {
		s.mExchangeCtxt = ec
	}
	s.mExchangeCtxt.UseSuggestedResponseTimeout(kExpectedHighProcessingTime)
	if msg.IsNull() {
		return lib.MATTER_ERROR_INVALID_ARGUMENT
	}
	return nil
}
