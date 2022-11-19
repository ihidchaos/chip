package dnssd

import (
	"fmt"
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/credentials"
	DeviceLayer "github.com/galenliu/chip/device"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/session"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"sync"
)

type Base interface {
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

func Instance() *Dnssd {
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
				mdnsAdvertiser:               NewAdvertise(),
			}
		}
	})
	return _DnssdInstance
}
func New() Base {
	return Instance()
}

func (d *Dnssd) SetFabricTable(fabrics *credentials.FabricTable) {
	d.mFabricTable = fabrics
}

func (d *Dnssd) SetCommissioningModeProvider(provider CommissioningModeProvider) {
	d.mCommissioningModeProvider = provider
}

func (d *Dnssd) SetSecuredPort(port uint16) {
	d.mSecuredPort = port
}

func (d *Dnssd) SetUnsecuredPort(port uint16) {
	d.mUnsecuredPort = port
}

func (d *Dnssd) SetInterfaceId(n net.Interface) {
	d.mInterfaceId = n
}

func (d *Dnssd) StartServer() {
	mode := CommissioningModeDisabled
	if d.mCommissioningModeProvider != nil {
		mode = d.mCommissioningModeProvider.GetCommissioningMode()
	}
	d.start(mode)
}

func (d *Dnssd) AdvertiseOperational() error {
	if d.mFabricTable == nil {
		return lib.IncorrectState
	}
	for _, info := range d.mFabricTable.Fabrics() {
		mac, err := config.DefaultManager().GetPrimaryMACAddress()
		if mac == "" || err != nil {
			mac = fmt.Sprintf("%016X", rand.Uint64())
		}
		advertiseParameters := NewOperationalAdvertisingParameters()
		advertiseParameters.SetPeerId(PeerId{info.GetNodeId(), info.CompressedFabricId()})
		advertiseParameters.SetMaC(mac)
		advertiseParameters.SetPort(d.mSecuredPort)
		advertiseParameters.SetInterfaceId(d.mInterfaceId)
		advertiseParameters.SetLocalMRPConfig(session.GetLocalMRPConfig())
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

func (d *Dnssd) Advertise(commissionableNode bool, mode int) error {

	advertiseParameters := NewCommissionAdvertisingParameters()
	if commissionableNode {
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
	mac, err := config.DefaultManager().GetPrimaryMACAddress()
	if err != nil || mac == "" {
		mac = fmt.Sprintf("%016X", rand.Uint64())
	}
	advertiseParameters.SetMaC(mac)

	//Sets device vendor id
	vid, err := DeviceLayer.DefaultInstanceInfo().GetVendorId()
	if err != nil {
		log.Infof("Vendor ID not known")
	} else {
		advertiseParameters.SetVendorId(vid)
	}

	// set  productId
	pid, err := DeviceLayer.DefaultInstanceInfo().GetProductId()
	if err != nil {
		log.Infof("Product ID not known")
	} else {
		advertiseParameters.SetProductId(pid)
	}

	// set discriminator
	var discriminator uint16 = 0
	discriminator, err = DeviceLayer.DefaultCommissionableDateProvider().GetSetupDiscriminator()
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
	deviceTypeId, err := config.DefaultManager().GetDeviceTypeId()
	if config.DefaultManager().IsCommissionableDeviceTypeEnabled() && err == nil {
		if err != nil {
			advertiseParameters.SetDeviceType(deviceTypeId)
		}
	}

	//set device name
	if config.DefaultManager().IsCommissionableDeviceNameEnabled() {
		deviceName, err := config.DefaultManager().GetCommissionableDeviceName()
		if err != nil {
			advertiseParameters.SetDeviceName(deviceName)
		}
	}

	advertiseParameters.SetLocalMRPConfig(session.GetLocalMRPConfig()).SetTcpSupported(config.InetConfigEnableTcpEndpoint)

	if !d.haveOperationalCredentials() {
		value, err := config.DefaultManager().GetInitialPairingHint()
		if value != 0 && err == nil {
			advertiseParameters.SetPairingHint(value)
		} else {
			log.Infof("DNS-SD Pairing Hint not set")
		}
		str, err := config.DefaultManager().GetInitialPairingInstruction()
		if err != nil {
			log.Infof("DNS-SD Pairing Instruction not set")
		} else {
			advertiseParameters.SetPairingInstruction(str)
		}
	} else {
		hint, err := config.DefaultManager().GetSecondaryPairingHint()
		if err != nil {
			log.Infof("DNS-SD Pairing Hint not set")
		} else {
			advertiseParameters.SetPairingHint(hint)
		}

		str, err := config.DefaultManager().GetSecondaryPairingInstruction()
		if err != nil {
			log.Infof("DNS-SD Pairing Instruction not set")
		} else {
			advertiseParameters.SetPairingInstruction(str)
		}
	}
	return d.mdnsAdvertiser.AdvertiseCommission(advertiseParameters)
}

func (d *Dnssd) start(mode int) {

	log.Printf("updating services using commissioning mode %d", mode)
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

	//如果Commission模式没有禁用，则广播Node
	if mode != CommissioningModeDisabled {
		err := d.AdvertiseCommissionableNode(mode)
		if err != nil {
			log.Error("failed to advertise commissionable node: %s", err.Error())
			log.Infof(err.Error())
		}
	}

	if config.CommissionerDiscovery != 0 {
		err := d.AdvertiseCommissioner()
		if err != nil {
			log.Errorf("failed to advertise commissioner: %s", err.Error())
		}
	}
	err = d.mdnsAdvertiser.FinalizeServiceUpdate()
	if err != nil {
		log.Errorf("failed to finalize service update: %s", err.Error())
	}
}

func (d *Dnssd) AdvertiseCommissionableNode(mode int) error {
	if config.ExtendedDiscovery {
		d.mCurrentCommissioningMode = mode
	}
	if mode != CommissioningModeDisabled {
		d.mExtendedDiscoveryExpiration = nil
	}
	return d.Advertise(true, mode)
}

func (d *Dnssd) haveOperationalCredentials() bool {
	return d.mFabricTable != nil && d.mFabricTable.FabricCount() != 0
}
