package lib

import "fmt"

var (
	ChipErrorInvalidArgument                   = fmt.Errorf("CHIP_ERROR_INVALID_ARGUMENT")
	ChipErrorIncorrectState                    = fmt.Errorf("CHIP_ERROR_INCORRECT_STATE")
	ChipErrorNotImplemented                    = fmt.Errorf("CHIP_ERROR_NOT_IMPLEMENTED")
	ChipErrorInternal                          = fmt.Errorf("CHIP_ERROR_INTERNAL")
	ChipDeviceErrorConfigNotFound              = fmt.Errorf("CHIP_DEVICE_ERROR_CONFIG_NOT_FOUND")
	ChipErrorDuplicateMessageReceived          = fmt.Errorf("CHIP_ERROR_DUPLICATE_MESSAGE_RECEIVED")
	ChipErrorInvalidFabricIndex                = fmt.Errorf("CHIP_ERROR_INVALID_FABRIC_INDEX")
	ChipErrorTooManyUnsolicitedMessageHandlers = fmt.Errorf("CHIP_ERROR_TOO_MANY_UNSOLICITED_MESSAGE_HANDLERS")
	ChipErrorWrongEncryptionTypeFromPeer       = fmt.Errorf("CHIP_ERROR_WRONG_ENCRYPTION_TYPE_FROM_PEER")
	ChipErrorInvalidMessageType                = fmt.Errorf("CHIP_ERROR_INVALID_MESSAGE_TYPE")
	ChipErrorWrongTlvType                      = fmt.Errorf("CHIP_ERROR_WRONG_TLV_TYPE")
	ChipErrorUnexpectedTlvElement              = fmt.Errorf("CHIP_ERROR_UNEXPECTED_TLV_ELEMENT")
	ChipErrorInvalidCaseParameter              = fmt.Errorf("CHIP_ERROR_INVALID_CASE_PARAMETER")

	ChipErrorKeyNotFound = fmt.Errorf("CHIP_ERROR_KEY_NOT_FOUND")

	ChipErrorNoUnsolicitedMessageHandler = fmt.Errorf("CHIP_ERROR_NO_UNSOLICITED_MESSAGE_HANDLER")
)
