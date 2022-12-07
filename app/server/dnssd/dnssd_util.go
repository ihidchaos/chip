package dnssd

import (
	"fmt"
	"github.com/miekg/dns"
	"net/netip"
	"time"
)

type filterType uint8

const (
	KSubtypeServiceNamePart    = "_sub"
	KCommissionableServiceName = "_matterc"
	KCommissionerServiceName   = "_matterd"
	KOperationalServiceName    = "_matter"
	KCommissionProtocol        = "_udp"
	KLocalDomain               = "local"
	KOperationalProtocol       = "_tcp"
)

const kMaxRetryInterval = time.Millisecond * 3600000

type CommissioningModeProvider interface {
	GetCommissioningMode() int
}

// The mode of a Node in which it allows Commissioning.
const (
	CommissioningModeDisabled        = iota // Commissioning Mode is disabled, CM=0 in DNS-SD key/value pairs
	CommissioningModeEnableBasic            // Basic Commissioning Mode, CM=1 in DNS-SD key/value pairs
	CommissioningModeEnabledEnhanced        // Enhanced Commissioning Mode, CM=2 in DNS-SD key/value pairs

	AdvertiseModeCommissionableNode = 0
	AdvertiseModeCommissioner       = 1

	AdvertiseTypeRemovingAll = 0
	AdvertiseTypeStarted     = 1
)

type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~int
}

const (
	Subtype_None filterType = iota
	Subtype_ShortDiscriminator
	Subtype_LongDiscriminator
	Subtype_VendorId
	Subtype_DeviceType
	Subtype_CommissioningMode
	Subtype_InstanceName
	Subtype_Commissioner
	Subtype_CompressedFabricId
)

const MdnsPort uint16 = 5353
const MaxCommissionRecords = 20 // 11

var IPv4LinkLocalMulticast = netip.AddrFrom4([4]byte{224, 0, 0, 251})
var IPv6LinkLocalMulticast = netip.AddrFrom16([16]byte{0xFF, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFB})

type MdnsHandler interface {
	ServeMdns(ResponseWriter, *QueryData) error
}

type ResponseWriter interface {
	WriteMsg(*dns.Msg) error
}

type DnsResponseWriter struct {
	destAddr string
	clint    *dns.Client
}

func (d *DnsResponseWriter) NewDnsResponseWriter(addr, net string) *DnsResponseWriter {
	return &DnsResponseWriter{
		destAddr: addr,
		clint: &dns.Client{
			Net: net,
		},
	}
}

func (d *DnsResponseWriter) WriterMsg(msg *dns.Msg) error {
	clint := &dns.Client{
		Net:            "",
		UDPSize:        0,
		TLSConfig:      nil,
		Dialer:         nil,
		Timeout:        0,
		DialTimeout:    0,
		ReadTimeout:    0,
		WriteTimeout:   0,
		TsigSecret:     nil,
		TsigProvider:   nil,
		SingleInflight: false,
	}
	_, _, err := clint.Exchange(msg, d.destAddr)
	return err
}

func makeServiceSubtype[T Uint](filter filterType, values ...T) string {
	var val T = 0
	if len(values) != 0 {
		val = values[0]
	}
	switch filter {
	case Subtype_ShortDiscriminator:
		return fmt.Sprintf("_S%d", val)
	case Subtype_LongDiscriminator:
		return fmt.Sprintf("_L%d", val)
	case Subtype_VendorId:
		return fmt.Sprintf("_V%d", val)
	case Subtype_DeviceType:
		return fmt.Sprintf("_T%d", val)
	case Subtype_CommissioningMode:
		return "_CM"
	case Subtype_Commissioner:
		return fmt.Sprintf("_D%d", val)
	case Subtype_CompressedFabricId:
		return fmt.Sprintf("_I%016X", val)
	case Subtype_InstanceName:
		return fmt.Sprintf("%016X", val)
	case Subtype_None:
		return ""
	default:
		return fmt.Sprintf("%s", val)
	}
}
