package server

import (
	"github.com/galenliu/chip/access"
	"github.com/galenliu/chip/config/costant"
	"github.com/galenliu/chip/controller"
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/inet/udp_endpoint"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/platform/storage"
	"github.com/galenliu/chip/server/dnssd"
	"github.com/galenliu/chip/server/dnssd/manager"
	"github.com/galenliu/chip/transport"
	"github.com/galenliu/gateway/pkg/errors"
	"log"
	"net"
)

type AppDelegate interface {
	OnCommissioningSessionStarted()
	OnCommissioningSessionStopped()
	OnCommissioningWindowOpened()
	OnCommissioningWindowClosed()
}

type Config struct {
	UnsecureServicePort int
	SecureServicePort   int
}

type Server struct {
	mSecuredServicePort            uint16
	mUnsecuredServicePort          uint16
	mOperationalServicePort        uint16
	mUserDirectedCommissioningPort uint16
	mInterfaceId                   net.Interface
	config                         Config
	mDnssd                         *dnssd.Server
	mFabrics                       *credentials.FabricTable
	mCommissioningWindowManager    *manager.CommissioningWindowManager
	mDeviceStorage                 storage.PersistentStorageDelegate //unknown
	mAccessControl                 access.AccessControler
	mSessionResumptionStorage      any
	mExchangeMgr                   messageing.ExchangeManager
	mAttributePersister            lib.AttributePersistenceProvider //unknown
	mAclStorage                    *AclStorage
	mTransports                    transport.TransportManager
	mListener                      any
}

func NewServer(initParams *CommonCaseDeviceServerInitParams) *Server {
	s := &Server{}
	log.Printf("app server initializing")

	var err error
	s.mUnsecuredServicePort = initParams.OperationalServicePort
	s.mSecuredServicePort = initParams.UserDirectedCommissioningPort
	s.mInterfaceId = initParams.InterfaceId

	s.mCommissioningWindowManager.SetAppDelegate(initParams.AppDelegate)

	s.mDnssd = dnssd.Server{}.Init()
	s.mDnssd.SetFabricTable(s.mFabrics)
	s.mCommissioningWindowManager = manager.CommissioningWindowManager{}.Init(&s)
	//s.mCommissioningWindowManager.SetAppDelegate(initParams.AppDelegate)

	// Initialize PersistentStorageDelegate-based storage
	s.mDeviceStorage = initParams.PersistentStorageDelegate
	s.mSessionResumptionStorage = initParams.SessionResumptionStorage

	// Set up attribute persistence before we try to bring up the data model
	// handler.
	if s.mAttributePersister != nil {
		err = s.mAttributePersister.Init(s.mDeviceStorage)
		errors.SuccessOrExit(err)
	}

	if s.mFabrics != nil {
		err = s.mFabrics.Init(s.mDeviceStorage)
		errors.SuccessOrExit(err)
	}

	//少sDeviceTypeResolver参数
	if s.mAccessControl != nil {
		err = s.mAccessControl.Init(initParams.AccessDelegate)
		errors.SuccessOrExit(err)
	}

	errors.SuccessOrExit(err)

	s.mDnssd.SetFabricTable(s.mFabrics)
	s.mDnssd.SetCommissioningModeProvider(s.mCommissioningWindowManager)

	//mGroupsProvider = initParams.groupDataProvider;
	//SetGroupDataProvider(mGroupsProvider);
	//
	//deviceInfoprovider = DeviceLayer::GetDeviceInfoProvider();
	//if (deviceInfoprovider)
	//{
	//	deviceInfoprovider->SetStorageDelegate(mDeviceStorage);
	//}

	// This initializes clusters, so should come after lower level initialization.
	//不知道干什么的
	controller.InitDataModelHandler(s.mExchangeMgr)

	params := transport.UdpListenParameters{}
	params.SetListenPort(s.mOperationalServicePort)
	params.SetNativeParams(initParams.EndpointNativeParams)
	s.mTransports, err = transport.NewUdpTransport(udp_endpoint.UDPEndpoint{}, params)

	//s.mListener, err = mdns.IntGroupDataProviderListener(s.mTransports)
	errors.SuccessOrExit(err)

	//dnssd.ResolverInstance().Init(udp_endpoint.UDPEndpoint{})

	s.mDnssd.SetSecuredPort(s.mOperationalServicePort)
	s.mDnssd.SetUnsecuredPort(s.mUserDirectedCommissioningPort)
	s.mDnssd.SetInterfaceId(s.mInterfaceId)

	if s.GetFabricTable() != nil {
		if s.GetFabricTable().FabricCount() != 0 {
			if costant.ConfigNetworkLayerBle {
				//TODO
				//如果Fabric不为零，那么设备已经被添加
				//可以在这里关闭蓝牙
			}
		}
	}

	//如果设备开启了自动配对模式，进入模式
	if costant.ChipDeviceConfigEnablePairingAutostart {
		s.GetFabricTable().DeleteAllFabrics()
		err = s.mCommissioningWindowManager.OpenBasicCommissioningWindow()
		errors.SuccessOrExit(err)
	}
	s.mDnssd.StartServer()

	return s
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
