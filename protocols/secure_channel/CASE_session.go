package secure_channel

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/galenliu/chip"
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/bitflags"
	"github.com/galenliu/chip/lib/tlv"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/platform/system"
	"github.com/moznion/go-optional"
	"golang.org/x/exp/rand"
	log "golang.org/x/exp/slog"
	"hash"
)

// SessionEstablishmentDelegate : CASEServer implementation
type SessionEstablishmentDelegate interface {
	OnSessionEstablishmentError(err error)
	OnSessionEstablishmentStarted()
	OnSessionEstablished(session *transport.SessionHandle)
}

// CASESession
// UnsolicitedMessageHandler,
// ExchangeDelegate,
// FabricTable::Delegate,
// pairingSessionBase
type CASESession struct {
	pairingSession
	mCommissioningHash hash.Hash

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

	mResumeResumptionId []byte              //会话恢复ID 16个字节
	mNewResumptionId    ResumptionIdStorage //会话恢复ID 16个字节

	mState state
}

func (s *CASESession) PeerCATs() lib.CATValues {
	return s.mPeerCATs
}

func (s *CASESession) DeriveSecureSession(ctx *transport.CryptoContext) error {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) Peer() lib.ScopedNodeId {
	return lib.NewScopedNodeId(s.mPeerNodeId, s.mFabricIndex)
}

func NewCASESession() *CASESession {
	caseSession := &CASESession{}
	caseSession.base = caseSession
	return caseSession
}

func (s *CASESession) GetMessageDispatch() messageing.ExchangeMessageDispatch {
	//TODO implement me
	panic("implement me")
}

func (s *CASESession) OnMessageReceived(ec *messageing.ExchangeContext, payloadHeader *raw.PayloadHeader, msg *system.PacketBufferHandle) error {

	err := s.ValidateReceivedMessage(ec, payloadHeader, msg)
	if err != nil {
		return err
	}
	msgType := MsgType(payloadHeader.MessageType)

	sha256.New()

	switch s.mState {
	case initialized:
		if msgType == CASESigma1 {
			return s.HandleSigma1AndSendSigma2(msg)
		}
	case sentSigma1:
	case sentSigma1Resume:
	case sentSigma2:
	case sentSigma3:
	case sentSigma2Resume:
	default:
		return chip.ErrorInvalidMessageType
	}
	return nil
}

func (s *CASESession) HandleSigma1(msg *system.PacketBufferHandle) error {

	log.Debug("CASESession HandleSigma1")

	s.mCommissioningHash.Write(msg.Bytes())

	tlvDecode := tlv.NewDecoder(msg)
	sigma1, err := ParseSigma1(tlvDecode, false)
	if err != nil {
		return err
	}

	log.Info("Peer assigned session key ID %d", sigma1.initiatorSessionId)
	s.peerSessionId = optional.Option[uint16]{sigma1.initiatorSessionId}

	if s.mFabricsTable == nil {
		return chip.ErrorIncorrectState
	}

	if sigma1.sessionResumptionRequested && len(sigma1.resumptionId) == 16 {
	}

	err = s.FindLocalNodeFromDestinationId(sigma1.destinationId, sigma1.initiatorRandom)
	if err == nil {
		log.Info("CASE matched destination", "FabricIndex", s.mFabricIndex, "NodeID", s.mLocalNodeId)
	} else {
		log.Info("CASE failed to match destination ID with local fabrics")
	}

	pubKey, err := crypto.UnmarshalPublicKey(sigma1.initiatorEphPubKey)
	if err != nil {
		return err
	}
	s.mRemotePubKey = pubKey

	err = s.SendSigma2()
	if err != nil {
		//s.sendStatusReport()
		s.mState = initialized
		return err
	}
	s.delegate.OnSessionEstablishmentStarted()
	return nil
}

func (s *CASESession) SendSigma2() error {

	if s.LocalSessionId().IsNone() {
		return chip.New(chip.ErrorInvalidArgument, "CASESession", "SendSigma2 LocalSessionId is nil")
	}
	if s.mFabricsTable == nil {
		return chip.New(chip.ErrorIncorrectState, "CASESession", "SendSigma2 FabriceTable is nil")
	}

	//获取ICA证书
	icaCert, err := s.mFabricsTable.FetchICACert(s.mFabricIndex)
	if err != nil && len(icaCert) < credentials.MaxCHIPCertLength {
		return chip.New(chip.ErrorInvalidArgument, "CASESession", "sigma2 icaCert err")
	}

	//获取节点操作证书(NOC)
	nocCert, err := s.mFabricsTable.FetchNOCCert(s.mFabricIndex)
	if err != nil && len(nocCert) < credentials.MaxCHIPCertLength {
		return fmt.Errorf("sigma2 nocCert err")
	}

	//生成一个32字节的随机数
	msgRand := make([]byte, sigmaParamRandomNumberSize)
	_, err = rand.Read(msgRand)

	//生成临时密钥对
	s.mEphemeralKey = s.mFabricsTable.AllocateEphemeralKeypairForCASE()

	// 生成共享密钥
	if s.mSharedSecret, err = s.mEphemeralKey.ECDHDeriveSecret(s.mRemotePubKey); err != nil {
		return err
	}

	//生成一个盐值
	salt, err := s.ConstructSaltSigma2(msgRand, s.mEphemeralKey.PrivateKey().PublicKey().Bytes(), s.mIPK)
	if err != nil {
		return chip.New(chip.ErrorInternal, "CASESession", "Construct salt sigma2 err")
	}

	// 使用HKDF函数派生出密钥sr2k,sr2k为16个字节
	sr2k := crypto.HKDFSha256(s.mSharedSecret, salt, KDFSR2Info)
	if sr2k == nil || len(sr2k) != crypto.SymmetricKeyLengthBytes {
		return chip.New(chip.ErrorInternal, "CASESession", "sigma2 hkdfSha256 error")
	}

	// Construct Sigma2 TBS Data
	msgR2Signed, err := s.ConstructTBSData(nocCert, icaCert, s.mEphemeralKey.PrivateKey().PublicKey().Bytes(), s.mRemotePubKey.PublicKey().Bytes())
	if err != nil {
		return err
	}

	// Generate a Signature
	tbsData2Signature, err := s.mFabricsTable.SignWithOpKeypair(s.mFabricIndex, msgR2Signed)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	e := tlv.NewEncoder(buf)
	outType := tlv.TypeNotSpecified
	if outType, err = e.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure); err != nil {
		return err
	}
	if err = e.Put(tlv.ContextTag(TagTBEDataSenderNOC), nocCert); err != nil {
		return err
	}
	if icaCert != nil {
		if err = e.Put(tlv.ContextTag(TagTBEDataSenderICAC), icaCert); err != nil {
			return err
		}
	}
	if err = e.Put(tlv.ContextTag(TagTBEDataSignature), tbsData2Signature); err != nil {
		return err
	}

	s.mNewResumptionId = make([]byte, resumptionIdSize)
	if _, err = rand.Read(s.mNewResumptionId); err != nil {
		return err
	}
	if err = e.Put(tlv.ContextTag(TagTBEDataResumptionID), s.mNewResumptionId); err != nil {
		return err
	}

	if err = e.EndContainer(outType); err != nil {
		return err
	}

	// 使用对称密钥sr2k 对 Sigma2数据进行加密
	msgR2Encrypted, _ := crypto.AES128CCMEncrypt(buf.Bytes(), sr2k, kTBEData2Nonce, nil, crypto.AEADMicLengthBytes)

	buf = new(bytes.Buffer)
	e2 := tlv.NewEncoder(buf)
	outType = tlv.TypeNotSpecified
	if outType, err = e2.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure); err != nil {
		return err
	}
	if err = e2.Put(tlv.ContextTag(1), msgRand); err != nil {
		return err
	}
	if err = e2.Put(tlv.ContextTag(2), s.LocalSessionId().Unwrap()); err != nil {
		return err
	}
	if err = e2.Put(tlv.ContextTag(3), s.mEphemeralKey.PubBytes()); err != nil {
		return err
	}
	if err = e2.Put(tlv.ContextTag(4), msgR2Encrypted); err != nil {
		return err
	}

	if s.LocalMRPConfig != nil {
		if err = s.encodeMRPParameters(e2, tlv.ContextTag(5), s.LocalMRPConfig); err != nil {
			return err
		}
	}

	if err = e2.EndContainer(outType); err != nil {
		return err
	}

	//记录下Hash值
	s.mCommissioningHash.Write(buf.Bytes())

	if err = s.exchangeContext.SendMessage(CASESigma2, buf.Bytes(), bitflags.Some(messageing.ExpectResponse)); err != nil {
		return err
	}

	s.mState = sentSigma2

	log.Info("CASESession", "SendSigma2", "Sent sigma2 message")

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
	return s, nil
}

func (s *CASESession) OnExchangeCreationFailed(delegate messageing.ExchangeDelegate) {
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
	s.mCommissioningHash = sha256.New()
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
	config *messageing.ReliableMessageProtocolConfig,
) error {
	err := s.Init(sessionManger, policy, delegate, previouslyEstablishedPeer)
	if err != nil {
		return err
	}
	s.mFabricsTable = fabrics
	s.role = transport.KSessionRoleResponder
	s.mSessionResumptionStorage = storage
	s.LocalMRPConfig = config

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

func (s *CASESession) ValidateReceivedMessage(ec *messageing.ExchangeContext, header *raw.PayloadHeader, msg *system.PacketBufferHandle) error {
	if s.exchangeContext != nil {
		if s.exchangeContext != ec {
			return chip.MATTER_ERROR_INVALID_ARGUMENT
		}
	} else {
		s.exchangeContext = ec
	}
	s.exchangeContext.UseSuggestedResponseTimeout(kExpectedHighProcessingTime)
	if msg.IsNull() {
		return chip.MATTER_ERROR_INVALID_ARGUMENT
	}
	return nil
}

func (s *CASESession) OnSuccessStatusReport() {
	s.finish()
}

func (s *CASESession) ConstructSaltSigma2(rand []byte, publicKey []byte, ipk []byte) (saltSpan []byte, err error) {
	//md := make([]byte, crypto.Sha256HashLength)
	size := ipkSize + sigmaParamRandomNumberSize + crypto.P256PublicKeyLength + crypto.Sha256HashLength
	buf := bytes.NewBuffer(saltSpan)
	_, err = buf.Write(ipk)
	_, err = buf.Write(rand)
	_, err = buf.Write(publicKey)
	_, err = buf.Write(s.mCommissioningHash.Sum(nil))
	if buf.Len() != size {
		return nil, chip.ErrorInvalidArgument
	}
	return saltSpan, nil
}

func (s *CASESession) ConstructTBSData(nocCert []byte, ICACert []byte, pubKey []byte, receiverPubKey []byte) (enc []byte, err error) {
	buf := new(bytes.Buffer)
	e := tlv.NewEncoder(buf)
	outContainerType := tlv.TypeNotSpecified
	if outContainerType, err = e.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure); err != nil {
		return
	}
	if err = e.Put(tlv.ContextTag(1), nocCert); err != nil {
		return
	}
	if err = e.Put(tlv.ContextTag(2), ICACert); err != nil {
		return
	}
	if err = e.Put(tlv.ContextTag(3), pubKey); err != nil {
		return
	}
	if err = e.Put(tlv.ContextTag(4), receiverPubKey); err != nil {
		return
	}

	if err = e.EndContainer(outContainerType); err != nil {
		return
	}
	return buf.Bytes(), nil
}
