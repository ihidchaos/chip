package lib

type MatterError string

var (
	MATTER_ERROR_SENDING_BLOCKED MatterError = "Sending blocked"

	MATTER_ERROR_CONNECTION_ABORTED MatterError = "Connection aborted"

	MATTER_ERROR_INCORRECT_STATE MatterError = "Incorrect state"

	MATTER_ERROR_MESSAGE_TOO_LONG MatterError = "Message too long"

	MATTER_ERROR_UNSUPPORTED_EXCHANGE_VERSION MatterError = "Unsupported exchange version"

	MATTER_ERROR_TOO_MANY_UNSOLICITED_MESSAGE_HANDLERS MatterError = "Too many unsolicited message handlers"

	MATTER_ERROR_NO_UNSOLICITED_MESSAGE_HANDLER MatterError = "No unsolicited message handler"

	MATTER_ERROR_NO_CONNECTION_HANDLER MatterError = "No connection handler"

	MATTER_ERROR_TOO_MANY_PEER_NODES MatterError = "Too many peer nodes"

	MATTER_ERROR_SENTINEL MatterError = "Internal sentinel"

	MATTER_ERROR_NO_MEMORY MatterError = "No memory"

	MATTER_ERROR_NO_MESSAGE_HANDLER MatterError = "No message handler"

	MATTER_ERROR_MESSAGE_INCOMPLETE MatterError = "Message incomplete"

	MATTER_ERROR_DATA_NOT_ALIGNED MatterError = "Data not aligned"

	MATTER_ERROR_UNKNOWN_KEY_TYPE MatterError = "Unknown key type"

	MATTER_ERROR_KEY_NOT_FOUND MatterError = "Key not found"

	MATTER_ERROR_WRONG_ENCRYPTION_TYPE MatterError = "Wrong encryption type"

	MATTER_ERROR_TOO_MANY_KEYS MatterError = "Too many keys"

	MATTER_ERROR_INTEGRITY_CHECK_FAILED MatterError = "Integrity check failed"

	MATTER_ERROR_INVALID_SIGNATURE MatterError = "Invalid signature"

	MATTER_ERROR_UNSUPPORTED_MESSAGE_VERSION MatterError = "Unsupported message version"

	MATTER_ERROR_UNSUPPORTED_ENCRYPTION_TYPE MatterError = "Unsupported encryption type"

	MATTER_ERROR_UNSUPPORTED_SIGNATURE_TYPE MatterError = "Unsupported signature type"

	MATTER_ERROR_INVALID_MESSAGE_LENGTH MatterError = "Invalid message length"

	MATTER_ERROR_BUFFER_TOO_SMALL MatterError = "Buffer too small"

	MATTER_ERROR_DUPLICATE_KEY_ID MatterError = "Duplicate key id"

	MATTER_ERROR_WRONG_KEY_TYPE MatterError = "Wrong key type"

	MATTER_ERROR_WELL_UNINITIALIZED MatterError = "Well uninitialized"

	MATTER_ERROR_WELL_EMPTY MatterError = "Well empty"

	MATTER_ERROR_INVALID_STRING_LENGTH MatterError = "Invalid string length"

	MATTER_ERROR_INVALID_LIST_LENGTH MatterError = "invalid list length"

	MATTER_ERROR_INVALID_INTEGRITY_TYPE MatterError = "Invalid integrity type"

	MATTER_END_OF_TLV MatterError = "End of TLV"

	MATTER_ERROR_TLV_UNDERRUN MatterError = "TLV underrun"

	MATTER_ERROR_INVALID_TLV_ELEMENT MatterError = "Invalid TLV element"

	MATTER_ERROR_INVALID_TLV_TAG MatterError = "Invalid TLV tag"

	MATTER_ERROR_UNKNOWN_IMPLICIT_TLV_TAG MatterError = "Unknown implicit TLV tag"

	MATTER_ERROR_WRONG_TLV_TYPE MatterError = "Wrong TLV type"

	MATTER_ERROR_TLV_CONTAINER_OPEN MatterError = "TLV container open"

	MATTER_ERROR_INVALID_TRANSFER_MODE MatterError = "Invalid transfer mode"

	MATTER_ERROR_INVALID_PROFILE_ID MatterError = "Invalid profile id"

	MATTER_ERROR_INVALID_MESSAGE_TYPE MatterError = "Invalid message type"

	MATTER_ERROR_UNEXPECTED_TLV_ELEMENT MatterError = "Unexpected TLV element"

	MATTER_ERROR_STATUS_REPORT_RECEIVED MatterError = "Status Report received from peer"

	MATTER_ERROR_NOT_IMPLEMENTED MatterError = "Not Implemented"

	MATTER_ERROR_INVALID_ADDRESS MatterError = "Invalid address"

	MATTER_ERROR_INVALID_ARGUMENT MatterError = "Invalid argument"

	MATTER_ERROR_TLV_TAG_NOT_FOUND MatterError = "TLV tag not found"

	MATTER_ERROR_MISSING_SECURE_SESSION MatterError = "Missing secure session"

	MATTER_ERROR_INVALID_ADMIN_SUBJECT MatterError = "CaseAdminSubject is not valid"

	MATTER_ERROR_INSUFFICIENT_PRIVILEGE MatterError = "Required privilege was insufficient during an operation"

	MATTER_ERROR_IM_MALFORMED_ATTRIBUTE_REPORT_IB MatterError = "Malformed Interacton Model Attribute Report IB"

	MATTER_ERROR_IM_MALFORMED_COMMAND_DATA_IB MatterError = "Malformed Interacton Model Command Data IB"

	MATTER_ERROR_IM_MALFORMED_EVENT_STATUS_IB MatterError = "Malformed Interacton Model Event Status IB"

	MATTER_ERROR_IM_MALFORMED_STATUS_RESPONSE_MESSAGE MatterError = "Malformed Interacton Model Status Response IB"

	MATTER_ERROR_INVALID_PATH_LIST MatterError = "Invalid TLV path list"

	MATTER_ERROR_INVALID_DATA_LIST MatterError = "Invalid TLV data list"

	MATTER_ERROR_TRANSACTION_CANCELED MatterError = "Transaction canceled"

	MATTER_ERROR_LISTENER_ALREADY_STARTED MatterError = "Listener already started"

	MATTER_ERROR_LISTENER_ALREADY_STOPPED MatterError = "Listener already stopped"

	MATTER_ERROR_INVALID_SUBSCRIPTION MatterError = "Invalid Subscription Id"

	MATTER_ERROR_TIMEOUT MatterError = "Timeout"

	MATTER_ERROR_INVALID_DEVICE_DESCRIPTOR MatterError = "Invalid device descriptor"

	MATTER_ERROR_UNSUPPORTED_DEVICE_DESCRIPTOR_VERSION MatterError = "Unsupported device descriptor version"

	CHIP_END_OF_INPUT MatterError = "End of input"

	MATTER_ERROR_RATE_LIMIT_EXCEEDED MatterError = "Rate limit exceeded"

	MATTER_ERROR_SECURITY_MANAGER_BUSY MatterError = "Security manager busy"

	MATTER_ERROR_INVALID_PASE_PARAMETER MatterError = "Invalid PASE parameter"

	MATTER_ERROR_PASE_SUPPORTS_ONLY_CONFIG1 MatterError = "PASE supports only Config1"

	MATTER_ERROR_NO_COMMON_PASE_CONFIGURATIONS MatterError = "No supported PASE configurations in common"

	MATTER_ERROR_INVALID_PASE_CONFIGURATION MatterError = "Invalid PASE configuration"

	MATTER_ERROR_KEY_CONFIRMATION_FAILED MatterError = "Key confirmation failed"

	MATTER_ERROR_INVALID_USE_OF_SESSION_KEY MatterError = "Invalid use of session key"

	MATTER_ERROR_CONNECTION_CLOSED_UNEXPECTEDLY MatterError = "Connection closed unexpectedly"

	MATTER_ERROR_MISSING_TLV_ELEMENT MatterError = "Missing TLV element"

	MATTER_ERROR_RANDOM_DATA_UNAVAILABLE MatterError = "Random data unavailable"

	MATTER_ERROR_UNSUPPORTED_HOST_PORT_ELEMENT MatterError = "Unsupported type in host/port list"

	MATTER_ERROR_INVALID_HOST_SUFFIX_INDEX MatterError = "Invalid suffix index in host/port list"

	MATTER_ERROR_HOST_PORT_LIST_EMPTY MatterError = "Host/port empty"

	MATTER_ERROR_UNSUPPORTED_AUTH_MODE MatterError = "Unsupported authentication mode"

	MATTER_ERROR_INVALID_SERVICE_EP MatterError = "Invalid service endpoint"

	MATTER_ERROR_INVALID_DIRECTORY_ENTRY_TYPE MatterError = "Invalid directory entry type"

	MATTER_ERROR_FORCED_RESET MatterError = "Service manager forced reset"

	MATTER_ERROR_NO_ENDPOINT MatterError = "No endpoint was available to send the message"

	MATTER_ERROR_INVALID_DESTINATION_NODE_ID MatterError = "Invalid destination node id"

	MATTER_ERROR_NOT_CONNECTED MatterError = "Not connected"

	MATTER_ERROR_NO_SW_UPDATE_AVAILABLE MatterError = "No SW update available"

	MATTER_ERROR_CA_CERT_NOT_FOUND MatterError = "CA certificate not found"

	MATTER_ERROR_CERT_PATH_LEN_CONSTRAINT_EXCEEDED MatterError = "Certificate path length constraint exceeded"

	MATTER_ERROR_CERT_PATH_TOO_LONG MatterError = "Certificate path too long"

	MATTER_ERROR_CERT_USAGE_NOT_ALLOWED MatterError = "Requested certificate usage is not allowed"

	MATTER_ERROR_CERT_EXPIRED MatterError = "Certificate expired"

	MATTER_ERROR_CERT_NOT_VALID_YET MatterError = "Certificate not yet valid"

	MATTER_ERROR_UNSUPPORTED_CERT_FORMAT MatterError = "Unsupported certificate format"

	MATTER_ERROR_UNSUPPORTED_ELLIPTIC_CURVE MatterError = "Unsupported elliptic curve"

	CHIP_CERT_NOT_USED MatterError = "Certificate was not used in chain validation"

	MATTER_ERROR_CERT_NOT_FOUND MatterError = "Certificate not found"

	MATTER_ERROR_INVALID_CASE_PARAMETER MatterError = "Invalid CASE parameter"

	MATTER_ERROR_UNSUPPORTED_CASE_CONFIGURATION MatterError = "Unsupported CASE configuration"

	MATTER_ERROR_CERT_LOAD_FAILED MatterError = "Unable to load certificate"

	MATTER_ERROR_CERT_NOT_TRUSTED MatterError = "Certificate not trusted"

	MATTER_ERROR_INVALID_ACCESS_TOKEN MatterError = "Invalid access token"

	MATTER_ERROR_WRONG_CERT_DN MatterError = "Wrong certificate distinguished name"

	MATTER_ERROR_INVALID_PROVISIONING_BUNDLE MatterError = "Invalid provisioning bundle"

	MATTER_ERROR_PROVISIONING_BUNDLE_DECRYPTION_ERROR MatterError = "Provisioning bundle decryption error"

	MATTER_ERROR_PASE_RECONFIGURE_REQUIRED MatterError = "PASE reconfiguration required"

	MATTER_ERROR_WRONG_NODE_ID MatterError = "Wrong node ID"

	MATTER_ERROR_CONN_ACCEPTED_ON_WRONG_PORT MatterError = "Connection accepted on wrong port"

	MATTER_ERROR_CALLBACK_REPLACED MatterError = "Application callback replaced"

	MATTER_ERROR_NO_CASE_AUTH_DELEGATE MatterError = "No CASE auth delegate set"

	MATTER_ERROR_DEVICE_LOCATE_TIMEOUT MatterError = "Timeout attempting to locate device"

	MATTER_ERROR_DEVICE_CONNECT_TIMEOUT MatterError = "Timeout connecting to device"

	MATTER_ERROR_DEVICE_AUTH_TIMEOUT MatterError = "Timeout authenticating device"

	MATTER_ERROR_MESSAGE_NOT_ACKNOWLEDGED MatterError = "Message not acknowledged after max retries"

	MATTER_ERROR_RETRANS_TABLE_FULL MatterError = "Retransmit Table is already full"

	MATTER_ERROR_INVALID_ACK_MESSAGE_COUNTER MatterError = "Invalid acknowledged message counter"

	MATTER_ERROR_SEND_THROTTLED MatterError = "Sending to peer is throttled on this Exchange"

	MATTER_ERROR_WRONG_MSG_VERSION_FOR_EXCHANGE MatterError = "Message version not supported by current exchange context"

	MATTER_ERROR_UNSUPPORTED_CHIP_FEATURE MatterError = "Required feature not supported by this configuration"

	MATTER_ERROR_UNSOLICITED_MSG_NO_ORIGINATOR MatterError = "Unsolicited msg with originator bit clear"

	MATTER_ERROR_INVALID_FABRIC_INDEX MatterError = "Invalid Fabric Index"

	MATTER_ERROR_TOO_MANY_CONNECTIONS MatterError = "Too many connections"

	MATTER_ERROR_SHUT_DOWN MatterError = "The operation was cancelled because a shut down was initiated"

	MATTER_ERROR_CANCELLED MatterError = "The operation has been cancelled"

	MATTER_ERROR_DRBG_ENTROPY_SOURCE_FAILED MatterError = "DRBG entropy source failed to generate entropy data"

	MATTER_ERROR_MESSAGE_COUNTER_EXHAUSTED MatterError = "Message counter exhausted"

	MATTER_ERROR_FABRIC_EXISTS MatterError = "Trying to add a NOC for a fabric that already exists"

	MATTER_ERROR_KEY_NOT_FOUND_FROM_PEER MatterError = "Key not found error code received from peer"

	MATTER_ERROR_WRONG_ENCRYPTION_TYPE_FROM_PEER MatterError = "Wrong encryption type error code received from peer"

	MATTER_ERROR_UNKNOWN_KEY_TYPE_FROM_PEER MatterError = "Unknown key type error code received from peer"

	MATTER_ERROR_INVALID_USE_OF_SESSION_KEY_FROM_PEER MatterError = "Invalid use of session key error code received from peer"

	MATTER_ERROR_UNSUPPORTED_ENCRYPTION_TYPE_FROM_PEER MatterError = "Unsupported encryption type error code received from peer"

	MATTER_ERROR_INTERNAL_KEY_ERROR_FROM_PEER MatterError = "Internal key error code received from peer"

	MATTER_ERROR_INVALID_KEY_ID MatterError = "Invalid key identifier"

	MATTER_ERROR_INVALID_TIME MatterError = "Valid time value is not available"

	MATTER_ERROR_LOCKING_FAILURE MatterError = "Failure to lock/unlock OS-provided lock"

	MATTER_ERROR_UNSUPPORTED_PASSCODE_CONFIG MatterError = "Unsupported passcode encryption configuration"

	MATTER_ERROR_PASSCODE_AUTHENTICATION_FAILED MatterError = "Passcode authentication failed"

	MATTER_ERROR_PASSCODE_FINGERPRINT_FAILED MatterError = "Passcode fingerprint failed"

	MATTER_ERROR_SERIALIZATION_ELEMENT_NULL MatterError = "Element requested is null"

	MATTER_ERROR_WRONG_CERT_SIGNATURE_ALGORITHM MatterError = "Certificate not signed with required signature algorithm"

	MATTER_ERROR_WRONG_CHIP_SIGNATURE_ALGORITHM MatterError = "CHIP signature not signed with required signature algorithm"

	MATTER_ERROR_SCHEMA_MISMATCH MatterError = "Schema mismatch"

	MATTER_ERROR_INVALID_INTEGER_VALUE MatterError = "Invalid integer value"

	MATTER_ERROR_CASE_RECONFIG_REQUIRED MatterError = "CASE reconfiguration required"

	MATTER_ERROR_TOO_MANY_CASE_RECONFIGURATIONS MatterError = "Too many CASE reconfigurations were received"

	MATTER_ERROR_BAD_REQUEST MatterError = "Request cannot be processed or fulfilled"

	MATTER_ERROR_INVALID_MESSAGE_FLAG MatterError = "Invalid message flag"

	MATTER_ERROR_KEY_EXPORT_RECONFIGURE_REQUIRED MatterError = "Key export protocol required to reconfigure"

	MATTER_ERROR_NO_COMMON_KEY_EXPORT_CONFIGURATIONS MatterError = "No supported key export protocol configurations in common"

	MATTER_ERROR_INVALID_KEY_EXPORT_CONFIGURATION MatterError = "Invalid key export protocol configuration"

	MATTER_ERROR_NO_KEY_EXPORT_DELEGATE MatterError = "No key export protocol delegate set"

	MATTER_ERROR_UNAUTHORIZED_KEY_EXPORT_REQUEST MatterError = "Unauthorized key export request"

	MATTER_ERROR_UNAUTHORIZED_KEY_EXPORT_RESPONSE MatterError = "Unauthorized key export response"

	MATTER_ERROR_EXPORTED_KEY_AUTHENTICATION_FAILED MatterError = "Exported key authentication failed"

	MATTER_ERROR_TOO_MANY_SHARED_SESSION_END_NODES MatterError = "Too many shared session end nodes"

	MATTER_ERROR_IM_MALFORMED_ATTRIBUTE_DATA_IB MatterError = "Malformed Interaction Model Attribute Data IB"

	MATTER_ERROR_WRONG_CERT_TYPE MatterError = "Wrong certificate type"

	MATTER_ERROR_DEFAULT_EVENT_HANDLER_NOT_CALLED MatterError = "Default event handler not called"

	MATTER_ERROR_PERSISTED_STORAGE_FAILED MatterError = "Persisted storage failed"

	MATTER_ERROR_PERSISTED_STORAGE_VALUE_NOT_FOUND MatterError = "Value not found in the persisted storage"

	MATTER_ERROR_IM_FABRIC_DELETED MatterError = "The fabric is deleted, and the corresponding IM resources are released"

	MATTER_ERROR_PROFILE_STRING_CONTEXT_NOT_REGISTERED MatterError = "String context not registered"

	MATTER_ERROR_INCOMPATIBLE_SCHEMA_VERSION MatterError = "Incompatible data schema version"

	MATTER_ERROR_ACCESS_DENIED MatterError = "The CHIP message is not granted access"

	MATTER_ERROR_UNKNOWN_RESOURCE_ID MatterError = "Unknown resource ID"

	MATTER_ERROR_VERSION_MISMATCH MatterError = "Version mismatch"

	MATTER_ERROR_UNSUPPORTED_THREAD_NETWORK_CREATE MatterError = "Legacy device doesn't support standalone Thread network creation"

	MATTER_ERROR_INCONSISTENT_CONDITIONALITY MatterError = "The Trait Instance is already being updated with a different conditionality"

	MATTER_ERROR_LOCAL_DATA_INCONSISTENT MatterError = "The local data does not match any known version of the Trait Instance"

	CHIP_EVENT_ID_FOUND MatterError = "Event ID matching criteria was found"

	MATTER_ERROR_INTERNAL MatterError = "Internal error"

	MATTER_ERROR_OPEN_FAILED MatterError = "Open file failed"

	MATTER_ERROR_READ_FAILED MatterError = "read from file failed"

	MATTER_ERROR_WRITE_FAILED MatterError = "Write to file failed"

	MATTER_ERROR_DECODE_FAILED MatterError = "Decoding failed"

	MATTER_ERROR_SESSION_KEY_SUSPENDED MatterError = "Session key suspended"

	MATTER_ERROR_UNSUPPORTED_WIRELESS_REGULATORY_DOMAIN MatterError = "Unsupported wireless regulatory domain"

	MATTER_ERROR_UNSUPPORTED_WIRELESS_OPERATING_LOCATION MatterError = "Unsupported wireless operating location"

	MATTER_ERROR_MDNS_COLLISION MatterError = "mDNS collision"

	MATTER_ERROR_IM_MALFORMED_ATTRIBUTE_PATH_IB MatterError = "Malformed Interacton Model Attribute Path IB"

	MATTER_ERROR_IM_MALFORMED_EVENT_PATH_IB MatterError = "Malformed Interacton Model Event Path IB"

	MATTER_ERROR_IM_MALFORMED_COMMAND_PATH_IB MatterError = "Malformed Interacton Model Command Path IB"

	MATTER_ERROR_IM_MALFORMED_ATTRIBUTE_STATUS_IB MatterError = "Malformed Interacton Model Attribute Status IB"

	MATTER_ERROR_IM_MALFORMED_EVENT_DATA_IB MatterError = "Malformed Interacton Model Event Data IB"

	MATTER_ERROR_IM_MALFORMED_STATUS_IB MatterError = "Malformed Interacton Model Status IB"

	MATTER_ERROR_PEER_NODE_NOT_FOUND MatterError = "Unable to find the peer node"

	MATTER_ERROR_HSM MatterError = "Hardware security module"

	MATTER_ERROR_INTERMEDIATE_CA_NOT_REQUIRED MatterError = "Intermediate CA not required"

	MATTER_ERROR_REAL_TIME_NOT_SYNCED MatterError = "Real time not synchronized"

	MATTER_ERROR_UNEXPECTED_EVENT MatterError = "Unexpected event"

	MATTER_ERROR_ENDPOINT_POOL_FULL MatterError = "Endpoint pool full"

	MATTER_ERROR_INBOUND_MESSAGE_TOO_BIG MatterError = "Inbound message too big"

	MATTER_ERROR_OUTBOUND_MESSAGE_TOO_BIG MatterError = "Outbound message too big"

	MATTER_ERROR_DUPLICATE_MESSAGE_RECEIVED MatterError = "Duplicate message received"

	MATTER_ERROR_INVALID_PUBLIC_KEY MatterError = "Invalid public key"

	MATTER_ERROR_FABRIC_MISMATCH_ON_ICA MatterError = "Fabric mismatch on ICA"

	MATTER_ERROR_MESSAGE_COUNTER_OUT_OF_WINDOW MatterError = "Message id out of window"

	MATTER_ERROR_REBOOT_SIGNAL_RECEIVED MatterError = "Termination signal is received"

	MATTER_ERROR_NO_SHARED_TRUSTED_ROOT MatterError = "No shared trusted root"

	MATTER_ERROR_IM_STATUS_CODE_RECEIVED MatterError = "Interaction Model Error"

	MATTER_ERROR_IM_MALFORMED_COMMAND_STATUS_IB MatterError = "Malformed Interaction Model Command Status IB"

	MATTER_ERROR_IM_MALFORMED_INVOKE_RESPONSE_IB MatterError = "Malformed Interaction Model Invoke Response IB"

	MATTER_ERROR_IM_MALFORMED_INVOKE_REQUEST_MESSAGE MatterError = "Malformed Interaction Model Invoke Request Message"

	MATTER_ERROR_IM_MALFORMED_INVOKE_RESPONSE_MESSAGE MatterError = "Malformed Interaction Model Invoke Response Message"

	MATTER_ERROR_IM_MALFORMED_ATTRIBUTE_REPORT_MESSAGE MatterError = "Malformed Interaction Model Attribute Report Message"

	MATTER_ERROR_IM_MALFORMED_WRITE_REQUEST_MESSAGE MatterError = "Malformed Interaction Model Write Request Message"

	MATTER_ERROR_IM_MALFORMED_EVENT_FILTER_IB MatterError = "Malformed Interaction Model Event Filter IB"

	MATTER_ERROR_IM_MALFORMED_READ_REQUEST_MESSAGE MatterError = "Malformed Interaction Model read Request Message"

	MATTER_ERROR_IM_MALFORMED_SUBSCRIBE_REQUEST_MESSAGE MatterError = "Malformed Interaction Model Subscribe Request Message"

	MATTER_ERROR_IM_MALFORMED_SUBSCRIBE_RESPONSE_MESSAGE MatterError = "Malformed Interaction Model Subscribe Response Message"

	MATTER_ERROR_IM_MALFORMED_EVENT_REPORT_IB MatterError = "Malformed Interaction Model Event Report IB"

	MATTER_ERROR_IM_MALFORMED_CLUSTER_PATH_IB MatterError = "Malformed Interaction Model Cluster Path IB"

	MATTER_ERROR_IM_MALFORMED_DATA_VERSION_FILTER_IB MatterError = "Malformed Interaction Model Data Version Filter IB"

	MATTER_ERROR_NOT_FOUND MatterError = "The item referenced in the function call was not found"

	MATTER_ERROR_IM_MALFORMED_TIMED_REQUEST_MESSAGE MatterError = "Malformed Interaction Model Timed Request Message"

	MATTER_ERROR_INVALID_FILE_IDENTIFIER MatterError = "The file identifier, encoded in the first few bytes of a processed file, has unexpected value"

	MATTER_ERROR_BUSY MatterError = "The Resource is busy and cannot process the request"

	MATTER_ERROR_MAX_RETRY_EXCEEDED MatterError = "The maximum retry limit has been exceeded"

	MATTER_ERROR_PROVIDER_LIST_EXHAUSTED MatterError = "The provider list has been exhausted"

	MATTER_ERROR_ANOTHER_COMMISSIONING_IN_PROGRESS MatterError = "Another commissioning in progress"

	MATTER_ERROR_INVALID_SCHEME_PREFIX MatterError = "The scheme field contains an invalid prefix"

	MATTER_ERROR_MISSING_URI_SEPARATOR MatterError = "The URI separator is missing"

	InvalidArgument                   MatterError = " MATTER_ERROR_INVALID_ARGUMENT"
	IncorrectState                    MatterError = " MATTER_ERROR_INCORRECT_STATE"
	NotImplemented                    MatterError = " MATTER_ERROR_NOT_IMPLEMENTED"
	NotMemory                         MatterError = "Not memory"
	ErrorInternal                     MatterError = " MATTER_ERROR_INTERNAL"
	DeviceErrorConfigNotFound         MatterError = "CHIP_DEVICE_ERROR_CONFIG_NOT_FOUND"
	DuplicateMessageReceived          MatterError = " MATTER_ERROR_DUPLICATE_MESSAGE_RECEIVED"
	InvalidFabricIndex                MatterError = " MATTER_ERROR_INVALID_FABRIC_INDEX"
	TooManyUnsolicitedMessageHandlers MatterError = " MATTER_ERROR_TOO_MANY_UNSOLICITED_MESSAGE_HANDLERS"
	WrongEncryptionTypeFromPeer       MatterError = " MATTER_ERROR_WRONG_ENCRYPTION_TYPE_FROM_PEER"
	InvalidMessageType                MatterError = " MATTER_ERROR_INVALID_MESSAGE_TYPE"
	WrongTlvType                      MatterError = " MATTER_ERROR_WRONG_TLV_TYPE"
	UnexpectedTlvElement              MatterError = " MATTER_ERROR_UNEXPECTED_TLV_ELEMENT"
	InvalidCaseParameter              MatterError = " MATTER_ERROR_INVALID_CASE_PARAMETER"
	InvalidTlvTag                     MatterError = " MATTER_ERROR_INVALID_TLV_TAG"
	KeyNotFound                       MatterError = " MATTER_ERROR_KEY_NOT_FOUND"
	NoUnsolicitedMessageHandler       MatterError = " MATTER_ERROR_NO_UNSOLICITED_MESSAGE_HANDLER"
)

func (e MatterError) Error() string {
	value := string(e)
	return value
}
