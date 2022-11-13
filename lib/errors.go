package lib

type MatterError string

var (
	InvalidArgument                   MatterError = "CHIP_ERROR_INVALID_ARGUMENT"
	IncorrectState                    MatterError = "CHIP_ERROR_INCORRECT_STATE"
	NotImplemented                    MatterError = "CHIP_ERROR_NOT_IMPLEMENTED"
	NotMemory                         MatterError = "Not memory"
	ErrorInternal                     MatterError = "CHIP_ERROR_INTERNAL"
	DeviceErrorConfigNotFound         MatterError = "CHIP_DEVICE_ERROR_CONFIG_NOT_FOUND"
	DuplicateMessageReceived          MatterError = "CHIP_ERROR_DUPLICATE_MESSAGE_RECEIVED"
	InvalidFabricIndex                MatterError = "CHIP_ERROR_INVALID_FABRIC_INDEX"
	TooManyUnsolicitedMessageHandlers MatterError = "CHIP_ERROR_TOO_MANY_UNSOLICITED_MESSAGE_HANDLERS"
	WrongEncryptionTypeFromPeer       MatterError = "CHIP_ERROR_WRONG_ENCRYPTION_TYPE_FROM_PEER"
	InvalidMessageType                MatterError = "CHIP_ERROR_INVALID_MESSAGE_TYPE"
	InvalidMessage                    MatterError = "INVALID MESSAGE"
	WrongTlvType                      MatterError = "CHIP_ERROR_WRONG_TLV_TYPE"
	UnexpectedTlvElement              MatterError = "CHIP_ERROR_UNEXPECTED_TLV_ELEMENT"
	InvalidCaseParameter              MatterError = "CHIP_ERROR_INVALID_CASE_PARAMETER"
	InvalidTlvTag                     MatterError = "CHIP_ERROR_INVALID_TLV_TAG"
	KeyNotFound                       MatterError = "CHIP_ERROR_KEY_NOT_FOUND"
	NoUnsolicitedMessageHandler       MatterError = "CHIP_ERROR_NO_UNSOLICITED_MESSAGE_HANDLER"
)

func (e MatterError) Error() string {
	value := string(e)
	return value
}
