package core

import (
	"github.com/galenliu/chip/access"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/device"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/pkg/dnssd"
	"github.com/galenliu/chip/pkg/store"
	DeviceLayer "github.com/galenliu/chip/platform/device_layer"
	sc "github.com/galenliu/chip/protocols/secure_channel"
	log "golang.org/x/exp/slog"
	"net"
	"net/netip"
	"sync/atomic"
)

var sDeviceTypeResolver = access.DeviceTypeResolver{}

type ServerTransportMgr struct {
	*transport.Manager
}

func NewServerTransportMgr(transports ...raw.TransportBase) *ServerTransportMgr {
	return &ServerTransportMgr{
		Manager: transport.NewManager(transports...),
	}
}

type AppDelegate interface {
	OnCommissioningSessionStarted()
	OnCommissioningSessionStopped()
	OnCommissioningWindowOpened()
	OnCommissioningWindowClosed()
}

type Server struct {
	mTransports *ServerTransportMgr
	mSessions   *transport.SessionManager
	mCASEServer *sc.CASEServer

	mCASESessionManager *CASESessionManager
	mCASEClientPool     *CASEClientPool

	mOperationalServicePort        uint16
	mUserDirectedCommissioningPort uint16
	mInterfaceId                   net.Interface
	mDnssd                         dnssd.Base
	mFabrics                       *credentials.FabricTable
	mCommissioningWindowManager    dnssd.CommissioningWindowManager
	mDeviceStorage                 store.KvsPersistentStorageBase
	mAccessControl                 access.Controller
	mOpCerStore                    credentials.PersistentStorageOpCertStore
	mOperationalKeystore           crypto.OperationalKeystore
	mCertificateValidityPolicy     credentials.CertificateValidityPolicy

	mGroupsProvider           *credentials.GroupDataProvider
	mTestEventTriggerDelegate TestEventTriggerDelegate
	mFabricDelegate           credentials.FabricTableDelegate

	mSessionResumptionStorage lib.SessionResumptionStorage
	mMessageCounterManager    *sc.MessageCounterManager
	mUnsolicitedStatusHandler sc.UnsolicitedStatusHandler

	mAttributePersister lib.AttributePersistenceProvider //unknown

	mDevicePool  OperationalDeviceProxyPool
	mAclStorage  AclStorage
	mExchangeMgr *messageing.ExchangeManager

	mListener    GroupDataProviderListener
	mInitialized bool
}

var defaultServer atomic.Value

func DefaultServer() *Server {
	server := defaultServer.Load().(*Server)
	return server
}

func init() {
	sev := NewServer()
	defaultServer.Store(sev)
}

func NewServer() *Server {
	return &Server{
		mInterfaceId:              net.Interface{},
		mTestEventTriggerDelegate: TestEventTriggerDelegate{},
		mCASEClientPool:           NewCASEClientPool(config.DeviceMaxActiveCASEClients),
		mCASEServer:               sc.NewCASEServer(),
		mMessageCounterManager:    sc.NewMessageCounterManager(),
		mCASESessionManager:       NewCASESessionManager(),
		mDevicePool:               OperationalDeviceProxyPool{},
	}
}

func (s *Server) Init(initParams *InitParams) error {
	log.Info("app server initializing")
	var err error
	s.mOperationalServicePort = initParams.OperationalServicePort
	s.mUserDirectedCommissioningPort = initParams.UserDirectedCommissioningPort
	s.mInterfaceId = initParams.InterfaceId

	//if initParams.PersistentStorage == nil ||
	//	initParams.AccessDelegate == nil ||
	//	initParams.GroupDataProvider == nil ||
	//	initParams.OperationalKeystore == nil ||
	//	initParams.OpCertStore == nil {
	//	return nil, lib.CHIP_ERROR_INVALID_ARGUMENT
	//}

	s.mDeviceStorage = initParams.PersistentStorageDelegate
	s.mSessionResumptionStorage = initParams.SessionResumptionStorage
	s.mOperationalKeystore = initParams.OperationalKeystore
	s.mOpCerStore = initParams.OpCertStore

	s.mCertificateValidityPolicy = initParams.CertificateValidityPolicy

	s.mAttributePersister = lib.NewAttributePersistence()
	err = s.mAttributePersister.Init(s.mDeviceStorage)
	if err != nil {
		return err
	}

	{
		fabricTableInitParams := credentials.NewFabricTableInitParams()
		fabricTableInitParams.Storage = s.mDeviceStorage
		fabricTableInitParams.OperationalKeystore = s.mOperationalKeystore
		fabricTableInitParams.OpCertStore = s.mOpCerStore

		s.mFabrics = credentials.NewFabricTable()
		err := s.mFabrics.Init(fabricTableInitParams)
		if err != nil {
			return err
		}
	}

	{
		accessControl := access.NewAccessControl()
		err = accessControl.Init(initParams.AccessDelegate, sDeviceTypeResolver)
		if err != nil {
			return err
		}
		access.SetAccessControl(s.mAccessControl)
		s.mAccessControl = accessControl
	}

	{
		aclStorage := initParams.AclStorage
		err = aclStorage.Init(s.mDeviceStorage, s.mFabrics)
		if err != nil {
			log.Error("AclStorage init", err)
		}
		s.mAclStorage = aclStorage
	}

	s.mGroupsProvider = initParams.GroupDataProvider
	credentials.SetGroupDataProvider(s.mGroupsProvider)

	s.mTestEventTriggerDelegate = initParams.TestEventTriggerDelegate

	deviceInfoProvider := device.DefaultInfoProvider()
	if deviceInfoProvider != nil {
		deviceInfoProvider.SetStorageDelegate(s.mDeviceStorage)
	}

	{
		udpParams := raw.UdpListenParameters{}
		udp := raw.NewUdpTransport()
		err = udp.Init(udpParams.SetAddress(netip.AddrPortFrom(netip.IPv6Unspecified(), s.mOperationalServicePort)))
		if err != nil {
			log.Error("UdpTransport init", err)
		}
		s.mTransports = NewServerTransportMgr(udp)
	}

	{
		s.mListener = &GroupDataProviderListenerImpl{s}
		err = s.mListener.Init(s) // TODO
		if err != nil {
			return err
		}
		s.mGroupsProvider.SetListener(s.mListener)
	}

	{
		sessionManager := transport.NewSessionManager()
		err = sessionManager.Init(DeviceLayer.SystemLayer(), s.mTransports, s.mMessageCounterManager, s.mDeviceStorage, s.GetFabricTable())
		if err != nil {
			return err
		}
		s.mSessions = sessionManager
	}

	{
		fabricDelegate := NewServerFabricDelegateImpl()
		_ = fabricDelegate.Init(s)
		s.mFabricDelegate = fabricDelegate
		err := s.mFabrics.AddFabricDelegate(s.mFabricDelegate)
		if err != nil {
			return err
		}
	}

	{
		exchangeMgr := messageing.NewExchangeManager()
		err = exchangeMgr.Init(s.mSessions)
		if err != nil {
			return err
		}
		s.mExchangeMgr = exchangeMgr
	}

	{
		err = s.mMessageCounterManager.Init(s.mExchangeMgr)
		if err != nil {
			return err
		}

	}

	s.mUnsolicitedStatusHandler = sc.NewUnsolicitedStatusHandler()
	err = s.mUnsolicitedStatusHandler.Init(s.mExchangeMgr)
	if err != nil {
		return err
	}

	s.mCommissioningWindowManager = dnssd.NewCommissioningWindowManagerImpl()
	err = s.mCommissioningWindowManager.Init(s)
	if err != nil {
		return err
	}
	s.mCommissioningWindowManager.SetAppDelegate(initParams.AppDelegate)

	dnssdServer := dnssd.New()
	dnssdServer.SetFabricTable(s.mFabrics)
	dnssdServer.SetCommissioningModeProvider(s.mCommissioningWindowManager)

	if config.PairingAutostart {
		err := s.mCommissioningWindowManager.OpenBasicCommissioningWindow()
		if err != nil {
			return err
		}
	}

	//err = chip::app::InteractionModelEngine::GetInstance()->initCommissionableData(&mExchangeMgr, &GetFabricTable());
	//SuccessOrExit(err);

	//chip::Dnssd::Resolver::Instance().initCommissionableData(DeviceLayer::UDPEndPointManager());

	//err = sGlobalEventIdCounter.initCommissionableData(mDeviceStorage, &DefaultStorageKeyAllocator::IMEventNumber,
	//	CHIP_DEVICE_CONFIG_EVENT_ID_COUNTER_EPOCH);
	//SuccessOrExit(err);

	//{
	//	::chip::app::LogStorageResources logStorageResources[] = {
	//	{ &sDebugEventBuffer[0], sizeof(sDebugEventBuffer), ::chip::app::PriorityLevel::Debug },
	//	{ &sInfoEventBuffer[0], sizeof(sInfoEventBuffer), ::chip::app::PriorityLevel::Info },
	//	{ &sCritEventBuffer[0], sizeof(sCritEventBuffer), ::chip::app::PriorityLevel::Critical }
	//};
	//
	//chip::app::EventManagement::GetInstance().initCommissionableData(&mExchangeMgr, CHIP_NUM_EVENT_LOGGING_BUFFERS, &sLoggingBuffer[0],
	//	&logStorageResources[0], &sGlobalEventIdCounter);
	//}
	//#endif // CHIP_CONFIG_ENABLE_SERVER_IM_EVENT

	//// This initializes clusters, so should come after lower level initialization.
	messageing.InitDataModelHandler(s.mExchangeMgr)

	//#if defined(CHIP_APP_USE_ECHO)
	//	err = InitEchoHandler(&mExchangeMgr);
	//SuccessOrExit(err);
	//#endif

	//
	// We need to advertise the port that we're listening to for unsolicited messages over UDP. However, we have both a IPv4
	// and IPv6 endpoint to pick from. Given that the listen port passed in may be set to 0 (which then has the kernel select
	// a valid port at bind time), that will result in two possible ports being provided back from the resultant endpoint
	// initializations. Since IPv6 is POR for Matter, let's go ahead and pick that port.
	if udpTransport := s.mTransports.GetUpdImpl(); udpTransport != nil {
		dnssdServer.SetSecuredPort(udpTransport.BoundPort())
	}
	dnssdServer.SetUnsecuredPort(s.mUserDirectedCommissioningPort)
	dnssdServer.SetInterfaceId(s.mInterfaceId)

	if s.GetFabricTable() != nil {
		if s.GetFabricTable().FabricCount() != 0 {
			if config.NetworkLayerBle {
				//TODO
				//如果Fabric不为零，那么设备已经被添加
				//可以在这里关闭蓝牙
			}
		}
	}
	caseSessionManagerConfig := &CASESessionManagerConfig{
		SessionInitParams: DeviceProxyInitParams{
			SessionManager:            s.mSessions,
			SessionResumptionStorage:  s.mSessionResumptionStorage,
			CertificateValidityPolicy: s.mCertificateValidityPolicy,
			ExchangeMgr:               s.mExchangeMgr,
			FabricTable:               s.mFabrics,
			ClientPool:                s.mCASEClientPool,
			GroupDataProvider:         s.mGroupsProvider,
			MrpLocalConfig:            session.GetLocalMRPConfig(),
		},
		DevicePool: s.mDevicePool,
	}

	err = s.mCASESessionManager.Init(SystemLayer(), caseSessionManagerConfig)
	if err != nil {
		return err
	}

	err = s.mCASEServer.ListenForSessionEstablishment(s.mExchangeMgr, s.mSessions, s.mFabrics, s.mSessionResumptionStorage, s.mCertificateValidityPolicy, s.mGroupsProvider)
	if err != nil {
		return err
	}

	//如果设备开启了自动配对模式，进入模式
	if config.PairingAutostart {
		s.GetFabricTable().DeleteAllFabrics()
		err = s.mCommissioningWindowManager.OpenBasicCommissioningWindow()
		if err != nil {
			return err
		}
	}
	dnssdServer.StartServer()
	s.mInitialized = true
	return nil
}

// GetFabricTable 返回CHIP服务中的Fabric
func (s *Server) GetFabricTable() *credentials.FabricTable {
	return s.mFabrics
}

func (s *Server) Shutdown() {

}

func (s *Server) StartServer() error {
	return nil
}
