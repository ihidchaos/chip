package pkg

import "fmt"

var (
	ChipErrorInvalidArgument          = fmt.Errorf("CHIP_ERROR_INVALID_ARGUMENT")
	ChipErrorIncorrectState           = fmt.Errorf("CHIP_ERROR_INCORRECT_STATE")
	ChipErrorNotImplemented           = fmt.Errorf("CHIP_ERROR_NOT_IMPLEMENTED")
	ChipErrorInternal                 = fmt.Errorf("CHIP_ERROR_INTERNAL")
	ChipDeviceErrorConfigNotFound     = fmt.Errorf("CHIP_DEVICE_ERROR_CONFIG_NOT_FOUND")
	ChipErrorDuplicateMessageReceived = fmt.Errorf("CHIP_ERROR_DUPLICATE_MESSAGE_RECEIVED")

	ChipErrorNoUnsolicitedMessageHandler = fmt.Errorf("CHIP_ERROR_NO_UNSOLICITED_MESSAGE_HANDLER")
)
