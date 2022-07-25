package server

import (
	"github.com/galenliu/chip/access"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/platform"
	"github.com/galenliu/chip/server/dnssd"
	"github.com/galenliu/chip/storage"
	"github.com/galenliu/chip/transport"
	log "github.com/sirupsen/logrus"
	"net"
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
	mFabrics                       *credentials.FabricTable
	mCommissioningWindowManager    dnssd.CommissioningWindowManager
	mDeviceStorage                 storage.PersistentStorageDelegate //unknown
	mAccessControl                 access.AccessControl
	mOpCerStore                    credentials.PersistentStorageOpCertStore
	mOperationalKeystore           crypto.PersistentStorageOperationalKeystore
	mCertificateValidityPolicy     credentials.CertificateValidityPolicy

	mGroupsProvider           credentials.GroupDataProvider
	mTestEventTriggerDelegate TestEventTriggerDelegate
	mFabricDelegate           credentials.ServerFabricDelegate
	mSessionResumptionStorage any
	mExchangeMgr              messageing.ExchangeManager
	mAttributePersister       lib.AttributePersistenceProvider //unknown
	mAclStorage               AclStorage
	mTransports               transport.TransportManager
	mSessions                 transport.SessionManager
	mListener                 credentials.GroupDataProviderListener
}

func NewCHIPServer(initParams *ServerInitParams) (*Server, error) {
	s := &Server{}
	log.Printf("app server initializing")

	var err error

	s.mOperationalServicePort = initParams.OperationalServicePort
	s.mUserDirectedCommissioningPort = initParams.UserDirectedCommissioningPort
	s.mInterfaceId = initParams.InterfaceId

	if initParams.PersistentStorageDelegate == nil ||
		initParams.AccessDelegate == nil ||
		initParams.GroupDataProvider == nil ||
		initParams.OperationalKeystore == nil ||
		initParams.OpCertStore == nil {
		return nil, lib.CHIP_ERROR_INVALID_ARGUMENT
	}

	s.mDeviceStorage = initParams.PersistentStorageDelegate
	s.mSessionResumptionStorage = initParams.SessionResumptionStorage
	s.mOperationalKeystore = initParams.OperationalKeystore
	s.mOpCerStore = initParams.OpCertStore

	s.mCertificateValidityPolicy = initParams.CertificateValidityPolicy

	err = s.mAttributePersister.Init(s.mDeviceStorage)
	if err != nil {
		return nil, err
	}

	{
		fabricTableInitParams := credentials.NewFabricTableInitParams()
		fabricTableInitParams.Storage = s.mDeviceStorage
		fabricTableInitParams.OperationalKeystore = s.mOperationalKeystore
		fabricTableInitParams.OpCertStore = s.mOpCerStore
		err := s.mFabrics.Init(fabricTableInitParams)
		if err != nil {
			return nil, err
		}
	}

	err = s.mAccessControl.Init(initParams.AccessDelegate, sDeviceTypeResolver)
	if err != nil {
		return nil, err
	}
	access.SetAccessControl(s.mAccessControl)

	s.mAclStorage = initParams.AclStorage
	err = s.mAclStorage.Init(s.mDeviceStorage, s.mFabrics)
	if err != nil {
		return nil, err
	}

	s.mGroupsProvider = initParams.GroupDataProvider
	SetGroupDataProvider(s.mGroupsProvider)

	s.mTestEventTriggerDelegate = initParams.TestEventTriggerDelegate

	deviceInfoProvider := platform.GetDeviceInfoProvider()
	if deviceInfoProvider != nil {
		deviceInfoProvider.SetStorageDelegate(s.mDeviceStorage)
	}

	err = s.mTransports.Init()

	err = s.mListener.Init(s)
	if err != nil {
		return nil, err
	}
	s.mGroupsProvider.SetListener(s.mListener)

	err = s.mSessions.Init(s.mTransports, s.mDeviceStorage, s.GetFabricTable())
	if err != nil {
		return nil, err
	}

	err = s.mFabricDelegate.Init(s)
	if err != nil {
		return nil, err
	}

	s.mFabrics.AddFabricDelegate(s.mFabricDelegate)

	err = s.mExchangeMgr.Init(s.mSessions)
	if err != nil {
		return nil, err
	}

	//err = mMessageCounterManager.Init(&mExchangeMgr);
	//SuccessOrExit(err);

	//err = mUnsolicitedStatusHandler.Init(&mExchangeMgr);
	//SuccessOrExit(err);

	err = s.mCommissioningWindowManager.Init(s)
	if err != nil {
		return nil, err
	}
	s.mCommissioningWindowManager.SetAppDelegate(initParams.AppDelegate)

	discoveryService := dnssd.NewDnssdServer()
	discoveryService.SetFabricTable(s.mFabrics)
	discoveryService.SetCommissioningModeProvider(s.mCommissioningWindowManager)

	//err = chip::app::InteractionModelEngine::GetInstance()->Init(&mExchangeMgr, &GetFabricTable());
	//SuccessOrExit(err);

	//chip::Dnssd::Resolver::Instance().Init(DeviceLayer::UDPEndPointManager());

	//err = sGlobalEventIdCounter.Init(mDeviceStorage, &DefaultStorageKeyAllocator::IMEventNumber,
	//	CHIP_DEVICE_CONFIG_EVENT_ID_COUNTER_EPOCH);
	//SuccessOrExit(err);

	//{
	//	::chip::app::LogStorageResources logStorageResources[] = {
	//	{ &sDebugEventBuffer[0], sizeof(sDebugEventBuffer), ::chip::app::PriorityLevel::Debug },
	//	{ &sInfoEventBuffer[0], sizeof(sInfoEventBuffer), ::chip::app::PriorityLevel::Info },
	//	{ &sCritEventBuffer[0], sizeof(sCritEventBuffer), ::chip::app::PriorityLevel::Critical }
	//};
	//
	//chip::app::EventManagement::GetInstance().Init(&mExchangeMgr, CHIP_NUM_EVENT_LOGGING_BUFFERS, &sLoggingBuffer[0],
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
			if config.ChipConfigNetworkLayerBle {
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
	return s, nil
}

// GetFabricTable 返回CHIP服务中的Fabric
func (s Server) GetFabricTable() *credentials.FabricTable {
	return s.mFabrics
}

func (s Server) Shutdown() {

}

func (s *Server) StartServer() error {
	return nil
}

func SetGroupDataProvider(provider credentials.GroupDataProvider) {

}
