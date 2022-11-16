package lib

import "fmt"

type ActionId uint8
type AttributeId uint32
type ClusterId uint32
type ClusterStatus uint8
type CommandId uint32
type CompressedFabricId uint64
type DataVersion uint32
type DeviceTypeId uint32
type EndpointId uint16
type EventId uint32
type EventNumber uint64
type FabricId uint64
type FabricIndex uint8
type FieldId uint32
type ListIndex uint16
type TransactionId uint32
type KeysetId uint16
type InteractionModelRevision uint8
type SubscriptionId uint32

const InvalidKeysetId KeysetId = 0xFFFF

const InvalidClusterId ClusterId = 0xFFFF_FFFF
const InvalidAttributeId AttributeId = 0xFFFF_FFFF
const InvalidCommandId CommandId = 0xFFFF_FFFF
const InvalidEventId EventId = 0xFFFF_FFFF
const InvalidFieldId FieldId = 0xFFFF_FFFF

func (f FabricId) String() string {
	value := uint64(f)
	return fmt.Sprintf("%016X", value)
}

func (f CompressedFabricId) String() string {
	value := uint64(f)
	return fmt.Sprintf("%016X", value)
}

func UndefinedFabricIndex() FabricIndex {
	return kUndefinedFabricIndex
}
