package secure_channel

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	transport2 "github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/protocols"
	log "github.com/sirupsen/logrus"
)

const (
	kInitialized       = 0
	kSentSigma1        = 1
	kSentSigma2        = 2
	kSentSigma3        = 3
	kSentSigma1Resume  = 4
	kSentSigma2Resume  = 5
	kFinished          = 6
	kFinishedViaResume = 7
)

type CASEServerBase interface {
	messageing.UnsolicitedMessageHandler
	SessionEstablishmentDelegate
	messageing.ExchangeDelegate
}

// CASEServer IMPLEMENT(UnsolicitedMessageHandler SessionEstablishmentDelegate ExchangeDelegate)
type CASEServer struct {
	mExchangeManager           messageing.ExchangeManager
	mSessionResumptionStorage  SessionResumptionStorage
	mCertificateValidityPolicy credentials.CertificateValidityPolicy

	mPinnedSecureSession transport2.SessionHandle

	mPairingSession *CASESession

	mSessionManager transport2.SessionManager

	mFabrics           *credentials.FabricTable
	mGroupDataProvider credentials.GroupDataProvider
}

func NewCASEServer() *CASEServer {
	return &CASEServer{
		mPairingSession:    NewCASESession(),
		mSessionManager:    nil,
		mFabrics:           nil,
		mGroupDataProvider: nil,
	}
}

func (s *CASEServer) ListenForSessionEstablishment(
	mgr messageing.ExchangeManager,
	sessionManager transport2.SessionManager,
	fabrics *credentials.FabricTable,
	storage lib.SessionResumptionStorage,
	policy credentials.CertificateValidityPolicy,
	responderGroupDataProvider credentials.GroupDataProvider,
) error {
	s.mSessionManager = sessionManager
	s.mSessionResumptionStorage = storage
	s.mCertificateValidityPolicy = policy
	s.mFabrics = fabrics
	s.mExchangeManager = mgr
	s.mGroupDataProvider = responderGroupDataProvider
	s.GetSession().SetGroupDataProvider(s.mGroupDataProvider)
	s.PrepareForSessionEstablishment(lib.NewScopedNodeId())
	return nil
}

func (s *CASEServer) PrepareForSessionEstablishment(previouslyEstablishedPeer *lib.ScopedNodeId) {
	log.Printf("CASE Server enabling CASE session setups")
	err := s.mExchangeManager.RegisterUnsolicitedMessageHandlerForType(protocols.StandardProtocolId, CASESigma1, s)
	if err != nil {
		log.Printf(err.Error())
	}
	s.GetSession().Clear()
	if s.mPinnedSecureSession != nil {
		s.mPinnedSecureSession.ClearValue()
	}
	err = s.GetSession().PrepareForSessionEstablishment(
		s.mSessionManager,
		s.mFabrics,
		s.mSessionResumptionStorage,
		s.mCertificateValidityPolicy,
		s,
		previouslyEstablishedPeer,
		transport2.GetLocalMRPConfig(),
	)
	if err != nil {
		log.Panic(err.Error())
	}
	s.mPinnedSecureSession = s.GetSession().CopySecureSession()
}

func (s *CASEServer) InitCASEHandshake(ec *messageing.ExchangeContext) {
	ec.SetDelegate(s.GetSession())
}

func (s *CASEServer) OnUnsolicitedMessageReceived(header *raw.PayloadHeader, delegate messageing.ExchangeDelegate) error {
	//TODO implement me
	panic("implement me")
}

func (s *CASEServer) OnMessageReceived(context *messageing.ExchangeContext, header raw.PayloadHeader, data []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s *CASEServer) OnResponseTimeout(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (s *CASEServer) OnExchangeClosing(ec *messageing.ExchangeContext) {
	//TODO implement me
	panic("implement me")
}

func (s *CASEServer) OnExchangeCreationFailed(delegate messageing.ExchangeDelegate) {
	//TODO implement me
	panic("implement me")
}

func (s *CASEServer) GetSession() *CASESession {
	return s.mPairingSession
}

func (s *CASEServer) OnSessionEstablishmentError() {
	//TODO implement me
	panic("implement me")
}

func (s *CASEServer) OnSessionEstablishmentStarted() {
	//TODO implement me
	panic("implement me")
}

func (s *CASEServer) OnSessionEstablished() {
	//TODO implement me
	panic("implement me")
}
