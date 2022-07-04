package server

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/inet"
	"github.com/galenliu/chip/inet/Interface"
	"github.com/galenliu/chip/inet/udp_endpoint"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/pkg"
	"github.com/galenliu/chip/platform"
	"github.com/galenliu/chip/server/parameters"
	"github.com/galenliu/dnssd"
	"github.com/galenliu/dnssd/chip"
	"github.com/galenliu/dnssd/core"
	"github.com/galenliu/dnssd/record"
	"github.com/galenliu/dnssd/responders"
	"github.com/galenliu/gateway/pkg/errors"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

type Advertiser interface {
	Shutdown()
	RemoveServices() error
	AddResponder(responder responders.RecordResponder) *responders.QueryResponderSettings
	RemoveRecords() error
	AdvertiseRecords() error
}

type Fabrics interface {
	FabricCount() int
}

type MDNS struct {
	mSecuredPort                uint16
	mAdvertiser                 Advertiser
	mUnsecuredPort              uint16
	mInterfaceId                Interface.Id
	mCommissioningModeProvider  *CommissioningWindowManager
	mFabrics                    Fabrics
	mCommissionableInstanceName string
	mCurrentCommissioningMode   config.CommissioningMode
}

var insDnssd *MDNS
var onceDnssd sync.Once

func DnssdInstance() *MDNS {
	onceDnssd.Do(func() {
		insDnssd = DefaultDnssd()
	})
	return insDnssd
}

func DefaultDnssd() *MDNS {
	i := &MDNS{}
	i.mAdvertiser = dnssd.GetAdvertiserInstance()
	return i
}

func (m MDNS) Shutdown() {
	m.mAdvertiser.Shutdown()
}

func (m MDNS) SetFabricTable(f Fabrics) {
	m.mFabrics = f
}

func (m *MDNS) SetSecuredPort(port uint16) {
	m.mSecuredPort = port
}

func (m *MDNS) SetUnsecuredPort(port uint16) {
	m.mUnsecuredPort = port
}

func (m *MDNS) GetSecuredPort() uint16 {
	return m.mSecuredPort
}

func (m *MDNS) GetUnsecuredPort() uint16 {
	return m.mUnsecuredPort
}

func (m *MDNS) SetInterfaceId(inter Interface.Id) {
	m.mInterfaceId = inter
}

func (m *MDNS) StartServer() error {
	var mode = config.KDisabled
	if m.mCommissioningModeProvider != nil {
		mode = m.mCommissioningModeProvider.GetCommissioningMode()
	}
	return m.startServer(mode)
}

func (m *MDNS) startServer(mode config.CommissioningMode) error {

	//使用UDPEndPointManager初始化一个Dnssd-Advertiser
	err := dnssd.GetAdvertiserInstance().Init(udp_endpoint.UDPEndpoint{})

	if err != nil {
		log.Info("failed initialize advertiser")
	}

	err = dnssd.GetAdvertiserInstance().RemoveServices()
	errors.LogError(err, "Discover", "failed to remove advertised services")

	err = m.AdvertiseOperational()
	errors.LogError(err, "Discover", "Failed to advertise operational node")

	if mode != config.KDisabled {
		err = m.AdvertiseCommissionableNode(mode)
		errors.LogError(err, "Discover", "Failed to advertise commissionable node")
	}

	// If any fabrics exist, the commissioning window must have been opened by the administrator
	// commissioning cluster commands which take care of the timeout.
	if !m.HaveOperationalCredentials() {
		m.ScheduleDiscoveryExpiration()
	}

	if config.ChipDeviceConfigEnableCommissionerDiscovery {
		err = m.AdvertiseCommissioner()
		errors.LogError(err, "Discover", "Failed to advertise commissioner")
	}

	err = dnssd.GetAdvertiserInstance().FinalizeServiceUpdate()
	errors.LogError(err, "Discover", "failed to finalize service update")

	return nil
}

func (m MDNS) AdvertiseOperational() error {

	return nil
}

func (m MDNS) AdvertiseCommissioner() error {
	return m.Advertise(false, config.KDisabled)
}

func (m MDNS) HaveOperationalCredentials() bool {
	if m.mFabrics == nil {
		return false
	}
	return m.mFabrics.FabricCount() != 0
}

func (m MDNS) Advertise(commissionAbleNode bool, mode config.CommissioningMode) error {

	advertiseParameters := parameters.CommissionAdvertisingParameters{}

	advertiseParameters.SetPort(util.ConditionFunc(commissionAbleNode, m.GetUnsecuredPort, m.GetUnsecuredPort))
	advertiseParameters.SetCommissionAdvertiseMode(util.ConditionValue(commissionAbleNode, config.KCommissionableNode, config.KCommissioner))

	advertiseParameters.SetCommissioningMode(mode)

	mac, err := platform.ConfigurationMgr().GetPrimaryMACAddress()
	errors.LogError(err, "Discovery", "Failed to get primary mac address of device. Generating a random one.")
	advertiseParameters.SetMaC(string(util.ConditionValue(err != nil, mac, util.GenerateMac())))

	vid, err := platform.ConfigurationMgr().GetVendorId()
	errors.LogError(err, "Discovery", "Vendor ID not known")
	if err != nil {
		advertiseParameters.SetVendorId(uint16(vid))
	}

	pid, err := platform.ConfigurationMgr().GetProductId()
	errors.LogError(err, "Discovery", "Product ID not known")
	if err != nil {
		advertiseParameters.SetProductId(uint16(pid))
	}

	//uint16_t discriminator = 0;
	//CHIP_ERROR error       = DeviceLayer::GetCommissionableDataProvider()->GetSetupDiscriminator(discriminator);
	//if (error != CHIP_NO_ERROR)
	//{
	//	ChipLogError(Discovery,
	//		"Setup discriminator read error (%" CHIP_ERROR_FORMAT ")! Critical error, will not be commissionable.",
	//	error.Format());
	//	return error;
	//}
	//
	// Override discriminator with temporary one if one is set
	//discriminator = mEphemeralDiscriminator.ValueOr(discriminator);
	//
	//advertiseParameters.SetShortDiscriminator(static_cast<uint8_t>((discriminator >> 8) & 0x0F))
	//.SetLongDiscriminator(discriminator);
	//

	if platform.ConfigurationMgr().IsCommissionableDeviceTypeEnabled() {
		did, err := platform.ConfigurationMgr().GetDeviceTypeId()
		if err != nil {
			advertiseParameters.SetDeviceType(int32(did))
		}
	}

	if platform.ConfigurationMgr().IsCommissionableDeviceNameEnabled() {
		name, err := platform.ConfigurationMgr().GetCommissionableDeviceName()
		if err != nil {
			advertiseParameters.SetDeviceName(name)
		}
	}

	advertiseParameters.SetMRPConfig(messageing.GetLocalMRPConfig())
	advertiseParameters.SetTcpSupported(inet.InetConfigEnableTcpEndpoint)

	if !m.HaveOperationalCredentials() {
		value := platform.ConfigurationMgr().GetInitialPairingHint()

		advertiseParameters.SetPairingHint(value)

		ist := platform.ConfigurationMgr().GetInitialPairingInstruction()
		advertiseParameters.SetPairingInstruction(ist)
	} else {
		hint := platform.ConfigurationMgr().GetSecondaryPairingHint()
		advertiseParameters.SetPairingHint(hint)

		ins := platform.ConfigurationMgr().GetSecondaryPairingInstruction()
		advertiseParameters.SetPairingInstruction(ins)
	}

	return m.AdvertiseCommission(advertiseParameters)
}

func (m *MDNS) SetCommissioningModeProvider(manager *CommissioningWindowManager) {
	m.mCommissioningModeProvider = manager
}

func (m MDNS) ScheduleDiscoveryExpiration() {
	//TODO
	return
}

func (m *MDNS) UpdateCommissionableInstanceName() {
	m.mCommissionableInstanceName = strconv.FormatUint(rand.Uint64(), 16)
	m.mCommissionableInstanceName = strings.ToUpper(m.mCommissionableInstanceName)
}

func (m *MDNS) GetCommissionableInstanceName() string {
	if m.mCommissionableInstanceName == "" {
		m.mCommissionableInstanceName = strings.Replace(pkg.Mac48Address(pkg.RandHex()), ":", "", -1)
	}
	return m.mCommissionableInstanceName
}

func (m *MDNS) AdvertiseCommission(params parameters.CommissionAdvertisingParameters) error {
	_ = m.mAdvertiser.RemoveRecords()

	var serviceType chip.ServiceType

	if params.GetCommissionAdvertiseMode() == config.KCommissionableNode {
		serviceType = chip.KCommissionableServiceName
	} else {
		serviceType = chip.KCommissionerServiceName
	}

	serviceName := core.ParseFullQName(serviceType.String(), chip.KCommissionProtocol, chip.KLocalDomain)
	instanceName := core.ParseFullQName(m.GetCommissionableInstanceName(), serviceType.String(), chip.KCommissionProtocol, chip.KLocalDomain)

	hostName := core.ParseFullQName(chip.KLocalDomain, m.GetCommissionableInstanceName())

	if !m.mAdvertiser.AddResponder(responders.NewPtrResponder(serviceName, instanceName)).
		SetReportAdditional(instanceName).
		SetReportInServiceListing(true).
		IsValid() {
		return errors.New("failed to add service PTR record mDNS responder")
	}

	if !m.mAdvertiser.AddResponder(responders.NewSrvResponder(record.NewSrvResourceRecord(instanceName, hostName, params.GetPort()))).
		SetReportAdditional(hostName).
		IsValid() {
		return errors.New("failed to add SRV record mDNS responder")
	}

	if !m.mAdvertiser.AddResponder(responders.NewIPv6Responder(hostName)).
		IsValid() {
		return errors.New("failed to add IPv6 mDNS responder")
	}

	if params.IsIPv4Enabled() {
		if !m.mAdvertiser.AddResponder(responders.NewIPv4Responder(hostName)).
			IsValid() {
			return errors.New("failed to add IPv6 mDNS responder")
		}
	}

	if params.GetVendorId() != nil {
		name := fmt.Sprintf("_V%m", *params.GetVendorId())
		vendorServiceName := core.ParseFullQName(name, chip.KSubtypeServiceNamePart, serviceType.String(), chip.KCommissionProtocol, chip.KLocalDomain)
		if !m.mAdvertiser.AddResponder(responders.NewPtrResponder(vendorServiceName, instanceName)).
			SetReportAdditional(instanceName).
			SetReportInServiceListing(true).
			IsValid() {
			return errors.New("failed to add vendor PTR record mDNS responder")
		}
	}

	if params.GetDeviceType() != nil {
		name := fmt.Sprintf("_T%m", *params.GetDeviceType())
		typeServiceName := core.ParseFullQName(name, chip.KSubtypeServiceNamePart, serviceType.String(), chip.KCommissionProtocol, chip.KLocalDomain)
		if !m.mAdvertiser.AddResponder(responders.NewPtrResponder(typeServiceName, instanceName)).
			SetReportAdditional(instanceName).
			SetReportInServiceListing(true).
			IsValid() {
			return errors.New("failed to add vendor PTR record mDNS responder")
		}
	}

	if params.GetCommissionAdvertiseMode() == config.KCommissionableNode {
		// TODO
	}

	if !m.mAdvertiser.AddResponder(responders.NewTxtResponder(instanceName, m.GetCommissioningTxtEntries(params))).
		SetReportAdditional(hostName).
		IsValid() {
		return errors.New("failed to add TXT record mDNS responder")
	}

	err := m.mAdvertiser.AdvertiseRecords()
	if err != nil {
		return err
	}
	log.Infof("mDNS service published: %a", instanceName.String())
	return nil
}

func (m *MDNS) GetCommissioningTxtEntries(params parameters.CommissionAdvertisingParameters) []string {

	var txtFields []string
	if params.GetProductId() != nil && params.GetVendorId() != nil {
		txtFields = append(txtFields, fmt.Sprintf("VP=%d+%d", *params.GetVendorId(), *params.GetProductId()))
	} else if params.GetVendorId() != nil {
		txtFields = append(txtFields, fmt.Sprintf("VP=%d", params.GetVendorId()))
	}
	if params.GetDeviceType() != nil {
		txtFields = append(txtFields, fmt.Sprintf("DT=%d", *params.GetDeviceType()))
	}
	if params.GetDeviceName() != "" {
		txtFields = append(txtFields, fmt.Sprintf("%m", params.GetDeviceName()))
	}

	commonTxt := m.AddCommonTxtEntries(*params.BaseAdvertisingParams)
	if commonTxt != nil && len(commonTxt) > 0 {
		txtFields = append(txtFields, commonTxt...)
	}

	if params.GetCommissionAdvertiseMode() == config.KCommissionableNode {
		txtFields = append(txtFields, fmt.Sprintf("D=%d", params.GetLongDiscriminator()))
	}

	txtFields = append(txtFields, fmt.Sprintf("CM=%d", params.GetCommissioningMode()))

	if value := params.GetRotatingDeviceId(); value != "" {
		txtFields = append(txtFields, fmt.Sprintf("RI=%m", value))
	}

	if value := params.GetPairingHint(); value != nil {
		txtFields = append(txtFields, fmt.Sprintf("PH=%d", *value))
	}

	if value := params.GetPairingInstruction(); value != "" {
		txtFields = append(txtFields, fmt.Sprintf("PI=%m", value))
	}
	return txtFields
}

func (m *MDNS) AddCommonTxtEntries(params parameters.BaseAdvertisingParams) []string {
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

func (m *MDNS) AdvertiseCommissionableNode(mode config.CommissioningMode) error {
	m.mCurrentCommissioningMode = mode
	return m.Advertise(true, mode)
}
