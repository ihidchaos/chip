package dnssd

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials"
	DeviceLayer "github.com/galenliu/chip/device_layer"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/platform"
	"github.com/galenliu/chip/server/dnssd/params"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"net/netip"
	"sync"
)

const MdnsPort uint16 = 5353
const MaxCommissionRecords = 20 // 11

var IPv4LinkLocalMulticast = netip.AddrFrom4([4]byte{224, 0, 0, 251})
var IPv6LinkLocalMulticast = netip.AddrFrom16([16]byte{0xFF, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFB})

type DnssdServer interface {
	SetSecuredPort(port uint16)
	GetSecuredPort() uint16
	SetUnsecuredPort(port uint16)
	GetUnsecuredPort() uint16

	SetInterfaceId(id net.Interface)
	GetInterfaceId() net.Interface

	SetFabricTable(fabrics *credentials.FabricTable)
	SetCommissioningModeProvider(manager CommissioningModeProvider)

	AdvertiseOperational() error
	StartServer()

	GetCommissionableInstanceName() string
	SetEphemeralDiscriminator(discriminator uint16) error

	Advertise(commissionableNode bool, commissionMode int) error
	AdvertiseCommissioner() error
	AdvertiseCommissionableNode(commissionMode int) error
}

type DnssdServerImpl struct {
	mSecuredPort                 uint16
	mUnsecuredPort               uint16
	mInterfaceId                 net.Interface
	mFabricTable                 *credentials.FabricTable
	mCommissioningModeProvider   CommissioningModeProvider
	mCurrentCommissioningMode    int
	mExtendedDiscoveryExpiration any
	mEphemeralDiscriminator      *uint16
	mdnsAdvertiser               Advertiser
}

var _serviceAdvertiserInstance *DnssdServerImpl
var _serviceAdvertiserInstanceOnce sync.Once

func GetInstance() *DnssdServerImpl {
	_serviceAdvertiserInstanceOnce.Do(func() {
		if _serviceAdvertiserInstance == nil {
			_serviceAdvertiserInstance = &DnssdServerImpl{
				mSecuredPort:                 0,
				mUnsecuredPort:               0,
				mInterfaceId:                 net.Interface{},
				mFabricTable:                 nil,
				mCommissioningModeProvider:   nil,
				mCurrentCommissioningMode:    0,
				mExtendedDiscoveryExpiration: nil,
				mEphemeralDiscriminator:      nil,
				mdnsAdvertiser:               NewAdvertise(),
			}
		}
	})
	return _serviceAdvertiserInstance
}

func NewDnssdServer() *DnssdServerImpl {
	return GetInstance()
}

func (d *DnssdServerImpl) GetSecuredPort() uint16 {
	return d.mSecuredPort
}

func (d *DnssdServerImpl) GetUnsecuredPort() uint16 {
	return d.mUnsecuredPort
}

func (d *DnssdServerImpl) GetInterfaceId() net.Interface {
	return d.mInterfaceId
}

func (d *DnssdServerImpl) AdvertiseOperational() error {
	if d.mFabricTable == nil {
		return lib.ChipErrorIncorrectState
	}
	for _, fabricInfo := range d.mFabricTable.GetFabricInfos() {
		mac, err := config.ConfigurationMgr().GetPrimaryMACAddress()
		if mac == "" || err != nil {
			mac = fmt.Sprintf("%016X", rand.Uint64())
		}
		advertiseParameters := params.NewOperationalAdvertisingParameters()
		advertiseParameters.SetPeerId(fabricInfo.GetPeerId())
		advertiseParameters.SetMaC(mac)
		advertiseParameters.SetPort(d.GetSecuredPort())
		advertiseParameters.SetInterfaceId(d.GetInterfaceId())
		advertiseParameters.SetLocalMRPConfig(messageing.GetLocalMRPConfig())
		advertiseParameters.SetTcpSupported(config.InetConfigEnableTcpEndpoint)
		advertiseParameters.EnableIpV4(true)
		if d.mdnsAdvertiser == nil {
			d.mdnsAdvertiser = NewAdvertise()
		}
		err = d.mdnsAdvertiser.advertiseOperational(advertiseParameters)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DnssdServerImpl) GetCommissionableInstanceName() string {
	name, _ := d.mdnsAdvertiser.GetCommissionableInstanceName()
	return name
}

func (d *DnssdServerImpl) SetEphemeralDiscriminator(discriminator uint16) error {
	if discriminator >= DeviceLayer.KMaxDiscriminatorValue {
		return lib.ChipErrorInvalidArgument
	}
	d.mEphemeralDiscriminator = &discriminator
	return nil
}

func (d *DnssdServerImpl) AdvertiseCommissioner() error {
	return d.Advertise(false, CommissioningMode_Disabled)
}

func (d *DnssdServerImpl) haveOperationalCredentials() bool {
	return d.mFabricTable.FabricCount() != 0
}

func (d *DnssdServerImpl) SetFabricTable(fabrics *credentials.FabricTable) {
	d.mFabricTable = fabrics
}

func (d *DnssdServerImpl) SetCommissioningModeProvider(manager CommissioningModeProvider) {
	d.mCommissioningModeProvider = manager
}

func (d *DnssdServerImpl) SetSecuredPort(port uint16) {
	d.mSecuredPort = port
}

func (d *DnssdServerImpl) SetUnsecuredPort(port uint16) {
	d.mUnsecuredPort = port
}

func (d *DnssdServerImpl) SetInterfaceId(id net.Interface) {
	d.mInterfaceId = id
}

func (d *DnssdServerImpl) StartServer() {
	mode := CommissioningMode_Disabled
	if d.mCommissioningModeProvider != nil {
		mode = d.mCommissioningModeProvider.GetCommissioningMode()
	}
	d.startServer(mode)
}

func (d *DnssdServerImpl) startServer(mode int) {

	log.Printf("Updating services using commissioning mode %d", mode)
	err := d.mdnsAdvertiser.Init()
	if err != nil {
		log.Error("failed to initialize advertiser: %s", err.Error())
	}
	err = d.mdnsAdvertiser.RemoveServices()
	if err != nil {
		log.Error("failed to remove advertised services: %s", err.Error())
	}
	err = d.AdvertiseOperational()
	if err != nil {
		log.Errorf("failed to advertise operational node: %s", err.Error())
	}

	if mode == CommissioningMode_Disabled {
		err := d.AdvertiseCommissionableNode(mode)
		if err != nil {
			log.Error("failed to advertise commissionable node: %s", err.Error())
			log.Infof(err.Error())
		}
	}

	if config.ChipDeviceConfigEnableCommissionerDiscovery != 0 {
		err := d.AdvertiseCommissioner()
		if err != nil {
			log.Errorf("failed to advertise commissioner: %s", err.Error())
		}
	}
	err = d.mdnsAdvertiser.FinalizeServiceUpdate()
	if err != nil {
		log.Errorf("Failed to finalize service update: %s", err.Error())
	}
}

func (d *DnssdServerImpl) AdvertiseCommissionableNode(mode int) error {
	if config.ChipDeviceConfigEnableExtendedDiscovery {
		d.mCurrentCommissioningMode = mode
	}
	if mode != CommissioningMode_Disabled {
		d.mExtendedDiscoveryExpiration = nil
	}
	return d.Advertise(true, mode)
}

func (d *DnssdServerImpl) Advertise(commissionAbleNode bool, mode int) error {

	advertiseParameters := params.NewCommissionAdvertisingParameters()
	if commissionAbleNode {
		advertiseParameters.SetPort(d.mSecuredPort)
		advertiseParameters.SetCommissionAdvertiseMode(AdvertiseMode_CommissionableNode)
	} else {
		advertiseParameters.SetPort(d.mUnsecuredPort)
		advertiseParameters.SetCommissionAdvertiseMode(AdvertiseMode_Commissioner)

	}
	advertiseParameters.SetInterfaceId(d.mInterfaceId)
	advertiseParameters.SetCommissioningMode(mode)
	advertiseParameters.EnableIpV4(true)

	//set  mac
	mac, err := config.ConfigurationMgr().GetPrimaryMACAddress()
	if err != nil || mac == "" {
		mac = fmt.Sprintf("%016X", rand.Uint64())
	}
	advertiseParameters.SetMaC(mac)

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
	discriminator, err = DeviceLayer.GetCommissionableDateProviderInstance().GetSetupDiscriminator()
	if err != nil {
		log.Infof(
			"Setup discriminator read error: (%s )! Critical error, will not be commissionable.",
			err.Error())
		return err
	}
	if d.mEphemeralDiscriminator != nil {
		discriminator = *d.mEphemeralDiscriminator
	}
	advertiseParameters.SetShortDiscriminator(uint8(discriminator>>8) & 0x0F).
		SetLongDiscriminator(discriminator)

	// set device type id
	deviceTypeId, err := config.ConfigurationMgr().GetDeviceTypeId()
	if config.ConfigurationMgr().IsCommissionableDeviceTypeEnabled() && err == nil {
		if err != nil {
			advertiseParameters.SetDeviceType(deviceTypeId)
		}
	}

	//set device name
	if config.ConfigurationMgr().IsCommissionableDeviceNameEnabled() {
		deviceName, err := config.ConfigurationMgr().GetCommissionableDeviceName()
		if err != nil {
			advertiseParameters.SetDeviceName(deviceName)
		}
	}

	advertiseParameters.SetLocalMRPConfig(messageing.GetLocalMRPConfig()).SetTcpSupported(config.InetConfigEnableTcpEndpoint)

	if !d.haveOperationalCredentials() {
		value, err := config.ConfigurationMgr().GetInitialPairingHint()
		if value != 0 && err == nil {
			advertiseParameters.SetPairingHint(value)
		} else {
			log.Infof("DNS-SD Pairing Hint not set")
		}
		str, err := config.ConfigurationMgr().GetInitialPairingInstruction()
		if err != nil {
			log.Infof("DNS-SD Pairing Instruction not set")
		} else {
			advertiseParameters.SetPairingInstruction(str)
		}
	} else {
		hint, err := config.ConfigurationMgr().GetSecondaryPairingHint()
		if err != nil {
			log.Infof("DNS-SD Pairing Hint not set")
		} else {
			advertiseParameters.SetPairingHint(hint)
		}

		str, err := config.ConfigurationMgr().GetSecondaryPairingInstruction()
		if err != nil {
			log.Infof("DNS-SD Pairing Instruction not set")
		} else {
			advertiseParameters.SetPairingInstruction(str)
		}
	}
	return d.mdnsAdvertiser.advertiseCommission(advertiseParameters)
}

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
