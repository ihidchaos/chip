package lib

const (
	kUndefinedFabricIndex FabricIndex = 0
)

type NodeId uint64

const KUndefinedNodeId NodeId = 0

//0xFFF1 Test Vendor #1
//0xFFF2 Test Vendor #2
//0xFFF3 Test Vendor #3
//0xFFF4 Test Vendor #4

type CompressedFabricId uint64
type FabricId uint64
type FabricIndex uint8

type VendorId uint16

type ProductId uint16

type GroupId uint16

const KUniversalGroupID GroupId = 0

type OperationalNodeId uint64

type GroupNodeID uint64

type TemporaryLocalNodeId uint64

type ScopedNodeId struct {
}
