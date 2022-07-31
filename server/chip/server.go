package chip

import (
	"github.com/galenliu/chip/access"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials"
	storage2 "github.com/galenliu/chip/crypto/persistent_storage"
	"github.com/galenliu/chip/device"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/secure_channel"
	"github.com/galenliu/chip/server"
	"github.com/galenliu/chip/server/dnssd"
	"github.com/galenliu/chip/storage"
	"github.com/galenliu/chip/transport"
	log "github.com/sirupsen/logrus"
	"net"
	"net/netip"
	"sync"
)

var sDeviceTypeResolver = access.DeviceTypeResolver{}

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
	mFabricTable                   *credentials.FabricTable
	mCommissioningWindowManager    dnssd.CommissioningWindowManager
	mDeviceStorage                 storage.StorageDelegate //unknown
	mAccessControl                 access.AccessControler
	mOpCerStore                    credentials.PersistentStorageOpCertStore
	mOperationalKeystore           storage2.PersistentStorageOperationalKeystore
	mCertificateValidityPolicy     credentials.CertificateValidityPolicy

	mGroupsProvider           credentials.GroupDataProvider
	mTestEventTriggerDelegate server.TestEventTriggerDelegate
	mFabricDelegate           credentials.ServerFabricDelegate
	mSessionResumptionStorage any
	mMessageCounterManager    *secure_channel.MessageCounterManager
	mUnsolicitedStatusHandler *secure_channel.UnsolicitedStatusHandler

	mAttributePersister lib.AttributePersistenceProvider //unknown
	mAclStorage         server.AclStorage
	mTransports         transport.Transport
	mExchangeMgr        messageing.ExchangeManager
	mSessions           transport.SessionManager
	mListener           credentials.GroupDataProviderListener
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

	//if initParams.StorageDelegate == nil ||
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

		s.mFabricTable = credentials.NewFabricTable()
		err := s.mFabricTable.Init(fabricTableInitParams)
		if err != nil {
			return nil, err
		}
	}

	s.mAccessControl = access.NewAccessControl()
	err = s.mAccessControl.Init(initParams.AccessDelegate, sDeviceTypeResolver)
	if err != nil {
		return nil, err
	}
	access.SetAccessControl(s.mAccessControl)

	s.mAclStorage = initParams.AclStorage
	err = s.mAclStorage.Init(s.mDeviceStorage, s.mFabricTable)
	if err != nil {
		return nil, err
	}

	s.mGroupsProvider = initParams.GroupDataProvider
	SetGroupDataProvider(s.mGroupsProvider)

	s.mTestEventTriggerDelegate = initParams.TestEventTriggerDelegate

	deviceInfoProvider := device.GetDeviceInfoProvider()
	if deviceInfoProvider != nil {
		deviceInfoProvider.SetStorageDelegate(s.mDeviceStorage)
	}

	var udpParams transport.UdpListenParameters
	s.mTransports = transport.NewUdbTransportImpl()
	err = s.mTransports.Init(udpParams.SetAddress(netip.AddrPortFrom(netip.IPv6Unspecified(), initParams.UserDirectedCommissioningPort)))

	s.mListener = credentials.NewGroupDataProviderListenerImpl()
	err = s.mListener.Init(s) // TODO
	if err != nil {
		return nil, err
	}
	s.mGroupsProvider.SetListener(s.mListener)

	s.mSessions = transport.NewSessionManagerImpl()
	err = s.mSessions.Init(s.mTransports, s.mDeviceStorage, s.GetFabricTable())
	if err != nil {
		return nil, err
	}

	s.mFabricDelegate = credentials.NewServerFabricDelegateImpl()
	err = s.mFabricDelegate.Init(s)
	if err != nil {
		return nil, err
	}

	s.mFabricTable.AddFabricDelegate(s.mFabricDelegate)

	s.mExchangeMgr = messageing.NewExchangeManagerImpl()
	err = s.mExchangeMgr.Init(s.mSessions)
	if err != nil {
		return nil, err
	}

	s.mMessageCounterManager = secure_channel.NewMessageCounterManager()
	err = s.mMessageCounterManager.Init(s.mExchangeMgr)
	if err != nil {
		return s, err
	}

	s.mUnsolicitedStatusHandler = secure_channel.NewUnsolicitedStatusHandler()
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

	discoveryService := dnssd.NewDnssdServer()
	discoveryService.SetFabricTable(s.mFabricTable)
	discoveryService.SetCommissioningModeProvider(s.mCommissioningWindowManager)

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

	//discoveryService.SetSecuredPort(s.mTransports.GetTransport().GetImplAtIndex().GetBoundPort())
	discoveryService.SetUnsecuredPort(s.mUserDirectedCommissioningPort)
	discoveryService.SetInterfaceId(s.mInterfaceId)

	if s.GetFabricTable() != nil {
		if s.GetFabricTable().FabricCount() != 0 {
			if config.NetworkLayerBle {
				//TODO
				//如果Fabric不为零，那么设备已经被添加
				//可以在这里关闭蓝牙
			}
		}
	}

	//如果设备开启了自动配对模式，进入模式
	if config.ChipDeviceConfigEnablePairingAutostart {
		s.GetFabricTable().DeleteAllFabrics()
		err = s.mCommissioningWindowManager.OpenBasicCommissioningWindow()
		if err != nil {
			log.Panic(err.Error())
		}
	}
	discoveryService.StartServer()
	s.mInitialized = true
	return s, nil
}

// GetFabricTable 返回CHIP服务中的Fabric
func (s Server) GetFabricTable() *credentials.FabricTable {
	return s.mFabricTable
}

func (s Server) Shutdown() {

}

func (s *Server) StartServer() error {
	return nil
}

func SetGroupDataProvider(provider credentials.GroupDataProvider) {

}
