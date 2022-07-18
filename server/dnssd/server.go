package dnssd

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/device_layer"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/platform"
	"github.com/galenliu/chip/server/dnssd/costants/commissioning_mode"
	"github.com/galenliu/chip/server/dnssd/costants/commssion_advertise_mode"
	"github.com/galenliu/chip/server/dnssd/costants/discovery"
	"github.com/galenliu/chip/server/dnssd/parameters"
	"github.com/galenliu/dnssd"
	"github.com/galenliu/dnssd/responders"
	"github.com/galenliu/gateway/pkg/errors"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/miekg/dns"
	"net"
	"strings"
	"time"
)

const kTimeoutCleared = 0
const kMdnsPort = 5353

type Advertiser interface {
	Shutdown()
	RemoveServices() error
	AddResponder(responder responders.RecordResponder) *responders.QueryResponderSettings
	RemoveRecords() error
	AdvertiseRecords(dnssd.BroadcastAdvertiseType) error
	FinalizeServiceUpdate() error
}

type FabricTable interface {
	FabricCount() int
	GetFabricInfos() []credentials.FabricInfo
}

func OnPlatformEvent(event *platform.ChipDeviceEvent) {

}

func OnPlatformEventWrapper(event *platform.ChipDeviceEvent, uint642 uint64) {
	OnPlatformEvent(event)
}

type Server struct {
	mSecuredPort                  uint16
	mAdvertiser                   Advertiser
	mUnsecuredPort                uint16
	mInterfaceId                  net.Interface
	mCommissioningModeProvider    CommissioningModeProvider
	mFabrics                      FabricTable
	mCommissionableInstanceName   uint64
	mCurrentCommissioningMode     uint8
	mEphemeralDiscriminator       *uint16
	mIsInitialized                bool
	mEmptyTextEntries             []string
	mExtendedDiscoveryTimeoutSecs time.Duration
	mExtendedDiscoveryExpiration  time.Duration
}

func NewServer() *Server {
	return &Server{}
}

func (s Server) Init() *Server {
	s.mEmptyTextEntries = append(s.mEmptyTextEntries, "=")
	return &s
}

func (s *Server) StartServer() {
	var mode = CommissioningMode.Disabled
	if s.mCommissioningModeProvider != nil {
		mode = s.mCommissioningModeProvider.GetCommissioningMode()
	}
	s.startServer(mode)
}

func (s *Server) startServer(mode uint8) {

	log.Infof("updating services using commissioning mode %d", mode)

	platform.PlatformMgr().AddEventHandler(OnPlatformEventWrapper, 0)

	//
	if !s.mIsInitialized {
		s.UpdateCommissionableInstanceName()
	}
	//使用UDPEndPointManager初始化一个Dnssd-Advertiser
	advertiser, err := dnssd.Advertiser{}.Init(nil, kMdnsPort)
	if err != nil {
		log.Infof("Failed to initialize advertiser: %s", err.Error())
	}

	s.mAdvertiser = advertiser
	err = s.removeServices()
	if err != nil {
		log.Infof("failed to remove advertised services: %s", err.Error())
	}

	err = s.AdvertiseOperational()
	if err != nil {
		log.Infof("Failed to advertise operational node: %s", err.Error())
	}

	if mode == CommissioningMode.Disabled {
		err = s.advertiseCommissionableNode(mode)
		if err != nil {
			log.Infof("Failed to advertise commissionable node: %s", err.Error())
		}
	}
	if config.ChipDeviceConfigEnableExtendedDiscovery {
		//if s.GetExtendedDiscoveryTimeoutSecs() != costant.ChipDeviceConfigDiscoveryDisabled {
		//	alwaysAdvertiseExtended := s.GetExtendedDiscoveryTimeoutSecs() == costant.ChipDeviceConfigDiscoveryNoTimeout
		//	if alwaysAdvertiseExtended || s.mCurrentCommissioningMode != CommissioningMode.Disabled || s.mExtendedDiscoveryExpiration != kTimeoutCleared {
		//		err := s.AdvertiseCommissionableNode(mode)
		//		if err != nil {
		//			log.Infof("failed to advertise extended commissionable node: %", err.Error())
		//		}
		//		if s.mExtendedDiscoveryExpiration == kTimeoutCleared {
		//			// set timeout
		//			s.ScheduleExtendedDiscoveryExpiration()
		//		}
		//	}
		//}
	}

	if config.ChipDeviceConfigEnableCommissionerDiscovery {
		err := s.advertiseCommissioner()
		if err != nil {
			log.Infof(err.Error())
		}
	}
	err = s.mAdvertiser.FinalizeServiceUpdate()
	if err != nil {
		log.Infof("failed to finalize service update: %s", err.Error())
	}
	s.mIsInitialized = true
}

func (s *Server) advertise(commissionAbleNode bool, mode uint8) error {

	advertiseParameters := parameters.CommissionAdvertisingParameters{}.Init()
	advertiseParameters.EnableIpV4(true)
	if commissionAbleNode {
		advertiseParameters.SetPort(s.GetSecuredPort())
		advertiseParameters.SetCommissionAdvertiseMode(CommssionAdvertiseMode.CommissionableNode)
	} else {
		advertiseParameters.SetPort(s.GetUnsecuredPort())
		advertiseParameters.SetCommissionAdvertiseMode(CommssionAdvertiseMode.Commissioner)
	}
	advertiseParameters.SetCommissioningMode(mode)

	//set  mac
	{
		mac := platform.ConfigurationMgr().GetPrimaryMACAddress()
		if mac != "" {
			advertiseParameters.SetMaC(mac)
		} else {
			log.Infof("failed to get primary mac address of device. Generating a random one.")
			advertiseParameters.SetMaC(util.GenerateMac())
		}
	}

	//Set device vendor id
	vid, err := platform.GetDeviceInstanceInfoProvider().GetVendorId()
	if err != nil {
		log.Infof("Vendor ID not known")
	} else {
		advertiseParameters.SetVendorId(vid)
	}

	// set  productId
	pid, err := platform.GetDeviceInstanceInfoProvider().GetProductId()
	if err != nil {
		log.Infof("Product ID not known")
	} else {
		advertiseParameters.SetProductId(pid)
	}

	// set discriminator
	var discriminator uint16 = 0
	discriminator, err = DeviceLayer.GetCommissionableDataProvider().GetSetupDiscriminator()
	if err != nil {
		log.Infof(
			"Setup discriminator read error: (%s )! Critical error, will not be commissionable.",
			err.Error())
	}
	if s.mEphemeralDiscriminator == nil {
		discriminator = *s.mEphemeralDiscriminator
	}
	advertiseParameters.SetShortDiscriminator(uint8(discriminator>>8) & 0x0F).
		SetLongDiscriminator(discriminator)

	// set device type id
	deviceTypeId, err := platform.ConfigurationMgr().GetDeviceTypeId()
	if platform.ConfigurationMgr().IsCommissionableDeviceTypeEnabled() && err == nil {
		if err != nil {
			advertiseParameters.SetDeviceType(deviceTypeId)
		}
	}

	//set device name
	if platform.ConfigurationMgr().IsCommissionableDeviceNameEnabled() {
		deviceName, err := platform.ConfigurationMgr().GetCommissionableDeviceName()
		if err != nil {
			advertiseParameters.SetDeviceName(deviceName)
		}
	}

	advertiseParameters.SetLocalMRPConfig(messageing.GetLocalMRPConfig()).SetTcpSupported(false)
	if !s.HaveOperationalCredentials() {
		value := platform.ConfigurationMgr().GetInitialPairingHint()
		if value != "" {
			advertiseParameters.SetPairingHint(value)
		} else {
			log.Infof("DNS-SD Pairing Hint not set")
		}
		ist := platform.ConfigurationMgr().GetInitialPairingInstruction()
		if ist == "" {
			log.Infof("DNS-SD Pairing Instruction not set")
		} else {
			advertiseParameters.SetPairingInstruction(ist)
		}
	} else {
		hint := platform.ConfigurationMgr().GetSecondaryPairingHint()
		if hint == "" {
			log.Infof("DNS-SD Pairing Hint not set")
		} else {
			advertiseParameters.SetPairingHint(hint)
		}

		ins := platform.ConfigurationMgr().GetSecondaryPairingInstruction()
		if ins == "" {
			log.Infof("DNS-SD Pairing Instruction not set")
		} else {
			advertiseParameters.SetPairingInstruction(ins)
		}
	}

	vid, _ = advertiseParameters.GetVendorId()
	pid, _ = advertiseParameters.GetProductId()
	log.Infof("advertise commission parameter vendorID=%d productID=%d discriminator=%04d/%02d",
		vid, pid,
		advertiseParameters.GetLongDiscriminator(), advertiseParameters.GetShortDiscriminator())

	return s.advCommission(advertiseParameters)
}

func (s *Server) AdvertiseOperational() error {
	if s.mFabrics == nil {
		return fmt.Errorf("fabrics nil")
	}
	for _, fabricInfo := range s.mFabrics.GetFabricInfos() {
		mac := platform.ConfigurationMgr().GetPrimaryMACAddress()
		if mac == "" {
			log.Infof("Failed to get primary mac address of device. Generating a random one.")
			mac = util.GenerateMac()
		}
		advertiseParameters := parameters.OperationalAdvertisingParameters{}.Init()
		advertiseParameters.SetPeerId(fabricInfo.GetPeerId())
		advertiseParameters.SetMaC(mac)
		advertiseParameters.SetPort(s.GetSecuredPort())
		advertiseParameters.SetInterfaceId(s.GetInterfaceId())
		advertiseParameters.SetLocalMRPConfig(messageing.GetLocalMRPConfig())
		advertiseParameters.SetTcpSupported(config.InetConfigEnableTcpEndpoint)
		advertiseParameters.EnableIpV4(true)

		log.Infof("advertise operational node %d - %d", advertiseParameters.GetPeerId().GetCompressedFabricId(),
			advertiseParameters.GetPeerId().GetNodeId())

		err := s.advOperational(advertiseParameters)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Server) advertiseCommissioner() error {
	return s.advertise(false, CommissioningMode.Disabled)
}

func (s *Server) GetCommissioningTxtEntries(params *parameters.CommissionAdvertisingParameters) []string {

	var txtFields []string
	pid, pErr := params.GetProductId()
	vid, vErr := params.GetVendorId()

	if pErr != nil && vErr != nil {
		txtFields = append(txtFields, fmt.Sprintf("VP=%d+%d", vid, pid))
	} else if vErr != nil {
		txtFields = append(txtFields, fmt.Sprintf("VP=%d", vid))
	}
	dType, err := params.GetDeviceType()
	if err != nil {
		txtFields = append(txtFields, fmt.Sprintf("DT=%d", dType))
	}
	if params.GetDeviceName() != "" {
		txtFields = append(txtFields, fmt.Sprintf("%s", params.GetDeviceName()))
	}

	commonTxt := s.AddCommonTxtEntries(*params.BaseAdvertisingParams)
	if commonTxt != nil && len(commonTxt) > 0 {
		txtFields = append(txtFields, commonTxt...)
	}

	if params.GetCommissionAdvertiseMode() == CommssionAdvertiseMode.CommissionableNode {
		name, _ := MakeServiceSubtype(discovery.LongDiscriminator, params.GetLongDiscriminator())
		txtFields = append(txtFields, name)
	}

	cm, err := MakeServiceSubtype(discovery.CommissioningMode, params.GetCommissioningMode())
	txtFields = append(txtFields, cm)

	if value := params.GetRotatingDeviceId(); value != "" {
		txtFields = append(txtFields, fmt.Sprintf("RI=%s", value))
	}

	if value := params.GetPairingHint(); value != "" {
		txtFields = append(txtFields, fmt.Sprintf("PH=%s", value))
	}

	if value := params.GetPairingInstruction(); value != "" {
		txtFields = append(txtFields, fmt.Sprintf("PI=%s", value))
	}
	return txtFields
}

func (s Server) GetOperationalTxtEntries(params *parameters.OperationalAdvertisingParameters) []string {
	txtFields := s.AddCommonTxtEntries(*params.BaseAdvertisingParams)
	if len(txtFields) == 0 || txtFields == nil {
		return s.mEmptyTextEntries
	}
	return txtFields
}

func (s *Server) AddCommonTxtEntries(params parameters.BaseAdvertisingParams) []string {
	var list []string

	if mrp := params.GetLocalMRPConfig(); mrp != nil {
		if mrp.IdleRetransTimeout > kMaxRetryInterval {
			log.Infof("MRP retry interval idle value exceeds allowed range of 1 hour, using maximum available")
			mrp.IdleRetransTimeout = kMaxRetryInterval
		}
		sleepyIdleIntervalBuf := fmt.Sprintf("SII=%d", mrp.IdleRetransTimeout)
		list = append(list, sleepyIdleIntervalBuf)

		if mrp.ActiveRetransTimeout > kMaxRetryInterval {
			log.Infof("MRP retry interval active value exceeds allowed range of 1 hour, using maximum available")
			mrp.ActiveRetransTimeout = kMaxRetryInterval
		}
		sleepyActiveIntervalBuf := fmt.Sprintf("SAI=%d", mrp.ActiveRetransTimeout)
		list = append(list, sleepyActiveIntervalBuf)
	}

	if value := params.GetTcpSupported(); value != nil {
		list = append(list, fmt.Sprintf("T=%d", func() int {
			if *value {
				return 1
			}
			return 0
		}()))
	}
	return list
}

func (s *Server) advertiseCommissionableNode(mode uint8) error {
	s.mCurrentCommissioningMode = mode
	return s.advertise(true, mode)
}

func (s Server) Shutdown() {
	s.advertiseRecords(dnssd.KRemovingAll)
	s.mAdvertiser.Shutdown()
	s.mIsInitialized = false
}

func (s *Server) advertiseRecords(advertiseType dnssd.BroadcastAdvertiseType) {
	_ = s.mAdvertiser.AdvertiseRecords(advertiseType)

}

func (s *Server) removeServices() error {
	return s.mAdvertiser.RemoveServices()
}

func (s *Server) advCommission(params *parameters.CommissionAdvertisingParameters) error {

	s.advertiseRecords(dnssd.KRemovingAll)

	var sType string
	if params.GetCommissionAdvertiseMode() == CommssionAdvertiseMode.CommissionableNode {
		sType = CommissionableServiceName
	} else {
		sType = CommissionerServiceName
	}

	serviceName := Fqdn(sType, CommissionProtocol, LocalDomain)
	name, _ := MakeServiceSubtype(discovery.InstanceName, s.GetCommissionableInstanceName())
	instanceName := Fqdn(name, sType, CommissionProtocol, LocalDomain)
	hostName := Fqdn(params.GetMac(), LocalDomain)

	if !s.mAdvertiser.AddResponder(responders.NewPtrResponder(serviceName, instanceName)).
		SetReportAdditional(instanceName).
		SetReportInServiceListing(true).
		IsValid() {
		return errors.New("failed to add service PTR record mDNS responder")
	}

	if !s.mAdvertiser.AddResponder(responders.NewSrvResponder(instanceName, hostName, params.GetPort())).
		SetReportAdditional(hostName).
		IsValid() {
		return errors.New("failed to add SRV record mDNS responder")
	}

	if !s.mAdvertiser.AddResponder(responders.NewIPv6Responder(hostName)).
		IsValid() {
		return errors.New("failed to add IPv6 mDNS responder")
	}

	if params.IsIPv4Enabled() {
		if !s.mAdvertiser.AddResponder(responders.NewIPv4Responder(hostName)).
			IsValid() {
			return errors.New("failed to add IPv6 mDNS responder")
		}
	}

	vid, err := params.GetVendorId()
	if err == nil {
		name, _ := MakeServiceSubtype(discovery.VendorId, vid)
		vendorServiceName := Fqdn(name, SubtypeServiceNamePart, sType, CommissionProtocol, LocalDomain)
		if !s.mAdvertiser.AddResponder(responders.NewPtrResponder(vendorServiceName, instanceName)).
			SetReportAdditional(instanceName).
			SetReportInServiceListing(true).
			IsValid() {
			return errors.New("failed to add vendor PTR record mDNS responder")
		}
	}

	dType, err := params.GetDeviceType()
	if err == nil {
		name, _ := MakeServiceSubtype(discovery.DeviceType, dType)
		typeServiceName := Fqdn(name, SubtypeServiceNamePart, sType, CommissionProtocol, LocalDomain)
		if !s.mAdvertiser.AddResponder(responders.NewPtrResponder(typeServiceName, instanceName)).
			SetReportAdditional(instanceName).
			SetReportInServiceListing(true).
			IsValid() {
			return errors.New("failed to add vendor PTR record mDNS responder")
		}
	}

	if params.GetCommissionAdvertiseMode() == CommssionAdvertiseMode.CommissionableNode {
		// TODO
	}

	if !s.mAdvertiser.AddResponder(responders.NewTxtResponder(instanceName, s.GetCommissioningTxtEntries(params))).
		SetReportAdditional(hostName).
		IsValid() {
		return errors.New("failed to add TXT record mDNS responder")
	}
	s.advertiseRecords(dnssd.KStarted)
	log.Infof("mDNS service published: %s", instanceName)
	return nil
}

func (s *Server) advOperational(params *parameters.OperationalAdvertisingParameters) error {

	s.advertiseRecords(dnssd.KRemovingAll)

	var name = params.GetPeerId().String()
	serviceName := Fqdn(OperationalServiceName, OperationalProtocol, LocalDomain)
	instanceName := Fqdn(name, OperationalServiceName, OperationalProtocol, LocalDomain)
	hostName := Fqdn(name, LocalDomain)

	if !s.mAdvertiser.AddResponder(responders.NewPtrResponder(serviceName, instanceName)).
		SetReportAdditional(instanceName).SetReportInServiceListing(true).IsValid() {
		return fmt.Errorf("failed to add service PTR record mDNS responder")
	}

	if !s.mAdvertiser.AddResponder(responders.NewSrvResponder(instanceName, hostName, params.GetPort())).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add SRV record mDNS responder")
	}

	if !s.mAdvertiser.AddResponder(responders.NewTxtResponder(instanceName, s.GetOperationalTxtEntries(params))).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add TXT record mDNS responder")
	}

	if !s.mAdvertiser.AddResponder(responders.NewIPv6Responder(hostName)).
		IsValid() {
		return fmt.Errorf("failed to add IPv6 mDNS responder")
	}

	if params.IsIPv4Enabled() {
		if !s.mAdvertiser.AddResponder(responders.NewIPv4Responder(hostName)).
			IsValid() {
			return fmt.Errorf("failed to add IPv4 mDNS responder")

		}
	}

	id := params.GetPeerId().GetCompressedFabricId()
	fabricId, _ := MakeServiceSubtype(discovery.CompressedFabricId, id)
	compressedFabricIdSubtype := Fqdn(fabricId, SubtypeServiceNamePart, OperationalServiceName, OperationalProtocol, LocalDomain)
	if !s.mAdvertiser.AddResponder(responders.NewPtrResponder(compressedFabricIdSubtype, instanceName)).
		SetReportAdditional(instanceName).
		IsValid() {
		log.Infof("Failed to add device type PTR record mDNS responder")
	}

	log.Infof("CHIP minimal mDNS configured as 'Operational device'.")
	s.advertiseRecords(dnssd.KStarted)
	log.Infof("mDNS service published: %s", instanceName)
	return nil
}

func (s *Server) AdvertiseCommissionableNode(mode uint8) error {
	if config.ChipDeviceConfigEnableExtendedDiscovery {
		s.mCurrentCommissioningMode = mode
		if mode == CommissioningMode.Disabled {
			s.HandleExtendedDiscoveryExpiration()
			s.mExtendedDiscoveryExpiration = kTimeoutCleared
		}
	}
	return s.advertise(true, mode)
}

func Fqdn(args ...string) string {
	var name = ""
	for _, arg := range args {
		name = name + dns.Fqdn(strings.TrimSpace(arg))
	}
	return dns.Fqdn(name)
}
