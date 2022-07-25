package dnssd

///*package DnssdServer
//
//import (
//	"fmt"
//	"github.com/galenliu/chip/config"
//	"github.com/galenliu/chip/credentials"
//	DeviceLayer "github.com/galenliu/chip/device_layer"
//	"github.com/galenliu/chip/lib"
//	"github.com/galenliu/chip/messageing"
//	"github.com/galenliu/chip/platform"
//	"github.com/galenliu/chip/server/dnssd/params"
//	"github.com/galenliu/gateway/pkg/util"
//	log "github.com/sirupsen/logrus"
//	"net"
//	"time"
//)
//
//type DnssdServer interface {
//	SetSecuredPort(port uint16)
//	GetSecuredPort() uint16
//	SetUnsecuredPort(port uint16)
//	GetUnsecuredPort() uint16
//	SetInterfaceId(interfaceId net.Interface)
//	GetInterfaceId() net.Interface
//	SetFabricTable(table FabricTable)
//	SetCommissioningModeProvider(provider CommissioningModeProvider)
//	//OnDiscoveryExpiration(aSystemLayer system.Layer, aAppState any)
//
//	AdvertiseOperational() error
//
//	StartServer()
//
//	StartServerMode(mode int)
//
//	GenerateRotatingDeviceId() (deviceId string, err error)
//
//	GetCommissionableInstanceName() (string, error)
//
//	SetEphemeralDiscriminator(discriminator uint16) error
//
//	advertise(commissionableNode bool, commissionMode int) error
//
//	advertiseCommissioner() error
//
//	advertiseCommissionableNode(commissionMode int) error
//
//	haveOperationalCredentials() bool
//}
//
//type ServerImpl struct {
//	mSecuredPort                  uint16
//	mUnsecuredPort                uint16
//	mInterfaceId                  net.Interface
//	mCommissioningModeProvider    CommissioningModeProvider
//	mFabricTable                  FabricTable
//	mCommissionableInstanceName   string
//	mCurrentCommissioningMode     int
//	mEphemeralDiscriminator       *uint16
//	mEmptyTextEntries             []string
//	mExtendedDiscoveryTimeoutSecs time.Duration
//	mExtendedDiscoveryExpiration  time.Duration
//
//	mAdvertiser MdnsServer
//
//	mQueryResponderAllocatorCommissionable QueryResponderAllocator
//	mQueryResponderAllocatorCommissioner   QueryResponderAllocator
//}
//
//type FabricTable interface {
//	FabricCount() int
//	GetFabricInfos() []credentials.FabricInfo
//}
//
//func OnPlatformEvent(event *platform.ChipDeviceEvent) {
//
//}
//
//func OnPlatformEventWrapper(event *platform.ChipDeviceEvent, uint642 uint64) {
//	OnPlatformEvent(event)
//}
//
//func NewServer() *ServerImpl {
//	return &ServerImpl{}
//}
//
//func (s *ServerImpl) Init() *ServerImpl {
//	s.mAdvertiser = MdnsInstance()
//	s.mEmptyTextEntries = append(s.mEmptyTextEntries, "=")
//	return s
//}
//
//func (s *ServerImpl) StartServer() {
//	var mode = CommissioningMode_Disabled
//	if s.mCommissioningModeProvider != nil {
//		mode = s.mCommissioningModeProvider.GetCommissioningMode()
//	}
//	s.StartServerMode(mode)
//}
//
//func (s *ServerImpl) StartServerMode(mode int) {
//
//	log.Infof("updating services using commissioning mode %d", mode)
//
//	platform.PlatformMgr().AddEventHandler(OnPlatformEventWrapper, 0)
//
//	err := s.mAdvertiser.Init()
//	if err != nil {
//		log.Infof("failed to initialize advertiser: %s", err.Error())
//	}
//
//	err = s.mAdvertiser.RemoveServices()
//	if err != nil {
//		log.Infof("failed to remove advertised services: %s", err.Error())
//	}
//
//	err = s.AdvertiseOperational()
//	if err != nil {
//		log.Infof("failed to advertise operational node: %s", err.Error())
//	}
//
//	if mode != CommissioningMode_Disabled {
//		err := s.AdvertiseCommissionableNode(mode)
//		if err != nil {
//			log.Infof("Failed to advertise commissionable node: %s", err.Error())
//		}
//	}
//
//	err = s.mAdvertiser.FinalizeServiceUpdate()
//	if err != nil {
//		log.Infof("failed to finalize service update: %s", err.Error())
//	}
//
//	//
//	//if !s.mIsInitialized {
//	//	s.UpdateCommissionableInstanceName()
//	//}
//	////使用UDPEndPointManager初始化一个Dnssd-DnssdServer
//	//advertiser, err := dnssd.DnssdServer{}.Init(nil, MdnsPort)
//	//if err != nil {
//	//	log.Infof("Failed to initialize advertiser: %s", err.Error())
//	//}
//	//
//	//err = s.removeServices()
//	//if err != nil {
//	//	log.Infof("failed to remove advertised services: %s", err.Error())
//	//}
//	//
//	//err = s.AdvertiseOperational()
//	//if err != nil {
//	//	log.Infof("Failed to advertise operational node: %s", err.Error())
//	//}
//	//
//	//if mode == CommissioningMode_Disabled {
//	//	err = s.advertiseCommissionableNode(mode)
//	//	if err != nil {
//	//		log.Infof("Failed to advertise commissionable node: %s", err.Error())
//	//	}
//	//}
//	//if config.ChipDeviceConfigEnableExtendedDiscovery {
//	//	//if s.GetExtendedDiscoveryTimeoutSecs() != costant.ChipDeviceConfigDiscoveryDisabled {
//	//	//	alwaysAdvertiseExtended := s.GetExtendedDiscoveryTimeoutSecs() == costant.ChipDeviceConfigDiscoveryNoTimeout
//	//	//	if alwaysAdvertiseExtended || s.mCurrentCommissioningMode != CommissioningMode.Disabled || s.mExtendedDiscoveryExpiration != kTimeoutCleared {
//	//	//		err := s.AdvertiseCommissionableNode(mode)
//	//	//		if err != nil {
//	//	//			log.Infof("failed to advertise extended commissionable node: %", err.Error())
//	//	//		}
//	//	//		if s.mExtendedDiscoveryExpiration == kTimeoutCleared {
//	//	//			// set timeout
//	//	//			s.ScheduleExtendedDiscoveryExpiration()
//	//	//		}
//	//	//	}
//	//	//}
//	//}
//	//
//	//if config.ChipDeviceConfigEnableCommissionerDiscovery {
//	//	err := s.advertiseCommissioner()
//	//	if err != nil {
//	//		log.Infof(err.Error())
//	//	}
//	//}
//	//err = s.mAdvertiser.FinalizeServiceUpdate()
//	//if err != nil {
//	//	log.Infof("failed to finalize service update: %s", err.Error())
//	//}
//
//}
//

//	vid, _ = advertiseParameters.GetVendorId()
//	pid, _ = advertiseParameters.GetProductId()
//	log.Infof("advertise commission parameter vendorID=%d productID=%d discriminator=%04X/%02X",
//		vid, pid,
//		advertiseParameters.GetLongDiscriminator(), advertiseParameters.GetShortDiscriminator())
//
//	return s.mAdvertiser.AdvertiseCommission(advertiseParameters)
//}
//
//func (s *ServerImpl) AdvertiseOperational() error {
//
//	if s.mFabricTable == nil {
//		return fmt.Errorf("fabrics nil")
//	}
//	for _, fabricInfo := range s.mFabricTable.GetFabricInfos() {
//
//		mac, err := platform.ConfigurationMgr().GetPrimaryMACAddress()
//		if err != nil {
//			log.Infof("Failed to get primary mac address of device. Generating a random one.")
//			mac = util.GenerateMac()
//		}
//
//		advertiseParameters := params.NewOperationalAdvertisingParameters()
//		advertiseParameters.SetPeerId(*fabricInfo.GetPeerId())
//		advertiseParameters.SetMaC(mac)
//		advertiseParameters.SetPort(s.GetSecuredPort())
//		advertiseParameters.SetInterfaceId(s.GetInterfaceId())
//		advertiseParameters.SetLocalMRPConfig(messageing.GetLocalMRPConfig())
//		advertiseParameters.SetTcpSupported(config.InetConfigEnableTcpEndpoint)
//		advertiseParameters.EnableIpV4(true)
//
//		log.Infof("advertise operational node %d - %d", advertiseParameters.GetPeerId().GetCompressedFabricId(),
//			advertiseParameters.GetPeerId().GetNodeId())
//
//		err = s.mAdvertiser.AdvertiseOperational(advertiseParameters)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (s ServerImpl) advertiseCommissioner() error {
//	return s.advertise(false, CommissioningMode_Disabled)
//}
//
////func (s *MdnsServerImpl) GetCommissioningTxtEntries(params *params.CommissionAdvertisingParameters) []string {
////
////	var txtFields []string
////	pid, pErr := params.GetProductId()
////	vid, vErr := params.GetVendorId()
////
////	if pErr != nil && vErr != nil {
////		txtFields = append(txtFields, fmt.Sprintf("VP=%d+%d", vid, pid))
////	} else if vErr != nil {
////		txtFields = append(txtFields, fmt.Sprintf("VP=%d", vid))
////	}
////	dType, err := params.GetDeviceType()
////	if err != nil {
////		txtFields = append(txtFields, fmt.Sprintf("DT=%d", dType))
////	}
////	if params.GetDeviceName() != "" {
////		txtFields = append(txtFields, fmt.Sprintf("%s", params.GetDeviceName()))
////	}
////
////	commonTxt := s.AddCommonTxtEntries(*params.BaseAdvertisingParams)
////	if commonTxt != nil && len(commonTxt) > 0 {
////		txtFields = append(txtFields, commonTxt...)
////	}
////
////	if params.GetCommissionAdvertiseMode() ==  AdvertiseMode_CommissionableNode {
////		name, _ := makeServiceSubtype(LongDiscriminator, params.GetLongDiscriminator())
////		txtFields = append(txtFields, name)
////	}
////
////	cm, err := makeServiceSubtype(CommissioningMode, params.GetCommissioningMode())
////	txtFields = append(txtFields, cm)
////
////	if value := params.GetRotatingDeviceId(); value != "" {
////		txtFields = append(txtFields, fmt.Sprintf("RI=%s", value))
////	}
////
////	if value := params.GetPairingHint(); value != 0 {
////		txtFields = append(txtFields, fmt.Sprintf("PH=%s", value))
////	}
////
////	if value := params.GetPairingInstruction(); value != "" {
////		txtFields = append(txtFields, fmt.Sprintf("PI=%s", value))
////	}
////	return txtFields
////}
//
////func (s MdnsServerImpl) GetOperationalTxtEntries(params *params.OperationalAdvertisingParameters) []string {
////	txtFields := s.AddCommonTxtEntries(*params.BaseAdvertisingParams)
////	if len(txtFields) == 0 || txtFields == nil {
////		return s.mEmptyTextEntries
////	}
////	return txtFields
////}
//
////func (s *MdnsServerImpl) AddCommonTxtEntries(params params.BaseAdvertisingParams) []string {
////	var list []string
////
////	if mrp := params.GetLocalMRPConfig(); mrp != nil {
////		if mrp.IdleRetransTimeout > kMaxRetryInterval {
////			log.Infof("MRP retry interval idle value exceeds allowed range of 1 hour, using maximum available")
////			mrp.IdleRetransTimeout = kMaxRetryInterval
////		}
////		sleepyIdleIntervalBuf := fmt.Sprintf("SII=%d", mrp.IdleRetransTimeout)
////		list = append(list, sleepyIdleIntervalBuf)
////
////		if mrp.ActiveRetransTimeout > kMaxRetryInterval {
////			log.Infof("MRP retry interval active value exceeds allowed range of 1 hour, using maximum available")
////			mrp.ActiveRetransTimeout = kMaxRetryInterval
////		}
////		sleepyActiveIntervalBuf := fmt.Sprintf("SAI=%d", mrp.ActiveRetransTimeout)
////		list = append(list, sleepyActiveIntervalBuf)
////	}
////
////	if value := params.GetTcpSupported(); value != nil {
////		list = append(list, fmt.Sprintf("T=%d", func() int {
////			if *value {
////				return 1
////			}
////			return 0
////		}()))
////	}
////	return list
////}
//
//func (s *ServerImpl) advertiseCommissionableNode(mode int) error {
//	s.mCurrentCommissioningMode = mode
//	return s.advertise(true, mode)
//}
//
////func (s MdnsServerImpl) Shutdown() {
////	s.advertiseRecords(dnssd.KRemovingAll)
////	s.mAdvertiser.Shutdown()
////	s.mIsInitialized = false
////}
//
//func (s *ServerImpl) removeServices() error {
//	return s.mAdvertiser.RemoveServices()
//}
//
//func (s *ServerImpl) AdvertiseCommissionableNode(mode int) error {
//	if config.ChipDeviceConfigEnableExtendedDiscovery {
//		s.mCurrentCommissioningMode = mode
//		if mode == CommissioningMode_Disabled {
//			//s.HandleExtendedDiscoveryExpiration()
//			// DeviceLayer::SystemLayer().CancelTimer(HandleExtendedDiscoveryExpiration, nullptr);
//			s.mExtendedDiscoveryExpiration = kTimeoutCleared
//		}
//	}
//	return s.advertise(true, mode)
//}
//
//func (s *ServerImpl) GenerateRotatingDeviceId() (deviceId string, err error) {
//	return "", nil
//}
//
//func (s *ServerImpl) SetEphemeralDiscriminator(discriminator uint16) error {
//	if discriminator >= DeviceLayer.KMaxDiscriminatorValue {
//		return lib.CHIP_ERROR_INVALID_ARGUMENT
//	}
//	s.mEphemeralDiscriminator = &discriminator
//	return nil
//}
//

//func (s *DnssdServerImpl) advertise(commissionAbleNode bool, mode int) error {
//
//
//}
