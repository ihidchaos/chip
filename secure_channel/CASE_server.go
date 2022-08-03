package secure_channel

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/transport"
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

type CASEServer struct {
	mExchangeManager           messageing.ExchangeManager
	mSessionResumptionStorage  SessionResumptionStorage
	mCertificateValidityPolicy credentials.CertificateValidityPolicy

	mPinnedSecureSession transport.SessionHandle

	mPairingSession *CASESession

	mSessionManger transport.SessionManager

	mFabrics           *credentials.FabricTable
	mGroupDataProvider credentials.GroupDataProvider
}

func NewCASEServer() *CASEServer {
	return &CASEServer{}
}

func (s *CASEServer) ListenForSessionEstablishment(mgr messageing.ExchangeManager, sessionManager transport.SessionManager, fabrics *credentials.FabricTable, storage lib.SessionResumptionStorage, policy credentials.CertificateValidityPolicy, provider credentials.GroupDataProvider) error {
	s.mSessionManger = sessionManager
	s.mSessionResumptionStorage = storage
	s.mCertificateValidityPolicy = policy
	s.mFabrics = fabrics
	s.mExchangeManager = mgr
	s.mGroupDataProvider = provider
	s.GetSession().SetGroupDataProvider(s.mGroupDataProvider)
	s.PrepareForSessionEstablishment(lib.UndefinedScopedNodeId)
	return nil
}

func (s *CASEServer) PrepareForSessionEstablishment(nodeId lib.ScopedNodeId) {
	log.Printf("CASE Server enabling CASE session setups")
	err := s.mExchangeManager.RegisterUnsolicitedMessageHandlerForType(CaseSigma1, s)
	if err != nil {
		log.Printf(err.Error())
	}
	s.GetSession().Clear()
	s.mPinnedSecureSession.ClearValue()

	err = s.GetSession().PrepareForSessionEstablishment(s.mSessionManger, s.mFabrics, s.mSessionResumptionStorage, s.mCertificateValidityPolicy, s, nodeId, transport.GetLocalMRPConfig())
	if err != nil {
		log.Panic(err.Error())
	}
	//s.mPinnedSecureSession = s.GetSession().CopySecureSession()
}

func (s *CASEServer) GetSession() *CASESession {
	return s.mPairingSession
}
