package secure

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/buffer"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/protocols"
	log "github.com/sirupsen/logrus"
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

	mPinnedSecureSession transport.SessionHandleBase

	mPairingSession *CASESession

	mSessionManager transport.SessionManager

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
	sessionManager transport.SessionManager,
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
	err := s.mExchangeManager.RegisterUnsolicitedMessageHandlerForType(protocols.StandardSecureChannelProtocolId, uint8(messageing.CASESigma1), s)
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
		transport.GetLocalMRPConfig(),
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
	delegate = s
	return nil
}

func (s *CASEServer) OnMessageReceived(context *messageing.ExchangeContext, header *raw.PayloadHeader, buf *buffer.PacketBuffer) error {
	err := s.mExchangeManager.UnregisterUnsolicitedMessageHandlerForType(uint8(messageing.CASESigma1))
	if err != nil {
		return err
	}
	return s.GetSession().OnMessageReceived(context, header, buf)
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
