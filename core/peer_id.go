package core

import "fmt"

const kUndefinedCompressedFabricId CompressedFabricId = 0

const kUndefinedFabricId FabricId = 0

type PeerId struct {
	mNodeId             NodeId
	mCompressedFabricId CompressedFabricId
}

func (p PeerId) Default() *PeerId {
	return &PeerId{
		mNodeId:             0xFFFF_FFFF_FFFF_1234,
		mCompressedFabricId: 0,
	}
}

func (p PeerId) Init(compressedFabricId CompressedFabricId, nodeId NodeId) *PeerId {
	return &PeerId{
		mNodeId:             nodeId,
		mCompressedFabricId: compressedFabricId,
	}
}

func (p PeerId) SetNodeId(id uint64) {
	p.mNodeId = NodeId(id)
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

func (p PeerId) GetNodeId() NodeId {
	return p.mNodeId
}

func (p PeerId) GetCompressedFabricId() CompressedFabricId {
	return p.mCompressedFabricId
}

func (p PeerId) SetCompressedFabricId() NodeId {
	return p.mNodeId
}

func IsValidFabricId(aFabricId FabricId) bool {
	return aFabricId != kUndefinedFabricId
}
