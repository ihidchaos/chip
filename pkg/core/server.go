package core

import (
	"github.com/galenliu/chip/access"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials"
	storage2 "github.com/galenliu/chip/crypto/operational_storage"
	"github.com/galenliu/chip/device"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/messageing/transport"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/pkg/dnssd"
	"github.com/galenliu/chip/pkg/storage"
	secure_channel2 "github.com/galenliu/chip/protocols/secure_channel"
	log "github.com/sirupsen/logrus"
	"net"
	"net/netip"
	"sync"
)

var sDeviceTypeResolver = access.DeviceTypeResolver{}

type ServerTransportManager interface {
	transport.ManagerBase
	GetImplAtIndex(index int) raw.TransportBase
}

type AppDelegate interface {
	OnCommissioningSessionStarted()
	OnCommissioningSessionStopped()
	OnCommissioningWindowOpened()
	OnCommissioningWindowClosed()
}

type Server struct {
	mOperationalServicePort        uint16
	mUserDirectedCommissioningPort uint16
	mInterfaceId                   net.Interface
	mDnssd                         dnssd.DnssdServer
	mFabrics                       *credentials.FabricTable
	mCommissioningWindowManager    dnssd.CommissioningWindowManager
	mDeviceStorage                 storage.KvsPersistentStorageDelegate
	mAccessControl                 access.AccessControler
	mOpCerStore                    credentials.PersistentStorageOpCertStore
	mOperationalKeystore           storage2.OperationalKeystore
	mCertificateValidityPolicy     credentials.CertificateValidityPolicy

	mGroupsProvider           credentials.GroupDataProvider
	mTestEventTriggerDelegate TestEventTriggerDelegate
	mFabricDelegate           credentials.FabricTableDelegate
	mCASEClientPool           CASEClientPool
	mCASEServer               *secure_channel2.CASEServer
	mSessionResumptionStorage lib.SessionResumptionStorage
	mMessageCounterManager    *secure_channel2.MessageCounterManager
	mUnsolicitedStatusHandler *secure_channel2.UnsolicitedStatusHandler

	mAttributePersister lib.AttributePersistenceProvider //unknown
	mCASESessionManager *CASESessionManager
	mDevicePool         OperationalDeviceProxyPool
	mAclStorage         AclStorage
	mTransports         ServerTransportManager
	mExchangeMgr        messageing.ExchangeManager
	mSessions           transport.SessionManager
	mListener           GroupDataProviderListener
	mInitialized        bool
}

var _chipServerInstance *Server
var _chipServerOnce sync.Once

func GetChipServerInstance() *Server {
	_chipServerOnce.Do(func() {
		if _chipServerInstance == nil {
			_chipServerInstance = &Server{}
		}
	})
	return _chipServerInstance
}

func NewCHIPServer() *Server {
	return GetChipServerInstance()
}

func (s *Server) Init(initParams *InitParams) (*Server, error) {

	log.Printf("app server initializing")

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
		return nil, err
	}

	{
		fabricTableInitParams := credentials.NewFabricTableInitParams()
		fabricTableInitParams.Storage = s.mDeviceStorage
		fabricTableInitParams.OperationalKeystore = s.mOperationalKeystore
		fabricTableInitParams.OpCertStore = s.mOpCerStore

		s.mFabrics = credentials.NewFabricTable()
		err := s.mFabrics.Init(fabricTableInitParams)
		if err != nil {
			return nil, err
		}
	}

	{
		accessControl := access.NewAccessControl()
		err = accessControl.Init(initParams.AccessDelegate, sDeviceTypeResolver)
		if err != nil {
			return nil, err
		}
		access.SetAccessControl(s.mAccessControl)
		s.mAccessControl = accessControl
	}

	{
		aclStorage := initParams.AclStorage
		err = aclStorage.Init(s.mDeviceStorage, s.mFabrics)
		if err != nil {
			log.Panic(err.Error())
		}
		s.mAclStorage = aclStorage
	}

	s.mGroupsProvider = initParams.GroupDataProvider
	credentials.SetGroupDataProvider(s.mGroupsProvider)

	s.mTestEventTriggerDelegate = initParams.TestEventTriggerDelegate

	deviceInfoProvider := device.GetDeviceInfoProvider()
	if deviceInfoProvider != nil {
		deviceInfoProvider.SetStorageDelegate(s.mDeviceStorage)
	}

	{
		udpParams := raw.UdpListenParameters{}
		udp := raw.NewUdpTransportImpl()
		err = udp.Init(udpParams.SetAddress(netip.AddrPortFrom(netip.IPv6Unspecified(), s.mOperationalServicePort)))
		if err != nil {
			log.Panic(err.Error())
		}
		s.mTransports = transport.NewTransportManagerImpl(udp)
	}

	{
		s.mListener = &GroupDataProviderListenerImpl{s}
		err = s.mListener.Init(s) // TODO
		if err != nil {
			return nil, err
		}
		s.mGroupsProvider.SetListener(s.mListener)
	}

	{
		session := transport.NewSessionManagerImpl()
		err = session.Init(s.mTransports, s.mDeviceStorage, s.GetFabricTable())
		if err != nil {
			return nil, err
		}
		s.mSessions = session
	}

	{
		fabricDelegate := NewServerFabricDelegateImpl()
		_ = fabricDelegate.Init(s)
		s.mFabricDelegate = fabricDelegate
		err := s.mFabrics.AddFabricDelegate(s.mFabricDelegate)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	{
		exchangeMgr := messageing.NewExchangeManagerImpl()
		err = exchangeMgr.Init(s.mSessions)
		if err != nil {
			log.Panic(err.Error())
		}
		s.mExchangeMgr = exchangeMgr
	}

	{
		messageCounterManager := secure_channel2.NewMessageCounterManager()
		err = messageCounterManager.Init(s.mExchangeMgr)
		if err != nil {
			return s, err
		}
		s.mMessageCounterManager = messageCounterManager
	}

	s.mUnsolicitedStatusHandler = secure_channel2.NewUnsolicitedStatusHandler()
	err = s.mUnsolicitedStatusHandler.Init(s.mExchangeMgr)
	if err != nil {
		return s, err
	}

	s.mCommissioningWindowManager = dnssd.NewCommissioningWindowManagerImpl()
	err = s.mCommissioningWindowManager.Init(s)
	if err != nil {
		return nil, err
	}
	s.mCommissioningWindowManager.SetAppDelegate(initParams.AppDelegate)

	dnssdServer := dnssd.NewDnssdInstance()
	dnssdServer.SetFabricTable(s.mFabrics)
	dnssdServer.SetCommissioningModeProvider(s.mCommissioningWindowManager)

	if config.ChipDeviceConfigEnablePairingAutostart {
		err := s.mCommissioningWindowManager.OpenBasicCommissioningWindow()
		if err != nil {
			log.Panic(err.Error())
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

	dnssdServer.SetSecuredPort(s.mTransports.GetImplAtIndex(0).GetBoundPort())
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
			MrpLocalConfig:            transport.GetLocalMRPConfig(),
		},
		DevicePool: s.mDevicePool,
	}

	s.mCASESessionManager = NewCASESessionManager()
	err = s.mCASESessionManager.Init(SystemLayer(), caseSessionManagerConfig)
	if err != nil {
		log.Panic(err.Error())
	}

	s.mCASEServer = secure_channel2.NewCASEServer()
	err = s.mCASEServer.ListenForSessionEstablishment(s.mExchangeMgr, s.mSessions, s.mFabrics, s.mSessionResumptionStorage, s.mCertificateValidityPolicy, s.mGroupsProvider)
	if err != nil {
		log.Panic(err.Error())
	}

	//如果设备开启了自动配对模式，进入模式
	if config.ChipDeviceConfigEnablePairingAutostart {
		s.GetFabricTable().DeleteAllFabrics()
		err = s.mCommissioningWindowManager.OpenBasicCommissioningWindow()
		if err != nil {
			log.Panic(err.Error())
		}
	}
	dnssdServer.StartServer()
	s.mInitialized = true
	return s, nil
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

func (s *Server) GetTransportManager() transport.ManagerBase {
	return s.mTransports
}
