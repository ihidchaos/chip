package dnssd

import (
	"github.com/galenliu/chip/lib"
)

type PeerId struct {
	mNodeId             lib.NodeId
	mCompressedFabricId lib.CompressedFabricId
}

func NewPeerId(mNodeId lib.NodeId, mCompressedFabricId lib.CompressedFabricId) *PeerId {
	return &PeerId{mNodeId: mNodeId, mCompressedFabricId: mCompressedFabricId}
}

func (p *PeerId) GetCompressedFabricId() lib.CompressedFabricId {
	return p.mCompressedFabricId
}

func (p *PeerId) GetNodeId() lib.NodeId {
	return p.mNodeId
}

func (p *PeerId) Equal(p2 *PeerId) bool {
	return p.mNodeId == p2.mNodeId && p.mCompressedFabricId == p2.mCompressedFabricId
}
