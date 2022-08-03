package dnssd

import (
	"fmt"
	"github.com/galenliu/chip/credentials"
	params2 "github.com/galenliu/chip/pkg/dnssd/params"
	responder2 "github.com/galenliu/chip/pkg/dnssd/responder"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"net/netip"
	"strings"
)

type Advertiser interface {
	Init() error
	RemoveServices() error
	GetCommissionableInstanceName() (string, error)
	UpdateCommissionableInstanceName() error
	advertiseOperational(params *params2.OperationalAdvertisingParameters) error
	advertiseCommission(params *params2.CommissionAdvertisingParameters) error
	FinalizeServiceUpdate() error
}

type AdvertiseImpl struct {
	mInitialized                           bool
	mResponseSender                        *ResponseSender
	mQueryResponderAllocatorCommissionable *QueryResponderAllocator
	mQueryResponderAllocatorCommissioner   *QueryResponderAllocator
	mOperationalResponders                 []*QueryResponderAllocator
	mCommissionableInstanceName            string
	mEmptyTextEntries                      string

	mSecuredPort                 uint16
	mUnsecuredPort               uint16
	mInterfaceId                 net.Interface
	mFabrice                     *credentials.FabricTable
	mCommissioningModeProvider   CommissioningModeProvider
	mCurrentCommissioningMode    int
	mExtendedDiscoveryExpiration any
	mEphemeralDiscriminator      *uint16
	mMdnsServer                  MdnsServer
}

func NewAdvertise() *AdvertiseImpl {
	a := &AdvertiseImpl{}
	a.mQueryResponderAllocatorCommissioner = NewQueryResponderAllocator()
	a.mQueryResponderAllocatorCommissionable = NewQueryResponderAllocator()
	a.mResponseSender = NewResponseSender()
	a.mEmptyTextEntries = "="
	return a
}

func (d *AdvertiseImpl) Init() error {
	d.mMdnsServer = NewMdnsServerImpl()
	d.mMdnsServer.Shutdown()
	if !d.mInitialized {
		_ = d.UpdateCommissionableInstanceName()
	}
	d.mResponseSender.SetServer(d.mMdnsServer)
	d.mMdnsServer.SetHandler(d.mResponseSender)
	err := d.mMdnsServer.StartServer(MdnsPort)
	if err != nil {
		return err
	}
	err = d.AdvertiseRecords(BroadcastAdvertiseType_Started)
	d.mInitialized = true
	return err
}

func (d *AdvertiseImpl) GetCommissionableInstanceName() (string, error) {
	return d.mCommissionableInstanceName, nil
}

func (d *AdvertiseImpl) UpdateCommissionableInstanceName() error {
	d.mCommissionableInstanceName = fmt.Sprintf("%016X", rand.Uint64())
	log.Infof("")
	return nil
}

func (a *AdvertiseImpl) advertiseCommission(params *params2.CommissionAdvertisingParameters) error {

	_ = a.AdvertiseRecords(BroadcastAdvertiseType_RemovingAll)
	allocator := func() *QueryResponderAllocator {
		if params.GetCommissioningMode() == AdvertiseMode_CommissionableNode {
			return a.mQueryResponderAllocatorCommissionable
		}
		return a.mQueryResponderAllocatorCommissioner
	}()

	serviceType := func() string {
		if params.GetCommissioningMode() == AdvertiseMode_CommissionableNode {
			return KCommissionableServiceName
		}
		return KCommissionerServiceName
	}()

	name, err := a.GetCommissionableInstanceName()
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

	if !allocator.AddResponder(responder2.NewPtrResponder(serviceName, instanceName)).
		SetReportAdditional(instanceName).
		SetReportInServiceListing(true).
		IsValid() {
		return fmt.Errorf("failed to add service PTR record mDNS responder")
	}

	if !allocator.AddResponder(responder2.NewSrvResponder(instanceName, hostName, params.GetPort())).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add SRV record mDNS responder")
	}

	if !allocator.AddResponder(responder2.NewIPv6Responder(hostName)).
		IsValid() {
		return fmt.Errorf("failed to add IPv6 mDNS responder")
	}

	if params.IsIPv4Enabled() {
		if !allocator.AddResponder(responder2.NewIPv4Responder(hostName)).
			IsValid() {
			return fmt.Errorf("failed to add IPv4 mDNS responder")
		}
	}

	if vid, err := params.GetVendorId(); err == nil && vid != 0 {
		name := makeServiceSubtype(Subtype_VendorId, vid)
		vendorServiceName := Fqdn(name, KSubtypeServiceNamePart, serviceType, KCommissionProtocol, KLocalDomain)
		if !allocator.AddResponder(responder2.NewPtrResponder(vendorServiceName, instanceName)).
			SetReportAdditional(instanceName).
			SetReportInServiceListing(true).
			IsValid() {
			return fmt.Errorf("failed to add vendor PTR record mDNS responder")
		}
	}

	if deviceType, err := params.GetDeviceType(); err == nil && deviceType != 0 {
		name := makeServiceSubtype(Subtype_DeviceType, deviceType)
		typeServiceName := Fqdn(name, KSubtypeServiceNamePart, serviceType, KCommissionProtocol, KLocalDomain)
		if !allocator.AddResponder(responder2.NewPtrResponder(typeServiceName, instanceName)).
			SetReportAdditional(instanceName).
			SetReportInServiceListing(true).
			IsValid() {
			return fmt.Errorf("failed to add vendor PTR record mDNS responder")
		}
	}

	if params.GetCommissionAdvertiseMode() == AdvertiseMode_CommissionableNode {
		if name := makeServiceSubtype(Subtype_ShortDiscriminator, params.GetShortDiscriminator()); name != "" {
			shortServiceName := Fqdn(name, KSubtypeServiceNamePart, serviceType, KCommissionProtocol, KLocalDomain)
			if !allocator.AddResponder(responder2.NewPtrResponder(shortServiceName, instanceName)).
				SetReportAdditional(instanceName).
				SetReportInServiceListing(true).
				IsValid() {
				return fmt.Errorf("failed to add short discriminator PTR record mDNS responder")
			}
		}

		if name := makeServiceSubtype(Subtype_LongDiscriminator, params.GetLongDiscriminator()); name != "" {
			shortServiceName := Fqdn(name, KSubtypeServiceNamePart, serviceType, KCommissionProtocol, KLocalDomain)
			if !allocator.AddResponder(responder2.NewPtrResponder(shortServiceName, instanceName)).
				SetReportAdditional(instanceName).
				SetReportInServiceListing(true).
				IsValid() {
				return fmt.Errorf("failed to add long discriminator PTR record mDNS responder")
			}
		}
	}

	if params.GetCommissioningMode() == CommissioningMode_Disabled {
		if name := makeServiceSubtype[int](Subtype_CommissioningMode); name != "" {
			commissioningModeServiceName := Fqdn(name, KSubtypeServiceNamePart, serviceType, KCommissionProtocol, KLocalDomain)
			if !allocator.AddResponder(responder2.NewPtrResponder(commissioningModeServiceName, instanceName)).
				SetReportAdditional(instanceName).SetReportInServiceListing(true).IsValid() {
				return fmt.Errorf("failed to add commissioning mode PTR record mDNS responder")
			}
		}
	}

	if !allocator.AddResponder(responder2.NewTxtResponder(instanceName, a.GetCommissioningTxtEntries(params))).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add TXT record mDNS responder")
	}

	if params.GetCommissionAdvertiseMode() == AdvertiseMode_CommissionableNode {
		log.Infof("CHIP minimal mDNS configured as 'Commissionable node device'.")
	} else {
		log.Infof("CHIP minimal mDNS configured as 'Commissioner device'.")
	}
	_ = a.AdvertiseRecords(BroadcastAdvertiseType_Started)
	log.Infof("mDNS service published: %s", instanceName)
	return nil
}

func (d *AdvertiseImpl) advertiseOperational(params *params2.OperationalAdvertisingParameters) error {

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

	if !operationalAllocator.AddResponder(responder2.NewPtrResponder(serviceName, instanceName)).
		SetReportAdditional(instanceName).
		SetReportInServiceListing(true).
		IsValid() {
		return fmt.Errorf("failed to add service PTR record mDNS responder")
	}

	if !operationalAllocator.AddResponder(responder2.NewSrvResponder(instanceName, hostName, params.GetPort())).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add SRV record mDNS responder")
	}

	if !operationalAllocator.AddResponder(responder2.NewTxtResponder(instanceName, d.GetOperationalTxtEntries(params))).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add TXT record mDNS responder")
	}

	if !operationalAllocator.AddResponder(responder2.NewIPv6Responder(hostName)).
		IsValid() {
		return fmt.Errorf("failed to add IPv6 mDNS responder")
	}

	if params.IsIPv4Enabled() {
		if !operationalAllocator.AddResponder(responder2.NewIPv4Responder(hostName)).
			IsValid() {
			return fmt.Errorf("failed to add IPv4 mDNS responder")

		}
	}

	id := params.GetPeerId().GetCompressedFabricId()
	fabricId := makeServiceSubtype(DiscoveryFilterType_kCompressedFabricId, id)
	compressedFabricIdSubtype := Fqdn(fabricId, KSubtypeServiceNamePart, KOperationalServiceName, KOperationalProtocol, KLocalDomain)
	if !operationalAllocator.AddResponder(responder2.NewPtrResponder(compressedFabricIdSubtype, instanceName)).
		SetReportAdditional(instanceName).
		IsValid() {
		log.Infof("Failed to add device type PTR record mDNS responder")
	}

	log.Infof("CHIP minimal mDNS configured as 'Operational device'.")
	_ = d.AdvertiseRecords(BroadcastAdvertiseType_Started)
	log.Infof("mDNS service published: %s", instanceName)
	return nil
}

func (d *AdvertiseImpl) FinalizeServiceUpdate() error {
	return nil
}

func (d *AdvertiseImpl) RemoveServices() error {
	return nil
}

func (d *AdvertiseImpl) GetCommissioningTxtEntries(params *params2.CommissionAdvertisingParameters) []string {

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

func (d *AdvertiseImpl) GetOperationalTxtEntries(params *params2.OperationalAdvertisingParameters) []string {
	txtFields := d.AddCommonTxtEntries(params.BaseAdvertisingParams)
	if len(txtFields) == 0 || txtFields == nil {
		return append(txtFields, d.mEmptyTextEntries)
	}
	return txtFields
}

func (d *AdvertiseImpl) AddCommonTxtEntries(params params2.BaseAdvertisingParams) []string {

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

func (d *AdvertiseImpl) FindOperationalAllocator(name string) *QueryResponderAllocator {
	for _, allocator := range d.mOperationalResponders {
		r := allocator.GetResponder(dns.TypeSRV, name)
		if r != nil {
			return allocator
		}
	}
	return nil
}

func (d *AdvertiseImpl) FindEmptyOperationalAllocator() *QueryResponderAllocator {
	OperationalQueryAllocator := NewQueryResponderAllocator()
	_ = d.mResponseSender.AddQueryResponder(OperationalQueryAllocator.GetQueryResponder())
	d.mOperationalResponders = append(d.mOperationalResponders, OperationalQueryAllocator)
	return OperationalQueryAllocator
}

func (d *AdvertiseImpl) AdvertiseRecords(typ int) error {

	var responseConfiguration = &responder2.ResponseConfiguration{}
	if typ == BroadcastAdvertiseType_RemovingAll {
		responseConfiguration.SetTtlSecondsOverride(0)
	}
	queryData := NewQueryData(dns.TypePTR, dns.ClassINET, false)
	queryData.configuration = responseConfiguration
	queryData.SetIsInternalBroadcast(true)

	interfaceIds, err := net.Interfaces()
	if err == nil {
		for _, interfaceId := range interfaceIds {
			adders, err := interfaceId.Addrs()
			if err == nil {
				for _, addr := range adders {
					if cidr, _, err := net.ParseCIDR(addr.String()); err == nil {
						if ip, err := netip.ParseAddr(cidr.String()); err == nil {
							if ip.IsGlobalUnicast() {
								if ip.Is4() {
									queryData.SetSrcAddr(netip.AddrPortFrom(ip, MdnsPort))
									queryData.SetDestAddr(netip.AddrPortFrom(IPv4LinkLocalMulticast, MdnsPort))
									err := d.mResponseSender.Respond(queryData, interfaceId)
									if err != nil {
										log.Errorf("failed to advertise records: %s", err.Error())
									}
								}
								if ip.Is6() {
									queryData.SetSrcAddr(netip.AddrPortFrom(ip, MdnsPort))
									queryData.SetDestAddr(netip.AddrPortFrom(IPv6LinkLocalMulticast, MdnsPort))
									err := d.mResponseSender.Respond(queryData, interfaceId)
									if err != nil {
										log.Errorf("failed to advertise records: %s", err.Error())
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func Fqdn(args ...string) string {
	var name = ""
	for _, arg := range args {
		name = name + dns.Fqdn(strings.TrimSpace(arg))
	}
	return dns.Fqdn(name)
}
