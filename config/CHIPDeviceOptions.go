package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
)

const (
	KDeviceOption_Version       = "version"
	KDeviceOption_Version_DEF   = 0
	KDeviceOption_Version_USAGE = "The version indication provides versioning of the setup payload.\n"

	KDeviceOption_VendorID              = "vendor-id"
	KDeviceOption_VendorID_DEF   uint64 = 0
	KDeviceOption_VendorID_USAGE        = "The Vendor ID is assigned by the Connectivity Standards Alliance.\n"

	KDeviceOption_ProductID              = "product-id"
	KDeviceOption_ProductID_DEF   uint64 = 0
	KDeviceOption_ProductID_USAGE        = "       The Product ID is specified by vendor.\n"

	KDeviceOption_CustomFlow             = "custom-flow"
	KDeviceOption_CustomFlow_DEF   uint8 = 0
	KDeviceOption_CustomFlow_USAGE       = "A 2-bit unsigned enumeration specifying manufacturer-specific custom flow options.\n"

	KDeviceOption_Capabilities             = "capabilities"
	KDeviceOption_Capabilities_DEF   uint8 = 0
	KDeviceOption_Capabilities_USAGE       = "Discovery Capabilities Bitmask which contains information about Device’s available technologies for device discovery.\\n\""

	KDeviceOption_Discriminator              = "discriminator"
	KDeviceOption_Discriminator_DEF   uint16 = 0
	KDeviceOption_Discriminator_USAGE        = "A 12-bit unsigned integer match the value which a device advertises during commissioning.\n"

	KDeviceOption_Passcode              = "passcode"
	KDeviceOption_Passcode_DEF   uint32 = 0xFFFFFFF
	KDeviceOption_Passcode_USAGE        = "A 27-bit unsigned integer, which serves as proof of possession during commissioning. If not provided to compute a verifier, the --spake2p-verifier-base64 must be provided. \n"

	KDeviceOption_Spake2pVerifierBase64              = "spake2p-verifier-base64"
	KDeviceOption_Spake2pVerifierBase64_DEF   uint32 = 0xFFFFF
	KDeviceOption_Spake2pVerifierBase64_USAGE        = "A raw concatenation of 'W0' and 'L' (67 bytes) as base64 to override the verifier auto-computed from the passcode, if provided."

	KDeviceOption_Spake2pSaltBase64              = "spake2p-salt-base64"
	KDeviceOption_Spake2pSaltBase64_DEF   uint32 = 0
	KDeviceOption_Spake2pSaltBase64_USAGE        = "16-32 bytes of salt to use for the PASE verifier, as base64. If omitted, will be generated randomly. If a --spake2p-verifier-base64 is passed, it must match against the salt otherwise failure will arise."

	KDeviceOption_Spake2pIterations              = "spake2p-iterations"
	KDeviceOption_Spake2pIterations_DEF   uint64 = 0
	KDeviceOption_Spake2pIterations_USAGE        = "Number of PBKDF iterations to use. If omitted, will be 1000. If a --spake2p-verifier-base64 is passed, the iteration counts must match that used to generate the verifier otherwise failure will arise."

	KDeviceOption_SecuredDevicePort              = "secured-device-port"
	KDeviceOption_SecuredDevicePort_DEF   uint16 = 5540
	KDeviceOption_SecuredDevicePort_USAGE        = "A 16-bit unsigned integer specifying the listen port to use for secure device messages (default is 5540)."

	KDeviceOption_SecuredCommissionerPort              = "secured-commissioner-port"
	KDeviceOption_SecuredCommissionerPort_DEF   uint16 = 5542
	KDeviceOption_SecuredCommissionerPort_USAGE        = "A 16-bit unsigned integer specifying the listen port to use for secure commissioner messages (default is 5542). Only valid when app is both device and commissioner"

	KDeviceOption_UnsecuredCommissionerPort       = "unsecured-commissioner-port"
	KDeviceOption_UnsecuredCommissionerPort_DEF   = 5550
	KDeviceOption_UnsecuredCommissionerPort_USAGE = "A 16-bit unsigned integer specifying the port to use for unsecured commissioner messages (default is 5550)."

	KDeviceOption_Command       = "command"
	KDeviceOption_Command_DEF   = "command"
	KDeviceOption_Command_USAGE = "A name for a command to execute during startup."

	KDeviceOption_PICS       = "PICS"
	KDeviceOption_PICS_DEF   = ""
	KDeviceOption_PICS_USAGE = "A file containing PICS items."

	KDeviceOption_KVS       = "KVS"
	KDeviceOption_KVS_DEF   = ""
	KDeviceOption_KVS_USAGE = "A file to store Key Value Store items."

	KDeviceOption_InterfaceId       = "interface-id"
	KDeviceOption_InterfaceId_DEF   = "interface-id"
	KDeviceOption_InterfaceId_USAGE = "A interface id to advertise on."
)

type DeviceOptions struct {
	Spake2pIterations         uint32
	Spake2pVerifier           []byte
	Spake2pSalt               []byte
	Discriminator             uint16
	Payload                   PayloadContents
	BleDevice                 uint32
	WiFi                      bool
	Thread                    bool
	SecuredDevicePort         uint16
	SecuredCommissionerPort   uint16
	UnsecuredCommissionerPort uint16
	Command                   string
	PICS                      string
	KVS                       string
	InterfaceId               net.Interface
	TraceStreamDecodeEnabled  bool
	TraceStreamToLogEnabled   bool
	TraceStreamFilename       string
	TestEventTriggerEnableKey []byte
}

func InitDeviceOptions(c *cobra.Command) {
	c.Flags().Uint8(KDeviceOption_Version, KDeviceOption_Version_DEF, KDeviceOption_Version_USAGE)
	c.Flags().Uint64(KDeviceOption_VendorID, KDeviceOption_VendorID_DEF, KDeviceOption_VendorID_USAGE)
	c.Flags().Uint64(KDeviceOption_ProductID, KDeviceOption_ProductID_DEF, KDeviceOption_ProductID_USAGE)
	c.Flags().Uint8(KDeviceOption_CustomFlow, KDeviceOption_CustomFlow_DEF, KDeviceOption_CustomFlow_USAGE)
	c.Flags().Uint8(KDeviceOption_Capabilities, KDeviceOption_Capabilities_DEF, KDeviceOption_Capabilities_USAGE)
	c.Flags().Uint16(KDeviceOption_Discriminator, KDeviceOption_Discriminator_DEF, KDeviceOption_Discriminator_USAGE)
	c.Flags().Uint32(KDeviceOption_Passcode, KDeviceOption_Passcode_DEF, KDeviceOption_Passcode_USAGE)
	c.Flags().Uint32(KDeviceOption_Spake2pVerifierBase64, KDeviceOption_Spake2pVerifierBase64_DEF, KDeviceOption_Spake2pVerifierBase64_USAGE)
	c.Flags().Uint32(KDeviceOption_Spake2pSaltBase64, KDeviceOption_Spake2pSaltBase64_DEF, KDeviceOption_Spake2pSaltBase64_USAGE)
	c.Flags().Uint64(KDeviceOption_Spake2pIterations, KDeviceOption_Spake2pIterations_DEF, KDeviceOption_Spake2pIterations_USAGE)
	c.Flags().Uint16(KDeviceOption_SecuredDevicePort, KDeviceOption_SecuredDevicePort_DEF, KDeviceOption_SecuredDevicePort_USAGE)
	c.Flags().Uint16(KDeviceOption_SecuredCommissionerPort, KDeviceOption_SecuredCommissionerPort_DEF, KDeviceOption_SecuredCommissionerPort_USAGE)
	c.Flags().Uint16(KDeviceOption_UnsecuredCommissionerPort, KDeviceOption_UnsecuredCommissionerPort_DEF, KDeviceOption_UnsecuredCommissionerPort_USAGE)
	c.Flags().String(KDeviceOption_Command, KDeviceOption_Command_DEF, KDeviceOption_Command_USAGE)
	c.Flags().String(KDeviceOption_PICS, KDeviceOption_PICS_DEF, KDeviceOption_PICS_USAGE)
	c.Flags().String(KDeviceOption_KVS, KDeviceOption_KVS_DEF, KDeviceOption_KVS_USAGE)
	c.Flags().String(KDeviceOption_InterfaceId, KDeviceOption_InterfaceId_DEF, KDeviceOption_InterfaceId_USAGE)
}

func GetDeviceOptions(config *viper.Viper) *DeviceOptions {
	deviceOptions := &DeviceOptions{
		Spake2pIterations: 0,
		Payload: PayloadContents{
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
