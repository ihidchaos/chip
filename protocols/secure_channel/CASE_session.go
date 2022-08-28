package secure_channel

import (
	"bytes"
	"crypto/elliptic"
	"errors"
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
	rand "math/rand"
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
	mRemotePubKey      *crypto.P256PublicKey
	mEphemeralKey      *crypto.P256Keypair
	mSharedSecret      []byte //crypto.P256ECDHDerivedSecret
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

	mResumeResumptionId []byte //会话恢复ID 16个字节
	mNewResumptionId    []byte //会话恢复ID 16个字节

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
	case StateInitialized:
		if msgType == messageing.CASESigma1 {
			return s.HandleSigma1AndSendSigma2(buf)
		}
	case StateSentSigma1:
	case StateSentSigma1Resume:
	case StateSentSigma2:
	case StateSentSigma3:
	case StateSentSigma2Resume:
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

	if sigma1.sessionResumptionRequested && len(sigma1.resumptionId) == 16 {
	}

	err = s.FindLocalNodeFromDestinationId(sigma1.destinationId, sigma1.initiatorRandom)
	if err == nil {
		log.Infof("CASE matched destination ID: fabricIndex %x, NodeID 0x", s.mFabricIndex, s.mLocalNodeId)
	} else {
		log.Infof("CASE failed to match destination ID with local fabrics")
	}

	pubKey, err := crypto.UnmarshalP256PublicKey(sigma1.initiatorEphPubKey)
	if err != nil {
		return err
	}
	s.mRemotePubKey = &pubKey

	err = s.SendSigma2()
	if err != nil {
		//s.SendStatusReport()
		s.mState = StateInitialized
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
		return errors.New("sigma2 icaCert err")
	}
	//获取节点操作证书(NOC)
	nocCert, err := s.mFabricsTable.FetchNOCCert(s.mFabricIndex)
	if err != nil && len(nocCert) != credentials.MaxCHIPCertLength {
		return errors.New("sigma2 nocCert err")
	}

	//生成一个32字节的随机数
	msgRand := make([]byte, kSigmaParamRandomNumberSize)
	_, err = rand.Read(msgRand)

	//生成临时密钥对
	s.mEphemeralKey = s.mFabricsTable.AllocateEphemeralKeypairForCASE()
	if s.mEphemeralKey == nil {
		return errors.New("sigma2 ephemeral key error")
	}

	// 生成共享密钥
	s.mSharedSecret = s.mEphemeralKey.ECDHDeriveSecret(*s.mRemotePubKey)

	//生成一个盐值

	salt, err := s.ConstructSaltSigma2(msgRand, elliptic.Marshal(elliptic.P256(), s.mEphemeralKey.PublicKey.X, s.mEphemeralKey.Y), s.mIPK)
	if err != nil || salt == nil {
		return errors.New("sigma2 salt error")
	}

	// 使用HKDF函数派生出密钥sr2k,sr2k为16个字节
	sr2k := crypto.HKDFSha256(s.mSharedSecret, salt, KDFSR2Info)
	if sr2k == nil || len(sr2k) != crypto.SymmetricKeyLengthBytes {
		return errors.New("sigma2 hkdfSha256 error")
	}

	//msgR2SignedLen := tlv.EstimateStructOverhead(credentials.MaxCHIPCertLength, credentials.MaxCHIPCertLength, crypto.KP256PublicKeyLength, crypto.KP256PublicKeyLength)

	//msgR2Signed := s.ConstructTBSData(nocCert, icaCert, s.mEphemeralKey.MarshalPublicKey(), s.mRemotePubKey.Marshal())

	tbsData2Signature := s.mFabricsTable.SignWithOpKeypair(s.mFabricIndex).Bytes()

	tlvWriterMsg1 := tlv.NewWriter()
	err = tlvWriterMsg1.StartContainer(tlv.AnonymousTag(), tlv.Type_Structure)
	if err != nil {
		return err
	}

	err = tlvWriterMsg1.PutBytes(tlv.ContextTag(kTag_TBEData_SenderNOC), nocCert)
	if err != nil {
		return err
	}
	if len(icaCert) > 0 {
		err = tlvWriterMsg1.PutBytes(tlv.ContextTag(kTag_TBEData_SenderICAC), icaCert)
		if err != nil {
			return err
		}
	}
	err = tlvWriterMsg1.PutBytes(tlv.ContextTag(kTag_TBEData_Signature), tbsData2Signature)
	if err != nil {
		return err
	}

	s.mNewResumptionId = make([]byte, kResumptionIdSize)
	err = crypto.DRBGBytes(s.mNewResumptionId)
	if err != nil {
		return err
	}
	err = tlvWriterMsg1.PutBytes(tlv.ContextTag(kTag_TBEData_ResumptionID), s.mNewResumptionId)
	if err != nil {
		return err
	}

	err = tlvWriterMsg1.EndContainer(tlv.Type_Structure)
	if err != nil {
		return err
	}

	// 使用对称密钥sr2k 对 Sigma2数量进行加密
	msgR2Encrypted, err := crypto.AesCcmEncrypt(tlvWriterMsg1.Bytes(), sr2k, kTBEData2Nonce, crypto.AEADMicLengthBytes)

	tlvWriterMsg2 := tlv.NewWriter()
	err = tlvWriterMsg2.StartContainer(tlv.AnonymousTag(), tlv.Type_Structure)
	if err != nil {
		return err
	}
	err = tlvWriterMsg2.PutBytes(tlv.ContextTag(1), msgRand)
	if err != nil {
		return err
	}
	err = tlvWriterMsg2.Put(tlv.ContextTag(2), uint64(sessionId))
	if err != nil {
		return err
	}
	err = tlvWriterMsg2.PutBytes(tlv.ContextTag(3), s.mEphemeralKey.MarshalPublicKey())
	if err != nil {
		return err
	}
	err = tlvWriterMsg2.PutBytes(tlv.ContextTag(4), msgR2Encrypted)
	if err != nil {
		return err
	}

	if s.mLocalMRPConfig != nil {

	}
	err = tlvWriterMsg2.EndContainer(tlv.Type_Structure)
	if err != nil {
		return err
	}

	msgR2 := s.mCommissioningHash.AddData(tlvWriterMsg2.Bytes())

	err = s.mExchangeCtxt.SendMessage(messageing.CASESigma2, msgR2, messageing.ExpectResponse)
	if err != nil {
		return err
	}

	s.mState = StateSentSigma2

	log.Infof("Sent sigma2 message")

	return nil
}

func (s *CASESession) FindLocalNodeFromDestinationId(destinationId []byte, initiatorRandom []byte) error {

	for _, fabricInfo := range s.mFabricsTable.GetFabrics() {

		fabriceId := fabricInfo.GetFabricId()
		nodeId := fabricInfo.GetNodeId()
		rootPubKey, err := s.mFabricsTable.FetchRootPubkey(fabricInfo.GetFabricIndex())
		if err != nil {
			return err
		}
		ipkKeySet, err := s.mGroupDataProvider.GetIpkKeySet(fabricInfo.GetFabricIndex())
		if err != nil || (ipkKeySet.NumKeysUsed == 0 || ipkKeySet.NumKeysUsed > credentials.KEpochKeysMax) {
			continue
		}

		// Try every IPK candidate we have for a match
		for keyIdx := 0; uint8(keyIdx) < ipkKeySet.NumKeysUsed; keyIdx++ {
			candidateDestinationId, err := GenerateCaseDestinationId(ipkKeySet.EpochKeys[keyIdx].Key, initiatorRandom, rootPubKey, fabriceId, nodeId)
			if err != nil && len(candidateDestinationId) == len(destinationId) {
				s.mFabricIndex = fabricInfo.GetFabricIndex()
				s.mLocalNodeId = nodeId
			}
		}
	}
	return nil
}

func (s *CASESession) ConstructSaltSigma2(rand []byte, publicKey []byte, ipk []byte) (saltSpan []byte, err error) {
	//md := make([]byte, crypto.KSha256HashLength)

	saltSpan = make([]byte, kIpkSize+kSigmaParamRandomNumberSize+crypto.KP256PublicKeyLength+crypto.KSha256HashLength)
	buf := bytes.NewBuffer(saltSpan)
	_, err = buf.Write(ipk)
	_, err = buf.Write(rand)
	_, err = buf.Write(publicKey)
	_, err = buf.Write(s.mCommissioningHash.Bytes())
	return saltSpan, nil
}

func (s *CASESession) ConstructTBSData(senderNOC []byte, senderICAC []byte, senderPubKey []byte, receiverPubKey []byte) []byte {

	return nil
}
