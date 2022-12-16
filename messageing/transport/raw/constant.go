package raw

const (
	kFixedUnencryptedHeaderSizeBytes = 8
	kNodeIdSizeBytes                 = 8
	kGroupIdSizeBytes                = 2
	kEncryptedHeaderSizeBytes        = 6
	kVendorIdSizeBytes               = 2
	kAckMessageCounterSizeBytes      = 4

	MaxAppMessageLen                     = 1200
	MaxTagLen                            = 16
	kMsgUnsecuredUnicastSessionId uint16 = 0x0000
	kMsgHeaderVersion             uint8  = 0x0000
)
