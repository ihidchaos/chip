package server

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/core"
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/pkg"
	"github.com/galenliu/chip/platform"
	"github.com/galenliu/chip/server/dnssd/costants/commissioning_mode"
	"github.com/galenliu/chip/server/dnssd/costants/commssion_advertise_mode"
	"github.com/galenliu/chip/server/dnssd/manager"
	"github.com/galenliu/chip/server/dnssd/parameters"
	"github.com/galenliu/dnssd"
	"github.com/galenliu/dnssd/QName"
	"github.com/galenliu/dnssd/chip"
	"github.com/galenliu/dnssd/record"
	"github.com/galenliu/dnssd/responders"
	"github.com/galenliu/gateway/pkg/errors"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"math/rand"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

const kMdnsPort = 5353

type Advertiser interface {
	Shutdown()
	Init([]netip.Addr, uint16) error
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

type DnssdServer struct {
	mSecuredPort                uint16
	mAdvertiser                 Advertiser
	mUnsecuredPort              uint16
	mInterfaceId                net.Interface
	mCommissioningModeProvider  *manager.CommissioningWindowManager
	mFabrics                    FabricTable
	mCommissionableInstanceName string
	mCurrentCommissioningMode   CommissioningMode.T
	mEphemeralDiscriminator     *uint16
	mIsInitialized              bool
}

func (s DnssdServer) Init() *DnssdServer {
	if !s.mIsInitialized {
		s.UpdateCommissionableInstanceName()
	}
	if s.mAdvertiser == nil {
		s.mAdvertiser = &dnssd.Advertiser{}
	}
	return &s
}

func (s *DnssdServer) StartServer() error {
	var mode = CommissioningMode.Disabled
	if s.mCommissioningModeProvider != nil {
		mode = s.mCommissioningModeProvider.GetCommissioningMode()
	}
	return s.startServer(mode)
}

func (s *DnssdServer) startServer(mode CommissioningMode.T) error {

	//使用UDPEndPointManager初始化一个Dnssd-Advertiser
	err := s.mAdvertiser.Init(nil, kMdnsPort)

	if err != nil {
		log.Info("failed initialize advertiser")
	}

	err = s.mAdvertiser.RemoveServices()
	errors.LogError(err, "Discover", "failed to remove advertised services")

	err = s.AdvertiseOperational()
	errors.LogError(err, "Discover", "Failed to advertise operational node")

	if mode == CommissioningMode.Disabled {
		err = s.AdvertiseCommissionableNode(mode)
		errors.LogError(err, "Discover", "Failed to advertise commissionable node")
	}

	// If any fabrics exist, the commissioning window must have been opened by the administrator
	// commissioning cluster commands which take care of the timeout.
	if !s.HaveOperationalCredentials() {
		s.ScheduleDiscoveryExpiration()
	}

	if config.ChipDeviceConfigEnableCommissionerDiscovery {
		err = s.AdvertiseCommissioner()
		errors.LogError(err, "Discover", "Failed to advertise commissioner")
	}

	err = s.mAdvertiser.FinalizeServiceUpdate()
	errors.LogError(err, "Discover", "failed to finalize service update")

	return nil
}

func (s *DnssdServer) Advertise(commissionAbleNode bool, mode CommissioningMode.T) error {

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
		mac, err := platform.ConfigurationMgr().GetPrimaryMACAddress()
		if err == nil {
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
	discriminator, err := platform.GetCommissionableDataProvider().GetSetupDiscriminator()
	if err != nil {
		log.Infof(
			"Setup discriminator read error: (%s )! Critical error, will not be commissionable.",
			err.Error())
	} else {
		advertiseParameters.SetProductId(pid)
	}
	if s.mEphemeralDiscriminator == nil {
		s.mEphemeralDiscriminator = &discriminator
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

	log.Infof("Advertise commission parameter vendorID=%d productID=%d discriminator=%04u/%02u",
		advertiseParameters.GetVendorId(), advertiseParameters.GetProductId(),
		advertiseParameters.GetLongDiscriminator(), advertiseParameters.GetShortDiscriminator())

	return s.advertiseCommission(advertiseParameters)
}

func (s *DnssdServer) AdvertiseOperational() error {
	if s.mFabrics == nil {
		return fmt.Errorf("fabrics nil")
	}
	for _, fabricInfo := range s.mFabrics.GetFabricInfos() {
		mac, err := platform.ConfigurationMgr().GetPrimaryMACAddress()
		if err != nil {
			log.Infof("Failed to get primary mac address of device. Generating a random one.")
		}
		advertiseParameters := parameters.OperationalAdvertisingParameters{}.Init()
		advertiseParameters.SetPeerId(fabricInfo.GetPeerId())
		advertiseParameters.SetMaC(mac)
		advertiseParameters.SetPort(s.GetSecuredPort())
		advertiseParameters.SetInterfaceId(s.GetInterfaceId())
		advertiseParameters.SetLocalMRPConfig(messageing.GetLocalMRPConfig())
		advertiseParameters.SetTcpSupported(config.InetConfigEnableTcpEndpoint)
		advertiseParameters.EnableIpV4(true)

		log.Infof("Advertise operational node %d - %d", advertiseParameters.GetPeerId().GetCompressedFabricId(),
			advertiseParameters.GetPeerId().GetNodeId())

		err = s.advertiseOperational(advertiseParameters)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s DnssdServer) AdvertiseCommissioner() error {
	return s.Advertise(false, CommissioningMode.Disabled)
}

func (s DnssdServer) HaveOperationalCredentials() bool {
	if s.mFabrics == nil {
		return false
	}
	return s.mFabrics.FabricCount() != 0
}

func (s *DnssdServer) SetCommissioningModeProvider(manager *manager.CommissioningWindowManager) {
	s.mCommissioningModeProvider = manager
}

func (s DnssdServer) ScheduleDiscoveryExpiration() {
	//TODO
	return
}

func (s *DnssdServer) UpdateCommissionableInstanceName() {
	s.mCommissionableInstanceName = strconv.FormatUint(rand.Uint64(), 16)
	s.mCommissionableInstanceName = strings.ToUpper(s.mCommissionableInstanceName)
}

func (s *DnssdServer) GetCommissionableInstanceName() string {
	if s.mCommissionableInstanceName == "" {
		s.mCommissionableInstanceName = strings.Replace(pkg.Mac48Address(pkg.RandHex()), ":", "", -1)
	}
	return s.mCommissionableInstanceName
}

func (s *DnssdServer) GetCommissioningTxtEntries(params *parameters.CommissionAdvertisingParameters) []string {

	var txtFields []string
	if params.GetProductId() != nil && params.GetVendorId() != nil {
		vid := params.GetVendorId()
		pid := params.GetVendorId()
		txtFields = append(txtFields, fmt.Sprintf("VP=%d+%d", *vid, *pid))
	} else if params.GetVendorId() != nil {
		txtFields = append(txtFields, fmt.Sprintf("VP=%d", params.GetVendorId()))
	}
	if params.GetDeviceType() != nil {
		txtFields = append(txtFields, fmt.Sprintf("DT=%d", *params.GetDeviceType()))
	}
	if params.GetDeviceName() != "" {
		txtFields = append(txtFields, fmt.Sprintf("%s", params.GetDeviceName()))
	}

	commonTxt := s.AddCommonTxtEntries(*params.BaseAdvertisingParams)
	if commonTxt != nil && len(commonTxt) > 0 {
		txtFields = append(txtFields, commonTxt...)
	}

	if params.GetCommissionAdvertiseMode() == CommssionAdvertiseMode.CommissionableNode {
		txtFields = append(txtFields, fmt.Sprintf("D=%d", params.GetLongDiscriminator()))
	}

	txtFields = append(txtFields, fmt.Sprintf("CM=%d", params.GetCommissioningMode()))

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

func (s DnssdServer) GetOperationalTxtEntries(params *parameters.OperationalAdvertisingParameters) []string {
	return nil
}

func (s *DnssdServer) AddCommonTxtEntries(params parameters.BaseAdvertisingParams) []string {
	var list []string

	if optionalMrp := params.GetLocalMRPConfig(); optionalMrp != nil {
		//TODO
	}

	if value := params.GetTcpSupported(); value != nil {
		i := 0
		if *value {
			i = 1
		}
		list = append(list, fmt.Sprintf("T=%d", i))
	}
	return list
}

func (s *DnssdServer) AdvertiseCommissionableNode(mode CommissioningMode.T) error {
	s.mCurrentCommissioningMode = mode
	return s.Advertise(true, mode)
}

func (s DnssdServer) Shutdown() {
	s.mAdvertiser.Shutdown()
}

func (s DnssdServer) SetFabricTable(f FabricTable) {
	s.mFabrics = f
}

func (s *DnssdServer) SetSecuredPort(port uint16) {
	s.mSecuredPort = port
}

func (s *DnssdServer) SetUnsecuredPort(port uint16) {
	s.mUnsecuredPort = port
}

func (s *DnssdServer) GetSecuredPort() uint16 {
	return s.mSecuredPort
}

func (s *DnssdServer) GetUnsecuredPort() uint16 {
	return s.mUnsecuredPort
}

func (s *DnssdServer) SetInterfaceId(inter net.Interface) {
	s.mInterfaceId = inter
}

func (s *DnssdServer) GetInterfaceId() net.Interface {
	return s.mInterfaceId
}

func (s *DnssdServer) advertiseOperational(params *parameters.OperationalAdvertisingParameters) error {
	_ = s.mAdvertiser.AdvertiseRecords(dnssd.KRemovingAll)
	var name = MakeInstanceName(params.GetPeerId())

	serviceName := QName.ParseFullQName(chip.KOperationalServiceName, chip.KOperationalProtocol, chip.KLocalDomain)
	instanceName := QName.ParseFullQName(name, chip.KOperationalServiceName, chip.KOperationalProtocol, chip.KLocalDomain)
	hostName := QName.ParseFullQName(name, chip.KLocalDomain)

	if !s.mAdvertiser.AddResponder(responders.NewPtrResponder(serviceName, instanceName)).
		SetReportAdditional(instanceName).SetReportInServiceListing(true).IsValid() {
		log.Infof("failed to add service PTR record mDNS responder")
	}

	if !s.mAdvertiser.AddResponder(responders.NewSrvResponder(record.NewSrvResourceRecord(instanceName, hostName, params.GetPort()))).
		SetReportAdditional(hostName).
		IsValid() {
		log.Infof("Failed to add SRV record mDNS responder")
	}

	if !s.mAdvertiser.AddResponder(responders.NewTxtResponder(instanceName, s.GetOperationalTxtEntries(params))).
		SetReportAdditional(hostName).
		IsValid() {
		log.Infof("Failed to add TXT record mDNS responder")
	}

	if !s.mAdvertiser.AddResponder(responders.NewIPv6Responder(hostName)).
		IsValid() {
		log.Infof("Failed to add IPv6 mDNS responder")
	}

	if params.IsIPv4Enabled() {
		if !s.mAdvertiser.AddResponder(responders.NewIPv4Responder(hostName)).
			IsValid() {
			log.Infof("Failed to add IPv4 mDNS responder")
		}
	}

	log.Infof("CHIP minimal mDNS configured as 'Operational device'.")
	_ = s.mAdvertiser.AdvertiseRecords(dnssd.KStarted)
	log.Infof("mDNS service published: %s", instanceName)
	return nil
}

func (s *DnssdServer) advertiseCommission(params *parameters.CommissionAdvertisingParameters) error {
	_ = s.mAdvertiser.RemoveRecords()

	var serviceType chip.ServiceType
	if params.GetCommissionAdvertiseMode() == CommssionAdvertiseMode.CommissionableNode {
		serviceType = chip.KCommissionableServiceName
	} else {
		serviceType = chip.KCommissionerServiceName
	}

	serviceName := QName.ParseFullQName(serviceType.String(), chip.KCommissionProtocol, chip.KLocalDomain)
	instanceName := QName.ParseFullQName(s.GetCommissionableInstanceName(), serviceType.String(), chip.KCommissionProtocol, chip.KLocalDomain)
	hostName := QName.ParseFullQName(chip.KLocalDomain, s.GetCommissionableInstanceName())

	if !s.mAdvertiser.AddResponder(responders.NewPtrResponder(serviceName, instanceName)).
		SetReportAdditional(instanceName).
		SetReportInServiceListing(true).
		IsValid() {
		return errors.New("failed to add service PTR record mDNS responder")
	}

	if !s.mAdvertiser.AddResponder(responders.NewSrvResponder(record.NewSrvResourceRecord(instanceName, hostName, params.GetPort()))).
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

	if params.GetVendorId() != nil {
		name := fmt.Sprintf("_V%d", *params.GetVendorId())
		vendorServiceName := QName.ParseFullQName(name, chip.KSubtypeServiceNamePart, serviceType.String(), chip.KCommissionProtocol, chip.KLocalDomain)
		if !s.mAdvertiser.AddResponder(responders.NewPtrResponder(vendorServiceName, instanceName)).
			SetReportAdditional(instanceName).
			SetReportInServiceListing(true).
			IsValid() {
			return errors.New("failed to add vendor PTR record mDNS responder")
		}
	}

	if params.GetDeviceType() != nil {
		name := fmt.Sprintf("_T%d", *params.GetDeviceType())
		typeServiceName := QName.ParseFullQName(name, chip.KSubtypeServiceNamePart, serviceType.String(), chip.KCommissionProtocol, chip.KLocalDomain)
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

	err := s.mAdvertiser.AdvertiseRecords(dnssd.KStarted)
	if err != nil {
		return err
	}
	log.Infof("mDNS service published: %s", instanceName.String())
	return nil
}

func MakeInstanceName(peerId *core.PeerId) string {
	nodeId := peerId.GetNodeId()               //uint64
	fabricId := peerId.GetCompressedFabricId() //uint64
	fabricIdH32 := uint32(fabricId >> 32)
	fabricIdL32 := uint32(fabricId)
	nodeIdH32 := uint32(nodeId >> 32)
	nodeIdL32 := uint32(nodeId)
	return fmt.Sprintf("%08x%08x%08x%08x", fabricIdH32, fabricIdL32, nodeIdH32, nodeIdL32)
}
