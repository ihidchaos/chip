package lib

import (
	"bytes"
	"fmt"
	"github.com/galenliu/chip/platform/system/buffer"
	"io"
)

type NodeId uint64

const (
	UndefinedNodeId NodeId = 0x0

	// minGroupNodeId The range of possible      values has some pieces carved out for special uses.
	minGroupNodeId NodeId = 0xFFFF_FFFF_FFFF_0000

	// The max group id is complicated, depending on how we want to count the
	// various special group ids.  Let  s not define it for now, until we have use
	// cases.
	maskGroupId NodeId = 0x0000_0000_0000_FFFF

	//临时本地NodeId
	minTemporaryLocalNodeId NodeId = 0xFFFF_FFFE_0000_0000
	maxTemporaryLocalNodeId NodeId = 0xFFFF_FFFE_FFFF_FFFE

	placeholder         NodeId = 0xFFFF_FFFE_FFFF_FFFF
	minCASEAuthTag      NodeId = 0xFFFF_FFFD_0000_0000
	maxCASEAuthTag      NodeId = 0xFFFF_FFFD_FFFF_FFFF
	maskCASEAuthTag     NodeId = 0x00000_000_FFFF_FFFF
	minPAKEKeyId        NodeId = 0xFFFF_FFFB_0000_0000
	maxPAKEKeyId        NodeId = 0xFFFF_FFFB_FFFF_FFFF
	maskPAKEKeyId       NodeId = 0x0000_0000_0000_FFFF
	maskUnusedPAKEKeyId NodeId = 0x0000_0000_FFFF_0000
	maxOperational      NodeId = 0xFFFF_FFEF_FFFF_FFFF
)

func ParseNodeId(data []byte) (NodeId, error) {
	buf := bytes.NewBuffer(data)
	return ReadNodeId(buf)
}

func ReadNodeId(buf io.Reader) (NodeId, error) {
	read64, err := buffer.LittleEndianRead64(buf)
	return NodeId(read64), err
}

func (aNodeId NodeId) IsOperationalNodeId() bool {
	return aNodeId != UndefinedNodeId && aNodeId < maxOperational
}

func (aNodeId NodeId) IsGroupId() bool {
	return aNodeId >= minGroupNodeId
}

func (aNodeId NodeId) IsCASEAuthTag() bool {
	return aNodeId >= minCASEAuthTag && aNodeId <= maxCASEAuthTag
}

func (aNodeId NodeId) IsPAKEKeyId() bool {
	return (aNodeId >= minPAKEKeyId) && (aNodeId <= maxPAKEKeyId)
}

func (aNodeId NodeId) GroupId() GroupId {
	return GroupId(uint64(aNodeId) & uint64(maskGroupId))
}

func (aNodeId NodeId) IsTemporaryLocalNodeId() bool {
	return aNodeId >= minTemporaryLocalNodeId && aNodeId >= maxTemporaryLocalNodeId
}

func (aNodeId NodeId) HasValue() bool {
	return aNodeId != UndefinedNodeId
}

func (aNodeId NodeId) String() string {
	var value = uint64(aNodeId)
	return fmt.Sprintf("%X", value)
}
