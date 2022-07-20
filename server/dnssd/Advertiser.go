package DnssdServer

import (
	"fmt"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/server/dnssd/params"
	"github.com/galenliu/dnssd"
	"github.com/galenliu/dnssd/responder"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type ServiceAdvertiser interface {
	Init() error
	RemoveServices() error

	AdvertiseOperational(params *params.OperationalAdvertisingParameters) error
	AdvertiseCommission(params *params.CommissionAdvertisingParameters) error

	FinalizeServiceUpdate() error

	GetCommissionableInstanceName() (string, error)
	UpdateCommissionableInstanceName() error

	//Shutdown()
	//SetOperationalDelegate(delegate OperationalResolveDelegate)
	//SetCommissioningDelegate(delegate CommissioningResolveDelegate)
	//ResolveNodeId(peerId core.PeerId, isIpV6 bool)
	//DiscoverCommissionableNodes(filter DiscoveryFilter)
	//DiscoverCommissioners(filter DiscoveryFilter)

	AddResponder(responder responder.RecordResponder) *responder.QueryResponderSettings
	RemoveRecords() error

	AdvertiseRecords(mode int) error
}

type DiscoveryImplPlatform struct {
	Resolver
	mInitialized                           bool
	mResponseSender                        dnssd.ResponseSender
	mQueryResponderAllocatorCommissionable QueryResponderAllocator
	mQueryResponderAllocatorCommissioner   QueryResponderAllocator
	mOperationalResponders                 []*QueryResponderAllocator
	mCommissionableInstanceName            string
	mEmptyTextEntries                      string
}

var _serviceAdvertiserInstance *DiscoveryImplPlatform
var _serviceAdvertiserInstanceOnce sync.Once

func GetServiceAdvertiserInstance() *DiscoveryImplPlatform {
	_serviceAdvertiserInstanceOnce.Do(func() {
		if _serviceAdvertiserInstance == nil {
			_serviceAdvertiserInstance = &DiscoveryImplPlatform{}
			_serviceAdvertiserInstance.mEmptyTextEntries = "="
			_serviceAdvertiserInstance.mInitialized = true
		}
	})
	return _serviceAdvertiserInstance
}

func (d *DiscoveryImplPlatform) Init() error {
	if !d.mInitialized {
		return lib.CHIP_ERROR_INCORRECT_STATE
	}
	return nil
}

func (d *DiscoveryImplPlatform) RemoveServices() error {
	//TODO implement me
	panic("implement me")
}

func (d *DiscoveryImplPlatform) AdvertiseOperational(params *params.OperationalAdvertisingParameters) error {

	var name = params.GetPeerId().String()

	_ = d.AdvertiseRecords(BroadcastAdvertiseType_RemovingAll)
	instanceName := Fqdn(name, KOperationalServiceName, KOperationalProtocol, KLocalDomain)

	operationalAllocator := d.FindOperationalAllocator(instanceName)
	if operationalAllocator == nil {
		operationalAllocator := d.FindEmptyOperationalAllocator()
		if operationalAllocator == nil {
			return fmt.Errorf("failed to find an open operational allocator")
		}
	}

	serviceName := Fqdn(KOperationalServiceName, KOperationalProtocol, KLocalDomain)
	hostName := Fqdn(name, KLocalDomain)

	if !operationalAllocator.AddResponder(responder.NewPtrResponder(serviceName, instanceName)).
		SetReportAdditional(instanceName).
		SetReportInServiceListing(true).
		IsValid() {
		return fmt.Errorf("failed to add service PTR record mDNS responder")
	}

	if !operationalAllocator.AddResponder(responder.NewSrvResponder(instanceName, hostName, params.GetPort())).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add SRV record mDNS responder")
	}

	if !operationalAllocator.AddResponder(responder.NewTxtResponder(instanceName, d.GetOperationalTxtEntries(params))).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add TXT record mDNS responder")
	}

	if !operationalAllocator.AddResponder(responder.NewIPv6Responder(hostName, nil)).
		IsValid() {
		return fmt.Errorf("failed to add IPv6 mDNS responder")
	}

	if params.IsIPv4Enabled() {
		if !operationalAllocator.AddResponder(responder.NewIPv4Responder(hostName, nil)).
			IsValid() {
			return fmt.Errorf("failed to add IPv4 mDNS responder")

		}
	}

	id := params.GetPeerId().GetCompressedFabricId()
	fabricId := makeServiceSubtype(DiscoveryFilterType_kCompressedFabricId, id)
	compressedFabricIdSubtype := Fqdn(fabricId, KSubtypeServiceNamePart, KOperationalServiceName, KOperationalProtocol, KLocalDomain)
	if !operationalAllocator.AddResponder(responder.NewPtrResponder(compressedFabricIdSubtype, instanceName)).
		SetReportAdditional(instanceName).
		IsValid() {
		log.Infof("Failed to add device type PTR record mDNS responder")
	}

	log.Infof("CHIP minimal mDNS configured as 'Operational device'.")
	_ = d.AdvertiseRecords(BroadcastAdvertiseType_Started)
	log.Infof("mDNS service published: %s", instanceName)
	return nil
}

func (d *DiscoveryImplPlatform) AdvertiseCommission(params *params.CommissionAdvertisingParameters) error {

	_ = d.AdvertiseRecords(BroadcastAdvertiseType_RemovingAll)

	allocator := func() QueryResponderAllocator {
		if params.GetCommissioningMode() == AdvertiseMode_CommissionableNode {
			return d.mQueryResponderAllocatorCommissionable
		}
		return d.mQueryResponderAllocatorCommissioner
	}()

	serviceType := func() string {
		if params.GetCommissioningMode() == AdvertiseMode_CommissionableNode {
			return KCommissionableServiceName
		}
		return KCommissionerServiceName
	}()

	name, err := d.GetCommissionableInstanceName()
	if err != nil {
		return err
	}
	mac, err := params.GetMac()
	if err != nil {
		return err
	}

	serviceName := Fqdn(serviceType, KCommissionProtocol, KLocalDomain)
	instanceName := Fqdn(name, serviceType, KCommissionProtocol, KLocalDomain)
	hostName := Fqdn(mac, KLocalDomain)

	if !allocator.AddResponder(responder.NewPtrResponder(serviceName, instanceName)).
		SetReportAdditional(instanceName).
		SetReportInServiceListing(true).
		IsValid() {
		return fmt.Errorf("failed to add service PTR record mDNS responder")
	}

	if !allocator.AddResponder(responder.NewSrvResponder(instanceName, hostName, params.GetPort())).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add SRV record mDNS responder")
	}

	if !allocator.AddResponder(responder.NewIPv6Responder(hostName, nil)).
		IsValid() {
		return fmt.Errorf("failed to add IPv6 mDNS responder")
	}

	if params.IsIPv4Enabled() {
		if !!allocator.AddResponder(responder.NewIPv4Responder(hostName, nil)).
			IsValid() {
			return fmt.Errorf("failed to add IPv6 mDNS responder")
		}
	}

	if vid, err := params.GetVendorId(); err == nil && vid != 0 {
		name := makeServiceSubtype(Subtype_VendorId, vid)
		vendorServiceName := Fqdn(name, KSubtypeServiceNamePart, serviceType, KCommissionProtocol, KLocalDomain)
		if !allocator.AddResponder(responder.NewPtrResponder(vendorServiceName, instanceName)).
			SetReportAdditional(instanceName).
			SetReportInServiceListing(true).
			IsValid() {
			return fmt.Errorf("failed to add vendor PTR record mDNS responder")
		}
	}

	if deviceType, err := params.GetDeviceType(); err == nil && deviceType != 0 {
		name := makeServiceSubtype(Subtype_DeviceType, deviceType)
		typeServiceName := Fqdn(name, KSubtypeServiceNamePart, serviceType, KCommissionProtocol, KLocalDomain)
		if !allocator.AddResponder(responder.NewPtrResponder(typeServiceName, instanceName)).
			SetReportAdditional(instanceName).
			SetReportInServiceListing(true).
			IsValid() {
			return fmt.Errorf("failed to add vendor PTR record mDNS responder")
		}
	}

	if params.GetCommissionAdvertiseMode() == AdvertiseMode_CommissionableNode {
		if name := makeServiceSubtype(Subtype_ShortDiscriminator, params.GetShortDiscriminator()); name != "" {
			shortServiceName := Fqdn(name, KSubtypeServiceNamePart, serviceType, KCommissionProtocol, KLocalDomain)
			if !allocator.AddResponder(responder.NewPtrResponder(shortServiceName, instanceName)).
				SetReportAdditional(instanceName).
				SetReportInServiceListing(true).
				IsValid() {
				return fmt.Errorf("failed to add short discriminator PTR record mDNS responder")
			}
		}

		if name := makeServiceSubtype(Subtype_LongDiscriminator, params.GetLongDiscriminator()); name != "" {
			shortServiceName := Fqdn(name, KSubtypeServiceNamePart, serviceType, KCommissionProtocol, KLocalDomain)
			if !allocator.AddResponder(responder.NewPtrResponder(shortServiceName, instanceName)).
				SetReportAdditional(instanceName).
				SetReportInServiceListing(true).
				IsValid() {
				return fmt.Errorf("failed to add long discriminator PTR record mDNS responder")
			}
		}
	}

	if params.GetCommissioningMode() == CommissioningMode_Disabled {
		if name := makeServiceSubtype(Subtype_CommissioningMode); name != "" {
			commissioningModeServiceName := Fqdn(name, KSubtypeServiceNamePart, serviceType, KCommissionProtocol, KLocalDomain)
			if !allocator.AddResponder(responder.NewPtrResponder(commissioningModeServiceName, instanceName)).
				SetReportAdditional(instanceName).SetReportInServiceListing(true).IsValid() {
				return fmt.Errorf("failed to add commissioning mode PTR record mDNS responder")
			}
		}
	}

	if !allocator.AddResponder(responder.NewTxtResponder(instanceName, d.GetCommissioningTxtEntries(params))).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add TXT record mDNS responder")
	}
	_ = d.AdvertiseRecords(BroadcastAdvertiseType_Started)
	log.Infof("mDNS service published: %s", instanceName)
	return nil
}

func (d DiscoveryImplPlatform) FinalizeServiceUpdate() error {
	//TODO implement me
	panic("implement me")
}

func (d DiscoveryImplPlatform) GetCommissionableInstanceName() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DiscoveryImplPlatform) GetCommissioningTxtEntries(params *params.CommissionAdvertisingParameters) []string {

	var txtFields []string

	// set vid and pid
	if vid, err := params.GetVendorId(); vid != 0 && err == nil {
		if pid, err := params.GetProductId(); pid != 0 && err == nil {
			txtFields = append(txtFields, fmt.Sprintf("VP=%d+%d", vid, pid))
		}
		txtFields = append(txtFields, fmt.Sprintf("VP=%d", vid))
	}

	// set device type
	if dType, err := params.GetDeviceType(); err != nil && dType != 0 {
		txtFields = append(txtFields, fmt.Sprintf("DT=%x", dType))
	}

	// set device name
	if deviceType := params.GetDeviceName(); deviceType != "" {
		txtFields = append(txtFields, fmt.Sprintf("DN=%s", deviceType))
	}

	// set common txt
	commonTxt := d.AddCommonTxtEntries(params.BaseAdvertisingParams)
	if commonTxt != nil && len(commonTxt) > 0 {
		txtFields = append(txtFields, commonTxt...)
	}

	if params.GetCommissionAdvertiseMode() == AdvertiseMode_CommissionableNode {
		txtFields = append(txtFields, fmt.Sprintf("D=%d", params.GetLongDiscriminator()))
		txtFields = append(txtFields, fmt.Sprintf("CM=%d", params.GetCommissioningMode()))
	}

	if rid := params.GetRotatingDeviceId(); rid != "" {
		txtFields = append(txtFields, fmt.Sprintf("RI=%s", rid))
	}

	if ph := params.GetPairingHint(); ph != 0 {
		txtFields = append(txtFields, fmt.Sprintf("PH=%d", ph))
	}

	if pi := params.GetPairingInstruction(); pi != "" {
		txtFields = append(txtFields, fmt.Sprintf("PI=%s", pi))
	}

	return txtFields
}

func (d *DiscoveryImplPlatform) GetOperationalTxtEntries(params *params.OperationalAdvertisingParameters) []string {
	txtFields := d.AddCommonTxtEntries(params.BaseAdvertisingParams)
	if len(txtFields) == 0 || txtFields == nil {
		return append(txtFields, d.mEmptyTextEntries)
	}
	return txtFields
}

func (d *DiscoveryImplPlatform) AddCommonTxtEntries(params params.BaseAdvertisingParams) []string {

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

	if value, err := params.GetTcpSupported(); err == nil {
		list = append(list, fmt.Sprintf("T=%d", func() int {
			if value {
				return 1
			}
			return 0
		}()))
	}
	return list
}

func (d *DiscoveryImplPlatform) UpdateCommissionableInstanceName() error {
	//TODO implement me
	panic("implement me")
}

func (d *DiscoveryImplPlatform) AddResponder(responder responder.RecordResponder) *responder.QueryResponderSettings {
	//TODO implement me
	panic("implement me")
}

func (d *DiscoveryImplPlatform) RemoveRecords() error {
	//TODO implement me
	panic("implement me")
}

func (d *DiscoveryImplPlatform) AdvertiseRecords(mode int) error {
	return nil
}

func (d *DiscoveryImplPlatform) FindOperationalAllocator(name string) *QueryResponderAllocator {
	for _, allocator := range d.mOperationalResponders {
		r := allocator.GetResponder(dns.TypeSRV, name)
		if r != nil {
			return allocator
		}
	}
	return nil
}

func (d *DiscoveryImplPlatform) FindEmptyOperationalAllocator() *QueryResponderAllocator {
	OperationalQueryAllocator := NewQueryResponderAllocator()
	_ = d.mResponseSender.AddQueryResponder(OperationalQueryAllocator.GetQueryResponder())
	d.mOperationalResponders = append(d.mOperationalResponders, OperationalQueryAllocator)
	return OperationalQueryAllocator
}

func Fqdn(args ...string) string {
	var name = ""
	for _, arg := range args {
		name = name + dns.Fqdn(strings.TrimSpace(arg))
	}
	return dns.Fqdn(name)
}
