package secure_channel

import (
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/platform/system"
	"github.com/galenliu/chip/protocols"
	log "golang.org/x/exp/slog"
)

type CASEServerBase interface {
	messageing.UnsolicitedMessageHandler
	SessionEstablishmentDelegate
	messageing.ExchangeDelegate
}

// CASEServer IMPLEMENT(UnsolicitedMessageHandler SessionEstablishmentDelegate ExchangeDelegate)
type CASEServer struct {
	mExchangeManager           messageing.ExchangeManagerBase
	mSessionResumptionStorage  SessionResumptionStorage
	mCertificateValidityPolicy credentials.CertificateValidityPolicy

	mPinnedSecureSession *transport.SessionHandle

	mPairingSession *CASESession

	mSessionManager *transport.SessionManager

	mFabrics           *credentials.FabricTable
	mGroupDataProvider *credentials.GroupDataProvider
}

func NewCASEServer() *CASEServer {
	return &CASEServer{
		mPairingSession:    NewCASESession(),
		mSessionManager:    transport.NewSessionManager(),
		mFabrics:           credentials.NewFabricTable(),
		mGroupDataProvider: &credentials.GroupDataProvider{},
	}
}

func (s *CASEServer) GetMessageDispatch() messageing.ExchangeMessageDispatchBase {
	//TODO implement me
	panic("implement me")
}

func (s *CASEServer) ListenForSessionEstablishment(
	mgr *messageing.ExchangeManager,
	sessionManager *transport.SessionManager,
	fabrics *credentials.FabricTable,
	storage lib.SessionResumptionStorage,
	policy credentials.CertificateValidityPolicy,
	responderGroupDataProvider *credentials.GroupDataProvider,
) error {
	s.mSessionManager = sessionManager
	s.mSessionResumptionStorage = storage
	s.mCertificateValidityPolicy = policy
	s.mFabrics = fabrics
	s.mExchangeManager = mgr
	s.mGroupDataProvider = responderGroupDataProvider
	s.Session().SetGroupDataProvider(s.mGroupDataProvider)
	s.PrepareForSessionEstablishment(lib.UndefinedScopedNodeId())
	return nil
}

func (s *CASEServer) PrepareForSessionEstablishment(previouslyEstablishedPeer *lib.ScopedNodeId) {

	log.Info("CASE Server enabling CASE session setups")
	err := s.mExchangeManager.RegisterUnsolicitedMessageHandlerForType(protocols.New(protocolId, nil), CASE_Sigma1, s)
	if err != nil {
		log.Info(err.Error())
	}
	s.Session().Clear()
	if s.mPinnedSecureSession != nil {
		s.mPinnedSecureSession.ClearValue()
	}
	err = s.Session().PrepareForSessionEstablishment(
		s.mSessionManager,
		s.mFabrics,
		s.mSessionResumptionStorage,
		s.mCertificateValidityPolicy,
		s,
		previouslyEstablishedPeer,
		session.GetLocalMRPConfig(),
	)
	if err != nil {
		log.Info(err.Error())
	}
	s.mPinnedSecureSession = s.Session().CopySecureSession()
}

func (s *CASEServer) InitCASEHandshake(ec *messageing.ExchangeContext) error {
	if ec == nil {
		return lib.MATTER_ERROR_INVALID_ARGUMENT
	}
	ec.SetDelegate(s.Session())
	return nil
}

func (s *CASEServer) OnUnsolicitedMessageReceived(header *raw.PayloadHeader) (messageing.ExchangeDelegate, error) {
	return s, nil
}

func (s *CASEServer) OnMessageReceived(ec *messageing.ExchangeContext, header *raw.PayloadHeader, buf *system.PacketBufferHandle) error {
	log.Info("CASE Server received Sigma1 message. Starting handshake.", "EC", ec)
	err := s.InitCASEHandshake(ec)

	err = s.mExchangeManager.UnregisterUnsolicitedMessageHandlerForType(protocols.New(protocolId, nil), CASE_Sigma1)
	if err != nil {
		return err
	}
	return s.Session().OnMessageReceived(ec, header, buf)
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

func (s *CASEServer) Session() *CASESession {
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
