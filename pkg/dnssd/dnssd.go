package dnssd

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials"
	DeviceLayer "github.com/galenliu/chip/device"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"sync"
)

type DnssdServer interface {
	SetFabricTable(fabrics *credentials.FabricTable)
	SetCommissioningModeProvider(provider CommissioningModeProvider)
	SetSecuredPort(port uint16)
	SetUnsecuredPort(port uint16)
	SetInterfaceId(net.Interface)
	StartServer()
}

type Dnssd struct {
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

var _DnssdInstance *Dnssd
var _DnssdInstanceOnce sync.Once

func GetInstance() *Dnssd {
	_DnssdInstanceOnce.Do(func() {
		if _DnssdInstance == nil {
			_DnssdInstance = &Dnssd{
				mSecuredPort:                 0,
				mUnsecuredPort:               0,
				mInterfaceId:                 net.Interface{},
				mFabricTable:                 nil,
				mCommissioningModeProvider:   nil,
				mCurrentCommissioningMode:    0,
				mExtendedDiscoveryExpiration: nil,
				mEphemeralDiscriminator:      nil,
				mdnsAdvertiser:               nil,
			}
		}
	})
	return _DnssdInstance
}

func NewDnssdInstance() *Dnssd {
	return GetInstance()
}

func (d Dnssd) SetFabricTable(fabrics *credentials.FabricTable) {
	d.mFabricTable = fabrics
}

func (d Dnssd) SetCommissioningModeProvider(provider CommissioningModeProvider) {
	d.mCommissioningModeProvider = provider
}

func (d Dnssd) SetSecuredPort(port uint16) {
	d.mSecuredPort = port
}

func (d Dnssd) SetUnsecuredPort(port uint16) {
	d.mUnsecuredPort = port
}

func (d Dnssd) SetInterfaceId(n net.Interface) {
	d.mInterfaceId = n
}

func (d *Dnssd) StartServer() {
	mode := CommissioningModeDisabled
	if d.mCommissioningModeProvider != nil {
		mode = d.mCommissioningModeProvider.GetCommissioningMode()
	}
	d.startServer(mode)
}

func (d *Dnssd) AdvertiseOperational() error {
	if d.mFabricTable == nil {
		return lib.ChipErrorIncorrectState
	}
	for _, info := range d.mFabricTable.GetFabrics() {
		mac, err := config.ConfigurationMgr().GetPrimaryMACAddress()
		if mac == "" || err != nil {
			mac = fmt.Sprintf("%016X", rand.Uint64())
		}
		advertiseParameters := NewOperationalAdvertisingParameters()
		advertiseParameters.SetPeerId(PeerId{info.GetNodeId(), info.GetCompressedFabricId()})
		advertiseParameters.SetMaC(mac)
		advertiseParameters.SetPort(d.mSecuredPort)
		advertiseParameters.SetInterfaceId(d.mInterfaceId)
		advertiseParameters.SetLocalMRPConfig(transport.GetLocalMRPConfig())
		advertiseParameters.SetTcpSupported(config.InetConfigEnableTcpEndpoint)
		advertiseParameters.EnableIpV4(true)
		if d.mdnsAdvertiser == nil {
			d.mdnsAdvertiser = NewAdvertise()
		}
		err = d.mdnsAdvertiser.AdvertiseOperational(advertiseParameters)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dnssd) AdvertiseCommissioner() error {
	return d.Advertise(false, CommissioningModeDisabled)
}

func (d *Dnssd) Advertise(commissionAbleNode bool, mode int) error {

	advertiseParameters := NewCommissionAdvertisingParameters()
	if commissionAbleNode {
		advertiseParameters.SetPort(d.mSecuredPort)
		advertiseParameters.SetCommissionAdvertiseMode(AdvertiseModeCommissionableNode)
	} else {
		advertiseParameters.SetPort(d.mUnsecuredPort)
		advertiseParameters.SetCommissionAdvertiseMode(AdvertiseModeCommissioner)

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
	vid, err := DeviceLayer.GetDeviceInstanceInfoProvider().GetVendorId()
	if err != nil {
		log.Infof("Vendor ID not known")
	} else {
		advertiseParameters.SetVendorId(vid)
	}

	// set  productId
	pid, err := DeviceLayer.GetDeviceInstanceInfoProvider().GetProductId()
	if err != nil {
		log.Infof("Product ID not known")
	} else {
		advertiseParameters.SetProductId(pid)
	}

	// set discriminator
	var discriminator uint16 = 0
	discriminator, err = DeviceLayer.GetCommissionableDateProvider().GetSetupDiscriminator()
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

	advertiseParameters.SetLocalMRPConfig(transport.GetLocalMRPConfig()).SetTcpSupported(config.InetConfigEnableTcpEndpoint)

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
	return d.mdnsAdvertiser.AdvertiseCommission(advertiseParameters)
}

func (d *Dnssd) startServer(mode int) {

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

	if mode == CommissioningModeDisabled {
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

func (d *Dnssd) AdvertiseCommissionableNode(mode int) error {
	if config.ChipDeviceConfigEnableExtendedDiscovery {
		d.mCurrentCommissioningMode = mode
	}
	if mode != CommissioningModeDisabled {
		d.mExtendedDiscoveryExpiration = nil
	}
	return d.Advertise(true, mode)
}

func (d *Dnssd) haveOperationalCredentials() bool {
	return d.mFabricTable.FabricCount() != 0
}
