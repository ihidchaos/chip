package lib

const (
	UndefinedFabricIndex FabricIndex = 0
)

const TestVendorId1 = 0xFFF1 //Test Vendor #1
const TestVendorId2 = 0xFFF2 //Test Vendor #2
const TestVendorId3 = 0xFFF3 //Test Vendor #3
const TestVendorId4 = 0xFFF4 //Test Vendor #4

type CompressedFabricId uint64
type FabricId uint64
type FabricIndex uint8

type VendorId uint16

type ProductId uint16

type OperationalNodeId uint64

type GroupNodeID uint64

type TemporaryLocalNodeId uint64

type GroupId uint16

// 0xFF00-0xFFFC Reserved for future use
const (
	UniversalGroupID  GroupId = 0
	AllNode           GroupId = 0xFFFF
	AllNonSleepyNodes GroupId = 0xFFFE
	AllProxies        GroupId = 0xFFFD
)
