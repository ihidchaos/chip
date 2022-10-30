package lib

import "fmt"

var (
	MatterErrorInvalidArgument                   = fmt.Errorf("CHIP_ERROR_INVALID_ARGUMENT")
	MatterErrorIncorrectState                    = fmt.Errorf("CHIP_ERROR_INCORRECT_STATE")
	MatterErrorNotImplemented                    = fmt.Errorf("CHIP_ERROR_NOT_IMPLEMENTED")
	MatterErrorInternal                          = fmt.Errorf("CHIP_ERROR_INTERNAL")
	MatterDeviceErrorConfigNotFound              = fmt.Errorf("CHIP_DEVICE_ERROR_CONFIG_NOT_FOUND")
	MatterErrorDuplicateMessageReceived          = fmt.Errorf("CHIP_ERROR_DUPLICATE_MESSAGE_RECEIVED")
	MatterErrorInvalidFabricIndex                = fmt.Errorf("CHIP_ERROR_INVALID_FABRIC_INDEX")
	MatterErrorTooManyUnsolicitedMessageHandlers = fmt.Errorf("CHIP_ERROR_TOO_MANY_UNSOLICITED_MESSAGE_HANDLERS")
	MatterErrorWrongEncryptionTypeFromPeer       = fmt.Errorf("CHIP_ERROR_WRONG_ENCRYPTION_TYPE_FROM_PEER")
	MatterErrorInvalidMessageType                = fmt.Errorf("CHIP_ERROR_INVALID_MESSAGE_TYPE")
	MatterErrorWrongTlvType                      = fmt.Errorf("CHIP_ERROR_WRONG_TLV_TYPE")
	MatterErrorUnexpectedTlvElement              = fmt.Errorf("CHIP_ERROR_UNEXPECTED_TLV_ELEMENT")
	MatterErrorInvalidCaseParameter              = fmt.Errorf("CHIP_ERROR_INVALID_CASE_PARAMETER")
	MatterErrorInvalidTlvTag                     = fmt.Errorf("CHIP_ERROR_INVALID_TLV_TAG")

	MatterErrorKeyNotFound = fmt.Errorf("CHIP_ERROR_KEY_NOT_FOUND")

	MatterErrorNoUnsolicitedMessageHandler = fmt.Errorf("CHIP_ERROR_NO_UNSOLICITED_MESSAGE_HANDLER")
)
