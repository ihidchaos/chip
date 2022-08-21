package lib

const (
	UndefinedFabricIndex FabricIndex = 0
)

type CompressedFabricId uint64
type FabricId uint64
type FabricIndex uint8
type KeysetId uint16

const KInvalidKeysetId KeysetId = 0xFFFF

type ProductId uint16

type OperationalNodeId uint64

type GroupNodeID uint64

type TemporaryLocalNodeId uint64

const KUndefinedFabricIndex FabricIndex = 0

const KMinValidFabricIndex FabricIndex = 1

const KMaxValidFabricIndex FabricIndex = 0xfe

func (index FabricIndex) IsValidFabricIndex() bool {
	return index >= KMinValidFabricIndex && index <= KMaxValidFabricIndex
}
