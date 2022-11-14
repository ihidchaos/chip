package dnssd

import (
	"fmt"
	"github.com/galenliu/chip/credentials"
	"github.com/galenliu/chip/pkg/dnssd/responder"
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
	AdvertiseOperational(params *OperationalAdvertisingParameters) error
	AdvertiseCommission(params *CommissionAdvertisingParameters) error
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

func (a *AdvertiseImpl) AdvertiseOperational(params *OperationalAdvertisingParameters) error {
	return a.advertiseOperational(params)
}

func NewAdvertise() *AdvertiseImpl {
	a := &AdvertiseImpl{}
	a.mQueryResponderAllocatorCommissioner = NewQueryResponderAllocator()
	a.mQueryResponderAllocatorCommissionable = NewQueryResponderAllocator()
	a.mResponseSender = NewResponseSender()
	a.mResponseSender.AddQueryResponder(a.mQueryResponderAllocatorCommissioner.GetQueryResponder())
	a.mResponseSender.AddQueryResponder(a.mQueryResponderAllocatorCommissionable.GetQueryResponder())
	a.mEmptyTextEntries = "="
	return a
}

func (a *AdvertiseImpl) Init() error {
	a.mMdnsServer = NewMdnsServerImpl()
	if !a.mInitialized {
		_ = a.UpdateCommissionableInstanceName()
	}
	a.mResponseSender.SetServer(a.mMdnsServer)
	a.mMdnsServer.SetHandler(a.mResponseSender)
	err := a.mMdnsServer.StartServer(MdnsPort)
	if err != nil {
		return err
	}
	err = a.AdvertiseRecords(AdvertiseTypeStarted)
	a.mInitialized = true
	return err
}

func (a *AdvertiseImpl) GetCommissionableInstanceName() (string, error) {
	return a.mCommissionableInstanceName, nil
}

func (a *AdvertiseImpl) UpdateCommissionableInstanceName() error {
	a.mCommissionableInstanceName = fmt.Sprintf("%016X", rand.Uint64())
	log.Infof("")
	return nil
}

func (a *AdvertiseImpl) AdvertiseCommission(params *CommissionAdvertisingParameters) error {

	_ = a.AdvertiseRecords(AdvertiseTypeRemovingAll)
	allocator := func() *QueryResponderAllocator {
		if params.GetCommissioningMode() == AdvertiseModeCommissionableNode {
			return a.mQueryResponderAllocatorCommissionable
		}
		return a.mQueryResponderAllocatorCommissioner
	}()

	serviceType := func() string {
		if params.GetCommissioningMode() == AdvertiseModeCommissionableNode {
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

	if !allocator.AddResponder(responder.NewIPv6Responder(hostName)).
		IsValid() {
		return fmt.Errorf("failed to add IPv6 mDNS responder")
	}

	if params.IsIPv4Enabled() {
		if !allocator.AddResponder(responder.NewIPv4Responder(hostName)).
			IsValid() {
			return fmt.Errorf("failed to add IPv4 mDNS responder")
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

	if params.GetCommissionAdvertiseMode() == AdvertiseModeCommissionableNode {
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

	if params.GetCommissioningMode() == CommissioningModeDisabled {
		if name := makeServiceSubtype[int](Subtype_CommissioningMode); name != "" {
			commissioningModeServiceName := Fqdn(name, KSubtypeServiceNamePart, serviceType, KCommissionProtocol, KLocalDomain)
			if !allocator.AddResponder(responder.NewPtrResponder(commissioningModeServiceName, instanceName)).
				SetReportAdditional(instanceName).SetReportInServiceListing(true).IsValid() {
				return fmt.Errorf("failed to add commissioning mode PTR record mDNS responder")
			}
		}
	}

	if !allocator.AddResponder(responder.NewTxtResponder(instanceName, a.GetCommissioningTxtEntries(params))).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add TXT record mDNS responder")
	}

	if params.GetCommissionAdvertiseMode() == AdvertiseModeCommissionableNode {
		log.Infof("CHIP minimal mDNS configured as 'Commissionable node device'.")
	} else {
		log.Infof("CHIP minimal mDNS configured as 'Commissioner device'.")
	}
	_ = a.AdvertiseRecords(AdvertiseTypeStarted)
	log.Infof("mDNS service published: %s", instanceName)
	return nil
}

func (a *AdvertiseImpl) advertiseOperational(params *OperationalAdvertisingParameters) error {

	var name = fmt.Sprintf("%016X-%016X", params.GetPeerId().GetCompressedFabricId(), params.GetPeerId().GetNodeId())

	_ = a.AdvertiseRecords(AdvertiseTypeRemovingAll)
	instanceName := Fqdn(name, KOperationalServiceName, KOperationalProtocol, KLocalDomain)

	operationalAllocator := a.FindOperationalAllocator(instanceName)
	if operationalAllocator == nil {
		operationalAllocator := a.FindEmptyOperationalAllocator()
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

	if !operationalAllocator.AddResponder(responder.NewTxtResponder(instanceName, a.GetOperationalTxtEntries(params))).
		SetReportAdditional(hostName).
		IsValid() {
		return fmt.Errorf("failed to add TXT record mDNS responder")
	}

	if !operationalAllocator.AddResponder(responder.NewIPv6Responder(hostName)).
		IsValid() {
		return fmt.Errorf("failed to add IPv6 mDNS responder")
	}

	if params.IsIPv4Enabled() {
		if !operationalAllocator.AddResponder(responder.NewIPv4Responder(hostName)).
			IsValid() {
			return fmt.Errorf("failed to add IPv4 mDNS responder")

		}
	}

	id := params.GetPeerId().GetCompressedFabricId()
	fabricId := makeServiceSubtype(Subtype_CompressedFabricId, id)
	compressedFabricIdSubtype := Fqdn(fabricId, KSubtypeServiceNamePart, KOperationalServiceName, KOperationalProtocol, KLocalDomain)
	if !operationalAllocator.AddResponder(responder.NewPtrResponder(compressedFabricIdSubtype, instanceName)).
		SetReportAdditional(instanceName).
		IsValid() {
		log.Infof("Failed to add device type PTR record mDNS responder")
	}

	log.Infof("CHIP minimal mDNS configured as 'Operational device'.")
	_ = a.AdvertiseRecords(AdvertiseTypeStarted)
	log.Infof("mDNS service published: %s", instanceName)
	return nil
}

func (a *AdvertiseImpl) FinalizeServiceUpdate() error {
	return nil
}

func (a *AdvertiseImpl) RemoveServices() error {
	return nil
}

func (a *AdvertiseImpl) GetCommissioningTxtEntries(params *CommissionAdvertisingParameters) []string {

	var txtFields []string

	// set vid and pid
	if vid, err := params.GetVendorId(); vid != 0 && err == nil {
		if pid, err := params.GetProductId(); pid != 0 && err == nil {
			txtFields = append(txtFields, fmt.Sprintf("VP=%a+%a", vid, pid))
		}
		txtFields = append(txtFields, fmt.Sprintf("VP=%a", vid))
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
	commonTxt := a.AddCommonTxtEntries(params.BaseAdvertisingParams)
	if commonTxt != nil && len(commonTxt) > 0 {
		txtFields = append(txtFields, commonTxt...)
	}

	if params.GetCommissionAdvertiseMode() == AdvertiseModeCommissionableNode {
		txtFields = append(txtFields, fmt.Sprintf("D=%a", params.GetLongDiscriminator()))
		txtFields = append(txtFields, fmt.Sprintf("CM=%a", params.GetCommissioningMode()))
	}

	if rid := params.GetRotatingDeviceId(); rid != "" {
		txtFields = append(txtFields, fmt.Sprintf("RI=%s", rid))
	}

	if ph := params.GetPairingHint(); ph != 0 {
		txtFields = append(txtFields, fmt.Sprintf("PH=%a", ph))
	}

	if pi := params.GetPairingInstruction(); pi != "" {
		txtFields = append(txtFields, fmt.Sprintf("PI=%s", pi))
	}

	return txtFields
}

func (a *AdvertiseImpl) GetOperationalTxtEntries(params *OperationalAdvertisingParameters) []string {
	txtFields := a.AddCommonTxtEntries(params.BaseAdvertisingParams)
	if len(txtFields) == 0 || txtFields == nil {
		return append(txtFields, a.mEmptyTextEntries)
	}
	return txtFields
}

func (a *AdvertiseImpl) AddCommonTxtEntries(params BaseAdvertisingParams) []string {

	var list []string
	if mrp := params.GetLocalMRPConfig(); mrp != nil {
		if mrp.IdleRetransTimeout > kMaxRetryInterval {
			log.Infof("MRP retry interval idle value exceeds allowed range of 1 hour, using maximum available")
			mrp.IdleRetransTimeout = kMaxRetryInterval
		}
		sleepyIdleIntervalBuf := fmt.Sprintf("SII=%a", mrp.IdleRetransTimeout)
		list = append(list, sleepyIdleIntervalBuf)

		if mrp.ActiveRetransTimeout > kMaxRetryInterval {
			log.Infof("MRP retry interval active value exceeds allowed range of 1 hour, using maximum available")
			mrp.ActiveRetransTimeout = kMaxRetryInterval
		}
		sleepyActiveIntervalBuf := fmt.Sprintf("SAI=%a", mrp.ActiveRetransTimeout)
		list = append(list, sleepyActiveIntervalBuf)
	}

	if value, err := params.GetTcpSupported(); err == nil {
		list = append(list, fmt.Sprintf("T=%a", func() int {
			if value {
				return 1
			}
			return 0
		}()))
	}
	return list
}

func (a *AdvertiseImpl) FindOperationalAllocator(name string) *QueryResponderAllocator {
	for _, allocator := range a.mOperationalResponders {
		r := allocator.GetResponder(dns.TypeSRV, name)
		if r != nil {
			return allocator
		}
	}
	return nil
}

func (a *AdvertiseImpl) FindEmptyOperationalAllocator() *QueryResponderAllocator {
	OperationalQueryAllocator := NewQueryResponderAllocator()
	a.mResponseSender.AddQueryResponder(OperationalQueryAllocator.GetQueryResponder())
	a.mOperationalResponders = append(a.mOperationalResponders, OperationalQueryAllocator)
	return OperationalQueryAllocator
}

func (a *AdvertiseImpl) AdvertiseRecords(typ int) error {

	var responseConfiguration = &responder.ResponseConfiguration{}
	if typ == AdvertiseTypeRemovingAll {
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
									err := a.mResponseSender.Respond(queryData, interfaceId)
									if err != nil {
										log.Errorf("failed to advertise records: %s", err.Error())
									}
								}
								if ip.Is6() {
									queryData.SetSrcAddr(netip.AddrPortFrom(ip, MdnsPort))
									queryData.SetDestAddr(netip.AddrPortFrom(IPv6LinkLocalMulticast, MdnsPort))
									err := a.mResponseSender.Respond(queryData, interfaceId)
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
