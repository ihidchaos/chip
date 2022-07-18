package config

import (
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
	"os"
	"path"
	"sync"
)

func getDefaultKVS() string {
	dir, _ := os.UserHomeDir()
	return path.Join(dir, "/chip.ini")
}

type ConfigFlag struct {
	Key          string
	DefaultValue any
	Usage        string
}

var (
	DeviceOption_Version = ConfigFlag{
		Key:          "version",
		DefaultValue: 0,
		Usage:        "The version indication provides versioning of the setup payload.\n",
	}

	DeviceOption_VendorID = ConfigFlag{
		"vendor-id",
		0,
		"The Vendor ID is assigned by the Connectivity Standards Alliance.\n",
	}

	DeviceOption_ProductID = ConfigFlag{
		"product-id",
		0,
		"The Product ID is specified by vendor.\n",
	}

	DeviceOption_CustomFlow = ConfigFlag{
		"custom-flow",
		0,
		"A 2-bit unsigned enumeration specifying manufacturer-specific custom flow options.\n",
	}

	DeviceOption_Capabilities = ConfigFlag{
		"capabilities",
		0,
		"Discovery Capabilities Bitmask which contains information about Device’s available technologies for device discovery.\n",
	}

	DeviceOption_Discriminator = ConfigFlag{
		"discriminator",
		0,
		"A 12-bit unsigned integer match the value which a device advertises during commissioning.\n",
	}

	DeviceOption_Passcode = ConfigFlag{
		"passcode",
		0xFFFFFFF,
		"A 27-bit unsigned integer, which serves as proof of possession during commissioning. If not provided to compute a verifier, the --spake2p-verifier-base64 must be provided. \n",
	}

	DeviceOption_Spake2pVerifierBase64 = ConfigFlag{
		"spake2p-verifier-base64",
		0xFFFFF,
		"A raw concatenation of 'W0' and 'L' (67 bytes) as base64 to override the verifier auto-computed from the passcode, if provided.\n",
	}

	DeviceOption_Spake2pSaltBase64 = ConfigFlag{
		"spake2p-salt-base64",
		0,
		"16-32 bytes of salt to use for the PASE verifier, as base64. If omitted, will be generated randomly. If a --spake2p-verifier-base64 is passed, it must match against the salt otherwise failure will arise.\n",
	}

	DeviceOption_Spake2pIterations = ConfigFlag{
		"spake2p-iterations",
		0,
		"Number of PB DF iterations to use. If omitted, will be 1000. If a --spake2p-verifier-base64 is passed, the iteration counts must match that used to generate the verifier otherwise failure will arise.\n",
	}

	DeviceOption_SecuredDevicePort = ConfigFlag{
		"secured-device-port",
		5540,
		"A 16-bit unsigned integer specifying the listen port to use for secure device messages (default is 5540).\n",
	}

	DeviceOption_SecuredCommissionerPort = ConfigFlag{
		"secured-commissioner-port",
		5542,
		"A 16-bit unsigned integer specifying the listen port to use for secure commissioner messages (default is 5542). Only valid when app is both device and commissioner.\n",
	}

	DeviceOption_UnsecuredCommissionerPort = ConfigFlag{
		"unsecured-commissioner-port",
		5550,
		"A 16-bit unsigned integer specifying the port to use for unsecured commissioner messages (default is 5550).\n",
	}

	DeviceOption_Command = ConfigFlag{
		"command",
		"command",
		"A name for a command to execute during startup.\n"}

	DeviceOption_PICS = ConfigFlag{
		"PICS",
		"",
		"A file containing PICS items.\n"}

	DeviceOption_KVS = ConfigFlag{
		"KVS",
		getDefaultKVS(),
		"A file to store Key Value Store items.\n"}

	DeviceOption_InterfaceId = ConfigFlag{
		"interface-id",
		"interface-id",
		"A interface id to advertise on.\n"}
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

var _instance *DeviceOptions
var _once sync.Once

func GetDeviceOptionsInstance() *DeviceOptions {
	_once.Do(func() {
		if _instance == nil {
			_instance = &DeviceOptions{}
		}
	})
	return _instance
}

func FlagsDeviceOptions(c *cobra.Command) {
	c.Flags().Uint8(DeviceOption_Version.Key, cast.ToUint8(DeviceOption_Version.DefaultValue), DeviceOption_Version.Usage)
	c.Flags().Uint64(DeviceOption_VendorID.Key, cast.ToUint64(DeviceOption_VendorID.DefaultValue), DeviceOption_VendorID.Usage)
	c.Flags().Uint64(DeviceOption_ProductID.Key, cast.ToUint64(DeviceOption_ProductID.DefaultValue), DeviceOption_ProductID.Usage)
	c.Flags().Uint8(DeviceOption_CustomFlow.Key, cast.ToUint8(DeviceOption_CustomFlow.DefaultValue), DeviceOption_CustomFlow.Usage)
	c.Flags().Uint8(DeviceOption_Capabilities.Key, cast.ToUint8(DeviceOption_Capabilities.DefaultValue), DeviceOption_Capabilities.Usage)
	c.Flags().Uint16(DeviceOption_Discriminator.Key, cast.ToUint16(DeviceOption_Discriminator.DefaultValue), DeviceOption_Discriminator.Usage)
	c.Flags().Uint32(DeviceOption_Passcode.Key, cast.ToUint32(DeviceOption_Passcode.DefaultValue), DeviceOption_Passcode.Usage)
	c.Flags().Uint32(DeviceOption_Spake2pVerifierBase64.Key, cast.ToUint32(DeviceOption_Spake2pVerifierBase64.DefaultValue), DeviceOption_Spake2pVerifierBase64.Usage)
	c.Flags().Uint32(DeviceOption_Spake2pSaltBase64.Key, cast.ToUint32(DeviceOption_Spake2pSaltBase64.DefaultValue), DeviceOption_Spake2pSaltBase64.Usage)
	c.Flags().Uint64(DeviceOption_Spake2pIterations.Key, cast.ToUint64(DeviceOption_Spake2pIterations.DefaultValue), DeviceOption_Spake2pIterations.Usage)
	c.Flags().Uint16(DeviceOption_SecuredDevicePort.Key, cast.ToUint16(DeviceOption_SecuredDevicePort.DefaultValue), DeviceOption_SecuredDevicePort.Usage)
	c.Flags().Uint16(DeviceOption_SecuredCommissionerPort.Key, cast.ToUint16(DeviceOption_SecuredCommissionerPort.DefaultValue), DeviceOption_SecuredCommissionerPort.Usage)
	c.Flags().Uint16(DeviceOption_UnsecuredCommissionerPort.Key, cast.ToUint16(DeviceOption_UnsecuredCommissionerPort.DefaultValue), DeviceOption_UnsecuredCommissionerPort.Usage)
	c.Flags().String(DeviceOption_Command.Key, cast.ToString(DeviceOption_Command.DefaultValue), DeviceOption_Command.Usage)
	c.Flags().String(DeviceOption_PICS.Key, cast.ToString(DeviceOption_PICS.DefaultValue), DeviceOption_PICS.Usage)
	c.Flags().String(DeviceOption_KVS.Key, cast.ToString(DeviceOption_KVS.DefaultValue), DeviceOption_KVS.Usage)
	c.Flags().String(DeviceOption_InterfaceId.Key, cast.ToString(DeviceOption_InterfaceId.DefaultValue), DeviceOption_InterfaceId.Usage)
}

func InitDeviceOptions(config *viper.Viper) *DeviceOptions {

	GetDeviceOptionsInstance().Payload.Version = uint8(config.GetUint(DeviceOption_Version.Key))
	GetDeviceOptionsInstance().Payload.VendorID = uint16(config.GetUint32(DeviceOption_VendorID.Key))
	GetDeviceOptionsInstance().Payload.Discriminator = uint16(config.GetUint32(DeviceOption_Discriminator.Key))

	GetDeviceOptionsInstance().SecuredDevicePort = uint16(config.GetUint32(DeviceOption_SecuredDevicePort.Key))
	GetDeviceOptionsInstance().SecuredCommissionerPort = uint16(config.GetUint32(DeviceOption_SecuredCommissionerPort.Key))
	GetDeviceOptionsInstance().UnsecuredCommissionerPort = uint16(config.GetUint32(DeviceOption_UnsecuredCommissionerPort.Key))
	GetDeviceOptionsInstance().Command = config.GetString(DeviceOption_Command.Key)
	GetDeviceOptionsInstance().PICS = config.GetString(DeviceOption_PICS.Key)
	GetDeviceOptionsInstance().KVS = config.GetString(DeviceOption_KVS.Key)
	GetDeviceOptionsInstance().InterfaceId = net.Interface{}
	GetDeviceOptionsInstance().TraceStreamDecodeEnabled = false
	GetDeviceOptionsInstance().TraceStreamToLogEnabled = false

	return GetDeviceOptionsInstance()
}
