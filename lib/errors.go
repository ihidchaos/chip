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

	ChipErrorNoUnsolicitedMessageHandler = fmt.Errorf("CHIP_ERROR_NO_UNSOLICITED_MESSAGE_HANDLER")
)
