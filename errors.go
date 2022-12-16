package chip

type ErrorType string

var (
	MATTER_ERROR_SENDING_BLOCKED ErrorType = "Sending blocked"

	ErrorConnectionAborted ErrorType = "Connection aborted"

	MATTER_ERROR_MESSAGE_TOO_LONG ErrorType = "Message too long"

	MATTER_ERROR_UNSUPPORTED_EXCHANGE_VERSION ErrorType = "Unsupported exchange version"

	MATTER_ERROR_TOO_MANY_UNSOLICITED_MESSAGE_HANDLERS ErrorType = "Too many unsolicited message handlers"

	MATTER_ERROR_NO_UNSOLICITED_MESSAGE_HANDLER ErrorType = "No unsolicited message handler"

	MATTER_ERROR_NO_CONNECTION_HANDLER ErrorType = "No connection handler"

	MATTER_ERROR_TOO_MANY_PEER_NODES ErrorType = "Too many peer nodes"

	MATTER_ERROR_SENTINEL ErrorType = "ErrorInternal sentinel"

	ErrorNoMemory ErrorType = "No memory"

	MATTER_ERROR_NO_MESSAGE_HANDLER ErrorType = "No message handler"

	MATTER_ERROR_MESSAGE_INCOMPLETE ErrorType = "Message incomplete"

	MATTER_ERROR_DATA_NOT_ALIGNED ErrorType = "Data not aligned"

	MATTER_ERROR_UNKNOWN_KEY_TYPE ErrorType = "Unknown key type"

	MATTER_ERROR_KEY_NOT_FOUND ErrorType = "Key not found"

	MATTER_ERROR_WRONG_ENCRYPTION_TYPE ErrorType = "Wrong encryption type"

	MATTER_ERROR_TOO_MANY_KEYS ErrorType = "Too many keys"

	MATTER_ERROR_INTEGRITY_CHECK_FAILED ErrorType = "Integrity check failed"

	MATTER_ERROR_INVALID_SIGNATURE ErrorType = "Invalid signature"

	MATTER_ERROR_UNSUPPORTED_MESSAGE_VERSION ErrorType = "Unsupported message version"

	MATTER_ERROR_UNSUPPORTED_ENCRYPTION_TYPE ErrorType = "Unsupported encryption type"

	MATTER_ERROR_UNSUPPORTED_SIGNATURE_TYPE ErrorType = "Unsupported signature type"

	MATTER_ERROR_INVALID_MESSAGE_LENGTH ErrorType = "Invalid message length"

	MATTER_ERROR_BUFFER_TOO_SMALL ErrorType = "Buffer too small"

	MATTER_ERROR_DUPLICATE_KEY_ID ErrorType = "Duplicate key id"

	MATTER_ERROR_WRONG_KEY_TYPE ErrorType = "Wrong key type"

	MATTER_ERROR_WELL_UNINITIALIZED ErrorType = "Well uninitialized"

	MATTER_ERROR_WELL_EMPTY ErrorType = "Well empty"

	MATTER_ERROR_INVALID_STRING_LENGTH ErrorType = "Invalid string length"

	MATTER_ERROR_INVALID_LIST_LENGTH ErrorType = "invalid list length"

	MATTER_ERROR_INVALID_INTEGRITY_TYPE ErrorType = "Invalid integrity type"

	MATTER_END_OF_TLV ErrorType = "End of TLV"

	MATTER_ERROR_TLV_UNDERRUN ErrorType = "TLV underrun"

	MATTER_ERROR_INVALID_TLV_ELEMENT ErrorType = "Invalid TLV element"

	MATTER_ERROR_INVALID_TLV_TAG ErrorType = "Invalid TLV tag"

	MATTER_ERROR_UNKNOWN_IMPLICIT_TLV_TAG ErrorType = "Unknown implicit TLV tag"

	MATTER_ERROR_WRONG_TLV_TYPE ErrorType = "Wrong TLV type"

	MATTER_ERROR_TLV_CONTAINER_OPEN ErrorType = "TLV container open"

	MATTER_ERROR_INVALID_TRANSFER_MODE ErrorType = "Invalid transfer mode"

	MATTER_ERROR_INVALID_PROFILE_ID ErrorType = "Invalid profile id"

	MATTER_ERROR_UNEXPECTED_TLV_ELEMENT ErrorType = "Unexpected TLV element"

	MATTER_ERROR_STATUS_REPORT_RECEIVED ErrorType = "Status Report received from peer"

	MATTER_ERROR_NOT_IMPLEMENTED ErrorType = "Not Implemented"

	MATTER_ERROR_INVALID_ADDRESS ErrorType = "Invalid address"

	MATTER_ERROR_INVALID_ARGUMENT ErrorType = "Invalid argument"

	MATTER_ERROR_TLV_TAG_NOT_FOUND ErrorType = "TLV tag not found"

	MATTER_ERROR_MISSING_SECURE_SESSION ErrorType = "Missing secure session"

	MATTER_ERROR_INVALID_ADMIN_SUBJECT ErrorType = "CaseAdminSubject is not valid"

	MATTER_ERROR_INSUFFICIENT_PRIVILEGE ErrorType = "Required privilege was insufficient during an operation"

	MATTER_ERROR_IM_MALFORMED_ATTRIBUTE_REPORT_IB ErrorType = "Malformed Interacton Model Attribute Report IB"

	MATTER_ERROR_IM_MALFORMED_COMMAND_DATA_IB ErrorType = "Malformed Interacton Model Command Data IB"

	MATTER_ERROR_IM_MALFORMED_EVENT_STATUS_IB ErrorType = "Malformed Interacton Model Event Status IB"

	MATTER_ERROR_IM_MALFORMED_STATUS_RESPONSE_MESSAGE ErrorType = "Malformed Interacton Model Status Response IB"

	MATTER_ERROR_INVALID_PATH_LIST ErrorType = "Invalid TLV path list"

	MATTER_ERROR_INVALID_DATA_LIST ErrorType = "Invalid TLV data list"

	MATTER_ERROR_TRANSACTION_CANCELED ErrorType = "Transaction canceled"

	MATTER_ERROR_LISTENER_ALREADY_STARTED ErrorType = "Listener already started"

	MATTER_ERROR_LISTENER_ALREADY_STOPPED ErrorType = "Listener already stopped"

	MATTER_ERROR_INVALID_SUBSCRIPTION ErrorType = "Invalid Subscription Id"

	MATTER_ERROR_TIMEOUT ErrorType = "Timeout"

	MATTER_ERROR_INVALID_DEVICE_DESCRIPTOR ErrorType = "Invalid device descriptor"

	MATTER_ERROR_UNSUPPORTED_DEVICE_DESCRIPTOR_VERSION ErrorType = "Unsupported device descriptor version"

	CHIP_END_OF_INPUT ErrorType = "End of input"

	MATTER_ERROR_RATE_LIMIT_EXCEEDED ErrorType = "Rate limit exceeded"

	MATTER_ERROR_SECURITY_MANAGER_BUSY ErrorType = "Security manager busy"

	MATTER_ERROR_INVALID_PASE_PARAMETER ErrorType = "Invalid PASE parameter"

	MATTER_ERROR_PASE_SUPPORTS_ONLY_CONFIG1 ErrorType = "PASE supports only Config1"

	MATTER_ERROR_NO_COMMON_PASE_CONFIGURATIONS ErrorType = "No supported PASE configurations in common"

	MATTER_ERROR_INVALID_PASE_CONFIGURATION ErrorType = "Invalid PASE configuration"

	MATTER_ERROR_KEY_CONFIRMATION_FAILED ErrorType = "Key confirmation failed"

	ErrorInvalidUseOfSessionKey ErrorType = "Invalid use of session key"

	MATTER_ERROR_CONNECTION_CLOSED_UNEXPECTEDLY ErrorType = "Connection closed unexpectedly"

	MATTER_ERROR_MISSING_TLV_ELEMENT ErrorType = "Missing TLV element"

	MATTER_ERROR_RANDOM_DATA_UNAVAILABLE ErrorType = "Random data unavailable"

	MATTER_ERROR_UNSUPPORTED_HOST_PORT_ELEMENT ErrorType = "Unsupported type in host/port list"

	MATTER_ERROR_INVALID_HOST_SUFFIX_INDEX ErrorType = "Invalid suffix index in host/port list"

	MATTER_ERROR_HOST_PORT_LIST_EMPTY ErrorType = "Host/port empty"

	MATTER_ERROR_UNSUPPORTED_AUTH_MODE ErrorType = "Unsupported authentication mode"

	MATTER_ERROR_INVALID_SERVICE_EP ErrorType = "Invalid service endpoint"

	MATTER_ERROR_INVALID_DIRECTORY_ENTRY_TYPE ErrorType = "Invalid directory entry type"

	MATTER_ERROR_FORCED_RESET ErrorType = "Service manager forced reset"

	MATTER_ERROR_NO_ENDPOINT ErrorType = "No endpoint was available to send the message"

	MATTER_ERROR_INVALID_DESTINATION_NODE_ID ErrorType = "Invalid destination node id"

	MATTER_ERROR_NOT_CONNECTED ErrorType = "Not connected"

	MATTER_ERROR_NO_SW_UPDATE_AVAILABLE ErrorType = "No SW update available"

	MATTER_ERROR_CA_CERT_NOT_FOUND ErrorType = "CA certificate not found"

	MATTER_ERROR_CERT_PATH_LEN_CONSTRAINT_EXCEEDED ErrorType = "Certificate path length constraint exceeded"

	MATTER_ERROR_CERT_PATH_TOO_LONG ErrorType = "Certificate path too long"

	MATTER_ERROR_CERT_USAGE_NOT_ALLOWED ErrorType = "Requested certificate usage is not allowed"

	MATTER_ERROR_CERT_EXPIRED ErrorType = "Certificate expired"

	MATTER_ERROR_CERT_NOT_VALID_YET ErrorType = "Certificate not yet valid"

	MATTER_ERROR_UNSUPPORTED_CERT_FORMAT ErrorType = "Unsupported certificate format"

	MATTER_ERROR_UNSUPPORTED_ELLIPTIC_CURVE ErrorType = "Unsupported elliptic curve"

	CHIP_CERT_NOT_USED ErrorType = "Certificate was not used in chain validation"

	MATTER_ERROR_CERT_NOT_FOUND ErrorType = "Certificate not found"

	MATTER_ERROR_INVALID_CASE_PARAMETER ErrorType = "Invalid CASE parameter"

	MATTER_ERROR_UNSUPPORTED_CASE_CONFIGURATION ErrorType = "Unsupported CASE configuration"

	MATTER_ERROR_CERT_LOAD_FAILED ErrorType = "Unable to load certificate"

	MATTER_ERROR_CERT_NOT_TRUSTED ErrorType = "Certificate not trusted"

	MATTER_ERROR_INVALID_ACCESS_TOKEN ErrorType = "Invalid access token"

	MATTER_ERROR_WRONG_CERT_DN ErrorType = "Wrong certificate distinguished name"

	MATTER_ERROR_INVALID_PROVISIONING_BUNDLE ErrorType = "Invalid provisioning bundle"

	MATTER_ERROR_PROVISIONING_BUNDLE_DECRYPTION_ERROR ErrorType = "Provisioning bundle decryption error"

	MATTER_ERROR_PASE_RECONFIGURE_REQUIRED ErrorType = "PASE reconfiguration required"

	MATTER_ERROR_WRONG_NODE_ID ErrorType = "Wrong node ID"

	MATTER_ERROR_CONN_ACCEPTED_ON_WRONG_PORT ErrorType = "Connection accepted on wrong port"

	MATTER_ERROR_CALLBACK_REPLACED ErrorType = "Application callback replaced"

	MATTER_ERROR_NO_CASE_AUTH_DELEGATE ErrorType = "No CASE auth delegate set"

	MATTER_ERROR_DEVICE_LOCATE_TIMEOUT ErrorType = "Timeout attempting to locate device"

	MATTER_ERROR_DEVICE_CONNECT_TIMEOUT ErrorType = "Timeout connecting to device"

	MATTER_ERROR_DEVICE_AUTH_TIMEOUT ErrorType = "Timeout authenticating device"

	MATTER_ERROR_MESSAGE_NOT_ACKNOWLEDGED ErrorType = "Message not acknowledged after max retries"

	MATTER_ERROR_RETRANS_TABLE_FULL ErrorType = "Retransmit Table is already full"

	MATTER_ERROR_INVALID_ACK_MESSAGE_COUNTER ErrorType = "Invalid acknowledged message counter"

	MATTER_ERROR_SEND_THROTTLED ErrorType = "Sending to peer is throttled on this Exchange"

	MATTER_ERROR_WRONG_MSG_VERSION_FOR_EXCHANGE ErrorType = "Message version not supported by current exchange context"

	MATTER_ERROR_UNSUPPORTED_CHIP_FEATURE ErrorType = "Required feature not supported by this configuration"

	ErrorUnsolicitedMsgNoOriginator ErrorType = "Unsolicited msg with originator bit clear"

	MATTER_ERROR_TOO_MANY_CONNECTIONS ErrorType = "Too many connections"

	MATTER_ERROR_SHUT_DOWN ErrorType = "The operation was cancelled because a shut down was initiated"

	MATTER_ERROR_CANCELLED ErrorType = "The operation has been cancelled"

	MATTER_ERROR_DRBG_ENTROPY_SOURCE_FAILED ErrorType = "DRBG entropy source failed to generate entropy data"

	MATTER_ERROR_MESSAGE_COUNTER_EXHAUSTED ErrorType = "Message counter exhausted"

	MATTER_ERROR_FABRIC_EXISTS ErrorType = "Trying to add a NOC for a fabric that already exists"

	MATTER_ERROR_KEY_NOT_FOUND_FROM_PEER ErrorType = "Key not found error code received from peer"

	MATTER_ERROR_WRONG_ENCRYPTION_TYPE_FROM_PEER ErrorType = "Wrong encryption type error code received from peer"

	MATTER_ERROR_UNKNOWN_KEY_TYPE_FROM_PEER ErrorType = "Unknown key type error code received from peer"

	MATTER_ERROR_INVALID_USE_OF_SESSION_KEY_FROM_PEER ErrorType = "Invalid use of session key error code received from peer"

	MATTER_ERROR_UNSUPPORTED_ENCRYPTION_TYPE_FROM_PEER ErrorType = "Unsupported encryption type error code received from peer"

	MATTER_ERROR_INTERNAL_KEY_ERROR_FROM_PEER ErrorType = "ErrorInternal key error code received from peer"

	MATTER_ERROR_INVALID_KEY_ID ErrorType = "Invalid key identifier"

	MATTER_ERROR_INVALID_TIME ErrorType = "Valid time value is not available"

	MATTER_ERROR_LOCKING_FAILURE ErrorType = "Failure to lock/unlock OS-provided lock"

	MATTER_ERROR_UNSUPPORTED_PASSCODE_CONFIG ErrorType = "Unsupported passcode encryption configuration"

	MATTER_ERROR_PASSCODE_AUTHENTICATION_FAILED ErrorType = "Passcode authentication failed"

	MATTER_ERROR_PASSCODE_FINGERPRINT_FAILED ErrorType = "Passcode fingerprint failed"

	MATTER_ERROR_SERIALIZATION_ELEMENT_NULL ErrorType = "Element requested is null"

	MATTER_ERROR_WRONG_CERT_SIGNATURE_ALGORITHM ErrorType = "Certificate not signed with required signature algorithm"

	MATTER_ERROR_WRONG_CHIP_SIGNATURE_ALGORITHM ErrorType = "CHIP signature not signed with required signature algorithm"

	MATTER_ERROR_SCHEMA_MISMATCH ErrorType = "Schema mismatch"

	MATTER_ERROR_INVALID_INTEGER_VALUE ErrorType = "Invalid integer value"

	MATTER_ERROR_CASE_RECONFIG_REQUIRED ErrorType = "CASE reconfiguration required"

	MATTER_ERROR_TOO_MANY_CASE_RECONFIGURATIONS ErrorType = "Too many CASE reconfigurations were received"

	MATTER_ERROR_BAD_REQUEST ErrorType = "Request cannot be processed or fulfilled"

	MATTER_ERROR_INVALID_MESSAGE_FLAG ErrorType = "Invalid message flag"

	MATTER_ERROR_KEY_EXPORT_RECONFIGURE_REQUIRED ErrorType = "Key export protocol required to reconfigure"

	MATTER_ERROR_NO_COMMON_KEY_EXPORT_CONFIGURATIONS ErrorType = "No supported key export protocol configurations in common"

	MATTER_ERROR_INVALID_KEY_EXPORT_CONFIGURATION ErrorType = "Invalid key export protocol configuration"

	MATTER_ERROR_NO_KEY_EXPORT_DELEGATE ErrorType = "No key export protocol delegate set"

	MATTER_ERROR_UNAUTHORIZED_KEY_EXPORT_REQUEST ErrorType = "Unauthorized key export request"

	MATTER_ERROR_UNAUTHORIZED_KEY_EXPORT_RESPONSE ErrorType = "Unauthorized key export response"

	MATTER_ERROR_EXPORTED_KEY_AUTHENTICATION_FAILED ErrorType = "Exported key authentication failed"

	MATTER_ERROR_TOO_MANY_SHARED_SESSION_END_NODES ErrorType = "Too many shared session end nodes"

	MATTER_ERROR_IM_MALFORMED_ATTRIBUTE_DATA_IB ErrorType = "Malformed Interaction Model Attribute Data IB"

	MATTER_ERROR_WRONG_CERT_TYPE ErrorType = "Wrong certificate type"

	MATTER_ERROR_DEFAULT_EVENT_HANDLER_NOT_CALLED ErrorType = "Default event handler not called"

	MATTER_ERROR_PERSISTED_STORAGE_FAILED ErrorType = "Persisted storage failed"

	MATTER_ERROR_PERSISTED_STORAGE_VALUE_NOT_FOUND ErrorType = "Raw not found in the persisted storage"

	MATTER_ERROR_IM_FABRIC_DELETED ErrorType = "The fabric is deleted, and the corresponding IM resources are released"

	MATTER_ERROR_PROFILE_STRING_CONTEXT_NOT_REGISTERED ErrorType = "String context not registered"

	MATTER_ERROR_INCOMPATIBLE_SCHEMA_VERSION ErrorType = "Incompatible data schema version"

	MATTER_ERROR_ACCESS_DENIED ErrorType = "The CHIP message is not granted access"

	MATTER_ERROR_UNKNOWN_RESOURCE_ID ErrorType = "Unknown resource ID"

	MATTER_ERROR_VERSION_MISMATCH ErrorType = "Version mismatch"

	MATTER_ERROR_UNSUPPORTED_THREAD_NETWORK_CREATE ErrorType = "Legacy device doesn't support standalone Thread network creation"

	MATTER_ERROR_INCONSISTENT_CONDITIONALITY ErrorType = "The Trait Instance is already being updated with a different conditionality"

	MATTER_ERROR_LOCAL_DATA_INCONSISTENT ErrorType = "The local data does not match any known version of the Trait Instance"

	CHIP_EVENT_ID_FOUND ErrorType = "Event ID matching criteria was found"

	ErrorInternal ErrorType = "ErrorInternal error"

	MATTER_ERROR_OPEN_FAILED ErrorType = "Open file failed"

	MATTER_ERROR_READ_FAILED ErrorType = "read from file failed"

	MATTER_ERROR_WRITE_FAILED ErrorType = "Write to file failed"

	MATTER_ERROR_DECODE_FAILED ErrorType = "Decoding failed"

	MATTER_ERROR_SESSION_KEY_SUSPENDED ErrorType = "Session key suspended"

	MATTER_ERROR_UNSUPPORTED_WIRELESS_REGULATORY_DOMAIN ErrorType = "Unsupported wireless regulatory domain"

	MATTER_ERROR_UNSUPPORTED_WIRELESS_OPERATING_LOCATION ErrorType = "Unsupported wireless operating location"

	MATTER_ERROR_MDNS_COLLISION ErrorType = "mDNS collision"

	MATTER_ERROR_IM_MALFORMED_ATTRIBUTE_PATH_IB ErrorType = "Malformed Interacton Model Attribute Path IB"

	MATTER_ERROR_IM_MALFORMED_EVENT_PATH_IB ErrorType = "Malformed Interacton Model Event Path IB"

	MATTER_ERROR_IM_MALFORMED_COMMAND_PATH_IB ErrorType = "Malformed Interacton Model Command Path IB"

	MATTER_ERROR_IM_MALFORMED_ATTRIBUTE_STATUS_IB ErrorType = "Malformed Interacton Model Attribute Status IB"

	MATTER_ERROR_IM_MALFORMED_EVENT_DATA_IB ErrorType = "Malformed Interacton Model Event Data IB"

	MATTER_ERROR_IM_MALFORMED_STATUS_IB ErrorType = "Malformed Interacton Model Status IB"

	MATTER_ERROR_PEER_NODE_NOT_FOUND ErrorType = "Unable to find the peer node"

	MATTER_ERROR_HSM ErrorType = "Hardware security module"

	MATTER_ERROR_INTERMEDIATE_CA_NOT_REQUIRED ErrorType = "Intermediate CA not required"

	MATTER_ERROR_REAL_TIME_NOT_SYNCED ErrorType = "Real time not synchronized"

	MATTER_ERROR_UNEXPECTED_EVENT ErrorType = "Unexpected event"

	MATTER_ERROR_ENDPOINT_POOL_FULL ErrorType = "Endpoint pool full"

	MATTER_ERROR_INBOUND_MESSAGE_TOO_BIG ErrorType = "Inbound message too big"

	MATTER_ERROR_OUTBOUND_MESSAGE_TOO_BIG ErrorType = "Outbound message too big"

	MATTER_ERROR_DUPLICATE_MESSAGE_RECEIVED ErrorType = "Duplicate message received"

	MATTER_ERROR_INVALID_PUBLIC_KEY ErrorType = "Invalid public key"

	MATTER_ERROR_FABRIC_MISMATCH_ON_ICA ErrorType = "Fabric mismatch on ICA"

	MATTER_ERROR_MESSAGE_COUNTER_OUT_OF_WINDOW ErrorType = "Message id out of window"

	MATTER_ERROR_REBOOT_SIGNAL_RECEIVED ErrorType = "Termination signal is received"

	MATTER_ERROR_NO_SHARED_TRUSTED_ROOT ErrorType = "No shared trusted root"

	MATTER_ERROR_IM_STATUS_CODE_RECEIVED ErrorType = "Interaction Model Error"

	MATTER_ERROR_IM_MALFORMED_COMMAND_STATUS_IB ErrorType = "Malformed Interaction Model Command Status IB"

	MATTER_ERROR_IM_MALFORMED_INVOKE_RESPONSE_IB ErrorType = "Malformed Interaction Model Invoke Response IB"

	MATTER_ERROR_IM_MALFORMED_INVOKE_REQUEST_MESSAGE ErrorType = "Malformed Interaction Model Invoke Request Message"

	MATTER_ERROR_IM_MALFORMED_INVOKE_RESPONSE_MESSAGE ErrorType = "Malformed Interaction Model Invoke Response Message"

	MATTER_ERROR_IM_MALFORMED_ATTRIBUTE_REPORT_MESSAGE ErrorType = "Malformed Interaction Model Attribute Report Message"

	MATTER_ERROR_IM_MALFORMED_WRITE_REQUEST_MESSAGE ErrorType = "Malformed Interaction Model Write Request Message"

	MATTER_ERROR_IM_MALFORMED_EVENT_FILTER_IB ErrorType = "Malformed Interaction Model Event Filter IB"

	MATTER_ERROR_IM_MALFORMED_READ_REQUEST_MESSAGE ErrorType = "Malformed Interaction Model read Request Message"

	MATTER_ERROR_IM_MALFORMED_SUBSCRIBE_REQUEST_MESSAGE ErrorType = "Malformed Interaction Model Subscribe Request Message"

	MATTER_ERROR_IM_MALFORMED_SUBSCRIBE_RESPONSE_MESSAGE ErrorType = "Malformed Interaction Model Subscribe Response Message"

	MATTER_ERROR_IM_MALFORMED_EVENT_REPORT_IB ErrorType = "Malformed Interaction Model Event Report IB"

	MATTER_ERROR_IM_MALFORMED_CLUSTER_PATH_IB ErrorType = "Malformed Interaction Model Cluster Path IB"

	MATTER_ERROR_IM_MALFORMED_DATA_VERSION_FILTER_IB ErrorType = "Malformed Interaction Model Data Version Filter IB"

	MATTER_ERROR_NOT_FOUND ErrorType = "The item referenced in the function call was not found"

	MATTER_ERROR_IM_MALFORMED_TIMED_REQUEST_MESSAGE ErrorType = "Malformed Interaction Model Timed Request Message"

	MATTER_ERROR_INVALID_FILE_IDENTIFIER ErrorType = "The file identifier, encoded in the first few bytes of a processed file, has unexpected value"

	MATTER_ERROR_BUSY ErrorType = "The Resource is busy and cannot process the request"

	MATTER_ERROR_MAX_RETRY_EXCEEDED ErrorType = "The maximum retry limit has been exceeded"

	MATTER_ERROR_PROVIDER_LIST_EXHAUSTED ErrorType = "The provider list has been exhausted"

	MATTER_ERROR_ANOTHER_COMMISSIONING_IN_PROGRESS ErrorType = "Another commissioning in progress"

	MATTER_ERROR_INVALID_SCHEME_PREFIX ErrorType = "The scheme field contains an invalid prefix"

	MATTER_ERROR_MISSING_URI_SEPARATOR ErrorType = "The URI separator is missing"

	ErrorInvalidArgument              ErrorType = "Invalid Argument"
	ErrorIncorrectState               ErrorType = " Incorrect State"
	ErrorNotImplemented               ErrorType = " Not Implemented"
	ErrorNotMemory                    ErrorType = "Not memory"
	DeviceErrorConfigNotFound         ErrorType = "CHIP_DEVICE_ERROR_CONFIG_NOT_FOUND"
	ErrorDuplicateMessageReceived     ErrorType = " MATTER_ERROR_DUPLICATE_MESSAGE_RECEIVED"
	ErrorInvalidFabricIndex           ErrorType = " ErrorInvalidFabricIndex"
	TooManyUnsolicitedMessageHandlers ErrorType = " MATTER_ERROR_TOO_MANY_UNSOLICITED_MESSAGE_HANDLERS"
	ErrorWrongEncryptionTypeFromPeer  ErrorType = " MATTER_ERROR_WRONG_ENCRYPTION_TYPE_FROM_PEER"
	InvalidMessageType                ErrorType = " ErrorInvalidMessageType"
	ErrorInvalidMessageType           ErrorType = "ErrorInternal error"
	ErrorWrongTlvType                 ErrorType = " MATTER_ERROR_WRONG_TLV_TYPE"
	UnexpectedTlvElement              ErrorType = " MATTER_ERROR_UNEXPECTED_TLV_ELEMENT"
	InvalidCaseParameter              ErrorType = " MATTER_ERROR_INVALID_CASE_PARAMETER"
	InvalidTlvTag                     ErrorType = " MATTER_ERROR_INVALID_TLV_TAG"
	ErrorKeyNotFound                  ErrorType = " MATTER_ERROR_KEY_NOT_FOUND"
	ErrorNoUnsolicitedMessageHandler  ErrorType = " MATTER_ERROR_NO_UNSOLICITED_MESSAGE_HANDLER"
)

type MatterError struct {
	ErrorType ErrorType
	Mod       string
	Msg       string
}

func New(typ ErrorType, args ...string) MatterError {
	e := MatterError{
		ErrorType: typ,
	}
	if len(args) > 0 {
		e.Mod = args[0]
		if len(args) > 1 {
			e.Msg = args[1]
		}
	}
	return e
}

func (e ErrorType) Error() string {
	value := string(e)
	return value
}

func (e MatterError) Error() string {
	var str = e.ErrorType.Error()
	if e.Mod != "" {
		str = e.Mod + ":" + str
	}
	if e.Msg != "" {
		str = str + ":" + e.Msg
	}
	return str
}
