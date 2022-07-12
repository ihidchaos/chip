package config

import (
	"github.com/galenliu/chip/platform/options"
	"github.com/spf13/viper"
	"net"
)

const (
	KDeviceOption_Version                   = "version"
	VDeviceOption_Version                   = "0"
	KDeviceOption_VendorID                  = "vendor-id"
	KDeviceOption_ProductID                 = "product-id"
	KDeviceOption_CustomFlow                = "custom-flow"
	KDeviceOption_Capabilities              = "capabilities"
	KDeviceOption_Discriminator             = "discriminator"
	KDeviceOption_Passcode                  = "passcode"
	KDeviceOption_Spake2pVerifierBase64     = "spake2p-verifier-base64"
	KDeviceOption_Spake2pSaltBase64         = "spake2p-salt-base64"
	KDeviceOption_Spake2pIterations         = "spake2p-iterations"
	KDeviceOption_SecuredDevicePort         = "secured-device-port"
	KDeviceOption_SecuredCommissionerPort   = "secured-commissioner-port"
	KDeviceOption_UnsecuredCommissionerPort = "unsecured-commissioner-port"
	KDeviceOption_Command                   = "command"
	KDeviceOption_PICS                      = "PICS"
	KDeviceOption_KVS                       = "KVS"
	KDeviceOption_InterfaceId               = "interface-id"
)

func GetDeviceOptions(config *viper.Viper) *options.DeviceOptions {

	deviceOptions := &options.DeviceOptions{
		Spake2pIterations: 0,
		Payload: options.PayloadContents{
			Version:               uint8(config.GetUint(KDeviceOption_Version)),
			VendorID:              uint16(config.GetUint32(KDeviceOption_VendorID)),
			CommissioningFlow:     0,
			RendezvousInformation: 0,
			Discriminator:         uint16(config.GetUint32(KDeviceOption_Discriminator)),
			SetUpPINCode:          0,
			IsValidQRCodePayload:  false,
			IsValidManualCode:     false,
			IsShortDiscriminator:  false,
		},
		BleDevice:                 0,
		WiFi:                      false,
		Thread:                    false,
		SecuredDevicePort:         uint16(config.GetUint32(KDeviceOption_SecuredDevicePort)),
		SecuredCommissionerPort:   uint16(config.GetUint32(KDeviceOption_SecuredCommissionerPort)),
		UnsecuredCommissionerPort: uint16(config.GetUint32(KDeviceOption_UnsecuredCommissionerPort)),
		Command:                   config.GetString(KDeviceOption_Command),
		PICS:                      config.GetString(KDeviceOption_PICS),
		KVS:                       config.GetString(KDeviceOption_KVS),
		InterfaceId:               net.Interface{},
		TraceStreamDecodeEnabled:  false,
		TraceStreamToLogEnabled:   false,
		TraceStreamFilename:       "",
		TestEventTriggerEnableKey: nil,
	}
	return deviceOptions
}
