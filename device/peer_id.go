package device

import (
	"fmt"
	"github.com/galenliu/chip/lib"
)

const kUndefinedCompressedFabricId lib.CompressedFabricId = 0

const kUndefinedFabricId lib.FabricId = 0

type PeerId struct {
	mNodeId             lib.NodeID
	mCompressedFabricId lib.CompressedFabricId
}

func NewPeerId(mNodeId lib.NodeID, mCompressedFabricId lib.CompressedFabricId) PeerId {
	return PeerId{mNodeId: mNodeId, mCompressedFabricId: mCompressedFabricId}
}

func (p PeerId) Default() *PeerId {
	return &PeerId{
		mNodeId:             0xFFFF_FFFF_FFFF_1234,
		mCompressedFabricId: 0,
	}
}

func (p PeerId) Init(compressedFabricId lib.CompressedFabricId, nodeId lib.NodeID) *PeerId {
	return &PeerId{
		mNodeId:             nodeId,
		mCompressedFabricId: compressedFabricId,
	}
}

func (p PeerId) SetNodeId(id uint64) {
	p.mNodeId = lib.NodeID(id)
}

func (p PeerId) String() string {
	nodeId := p.GetNodeId()               //uint64
	fabricId := p.GetCompressedFabricId() //uint64
	fabricIdH32 := uint32(fabricId >> 32)
	fabricIdL32 := uint32(fabricId)
	nodeIdH32 := uint32(nodeId >> 32)
	nodeIdL32 := uint32(nodeId)
	return fmt.Sprintf("%08x%08x%08x%08x", fabricIdH32, fabricIdL32, nodeIdH32, nodeIdL32)
}

func (p PeerId) GetNodeId() lib.NodeID {
	return p.mNodeId
}

func (p PeerId) GetCompressedFabricId() lib.CompressedFabricId {
	return p.mCompressedFabricId
}

func (p PeerId) SetCompressedFabricId() lib.NodeID {
	return p.mNodeId
}

func IsValidFabricId(aFabricId lib.FabricId) bool {
	return aFabricId != kUndefinedFabricId
}
