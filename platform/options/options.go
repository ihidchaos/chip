package options

const (
	kDeviceOption_BleDevice                             = 0x1000
	kDeviceOption_WiFi                                  = 0x1001
	kDeviceOption_Thread                                = 0x1002
	kDeviceOption_Version                               = 0x1003
	kDeviceOption_VendorID                              = 0x1004
	kDeviceOption_ProductID                             = 0x1005
	kDeviceOption_CustomFlow                            = 0x1006
	kDeviceOption_Capabilities                          = 0x1007
	kDeviceOption_Discriminator                         = 0x1008
	kDeviceOption_Passcode                              = 0x1009
	kDeviceOption_SecuredDevicePort                     = 0x100a
	kDeviceOption_SecuredCommissionerPort               = 0x100b
	kDeviceOption_UnsecuredCommissionerPort             = 0x100c
	kDeviceOption_Command                               = 0x100d
	kDeviceOption_PICS                                  = 0x100e
	kDeviceOption_KVS                                   = 0x100f
	kDeviceOption_InterfaceId                           = 0x1010
	kDeviceOption_Spake2pVerifierBase64                 = 0x1011
	kDeviceOption_Spake2pSaltBase64                     = 0x1012
	kDeviceOption_Spake2pIterations                     = 0x1013
	kDeviceOption_TraceFile                             = 0x1014
	kDeviceOption_TraceLog                              = 0x1015
	kDeviceOption_TraceDecode                           = 0x1016
	kOptionCSRResponseCSRIncorrectType                  = 0x1017
	kOptionCSRResponseCSRNonceIncorrectType             = 0x1018
	kOptionCSRResponseCSRNonceTooLong                   = 0x1019
	kOptionCSRResponseCSRNonceInvalid                   = 0x101a
	kOptionCSRResponseNOCSRElementsTooLong              = 0x101b
	kOptionCSRResponseAttestationSignatureIncorrectType = 0x101c
	kOptionCSRResponseAttestationSignatureInvalid       = 0x101d
	kOptionCSRResponseCSRExistingKeyPair                = 0x101e
	kDeviceOption_TestEventTriggerEnableKey             = 0x101f
)

//type DeviceOptions struct {
//	Spake2pIterations         uint32
//	Spake2pVerifier           []byte
//	Spake2pSalt               []byte
//	Discriminator             uint16
//	Payload                   config.PayloadContents
//	BleDevice                 uint32
//	WiFi                      bool
//	Thread                    bool
//	SecuredDevicePort         uint16
//	SecuredCommissionerPort   uint16
//	UnsecuredCommissionerPort uint16
//	Command                   string
//	PICS                      string
//	KVS                       string
//	InterfaceId               net.Interface
//	TraceStreamDecodeEnabled  bool
//	TraceStreamToLogEnabled   bool
//	TraceStreamFilename       string
//	TestEventTriggerEnableKey []byte
//}

//var once sync.Once
//var _instance *DeviceOptions
//
//func GetInstance() *DeviceOptions {
//	once.Do(func() {
//		_instance = &DeviceOptions{
//			Spake2pIterations:         0,
//			Payload:                   config.PayloadContents{},
//			BleDevice:                 0,
//			WiFi:                      false,
//			Thread:                    false,
//			SecuredDevicePort:         0,
//			SecuredCommissionerPort:   0,
//			UnsecuredCommissionerPort: 0,
//			Command:                   "",
//			PICS:                      "",
//			KVS:                       "",
//			InterfaceId:               net.Interface{},
//			TraceStreamDecodeEnabled:  false,
//			TraceStreamToLogEnabled:   false,
//			TraceStreamFilename:       "",
//			TestEventTriggerEnableKey: nil,
//		}
//	})
//	return _instance
//}
//
//func (option *DeviceOptions) HandleOption(optionSet map[string]string, aIdentifier int, aValue string) {
//	switch aIdentifier {
//	case kDeviceOption_BleDevice:
//		break
//	case kDeviceOption_WiFi:
//		//option.mWiFi = true
//		break
//	case kDeviceOption_Thread:
//		//option.mThread = true
//		break
//	case kDeviceOption_Version:
//		//v, _ := strconv.Atoi(aValue)
//		//option.Payload.version = uint8(v)
//		break
//	}
//
//}
