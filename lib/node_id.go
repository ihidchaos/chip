package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type NodeId uint64

const (
	kUndefinedNodeId NodeId = 0x0000_0000_0000_0000

	// minGroupNodeId The range of possible      values has some pieces carved out for special uses.
	kMinGroupNodeId NodeId = 0xFFFF_FFFF_FFFF_0000

	// The max group id is complicated, depending on how we want to count the
	// various special group ids.  Let  s not define it for now, until we have use
	// cases.
	kMaskGroupId NodeId = 0x0000_0000_0000_FFFF

	//临时本地NodeId
	kMinTemporaryLocalNodeId NodeId = 0xFFFF_FFFE_0000_0000
	kMaxTemporaryLocalNodeId NodeId = 0xFFFF_FFFE_FFFF_FFFE
	kPlaceholderNodeId       NodeId = 0xFFFF_FFFE_FFFF_FFFF

	kMinCASEAuthTag NodeId = 0xFFFF_FFFD_0000_0000
	kMaxCASEAuthTag NodeId = 0xFFFF_FFFD_FFFF_FFFF
	maskCASEAuthTag NodeId = 0x00000_000_FFFF_FFFF

	kMinPAKEKeyId        NodeId = 0xFFFF_FFFB_0000_0000
	kMaxPAKEKeyId        NodeId = 0xFFFF_FFFB_FFFF_FFFF
	kMaskPAKEKeyId       NodeId = 0x0000_0000_0000_FFFF
	kMaskUnusedPAKEKeyId NodeId = 0xFFFF_FFEF_FFFF_FFFF

	kMaxOperationalNodeId NodeId = 0xFFFF_FFEF_FFFF_FFFF
)

func (aNodeId NodeId) IsOperationalNodeId() bool {
	return aNodeId != kUndefinedNodeId && aNodeId < kMaxOperationalNodeId
}

func ParseNodeId(data []byte) (NodeId, error) {
	buf := bytes.NewBuffer(data)
	return ReadNodeId(buf)
}

func (aNodeId NodeId) IsGroupId() bool {
	return aNodeId >= kMinGroupNodeId
}

func ReadNodeId(buf io.Reader) (NodeId, error) {
	data := make([]byte, 8)
	_, err := buf.Read(data)
	v := binary.LittleEndian.Uint64(data)
	return NodeId(v), err
}

func (aNodeId NodeId) IsCASEAuthTag() bool {
	return aNodeId >= kMinCASEAuthTag && aNodeId <= kMaxCASEAuthTag
}

func (aNodeId NodeId) IsPAKEKeyId() bool {
	return (aNodeId >= kMinPAKEKeyId) && (aNodeId <= kMaxPAKEKeyId)
}

func (aNodeId NodeId) GroupId() GroupId {
	return GroupId(uint64(aNodeId) & uint64(kMaskGroupId))
}

func (aNodeId NodeId) IsTemporaryLocalNodeId() bool {
	return aNodeId >= kMinTemporaryLocalNodeId && aNodeId <= kMaxTemporaryLocalNodeId
}

func (aNodeId NodeId) PasscodeId() PasscodeId {
	return PasscodeId(aNodeId & kMaskPAKEKeyId)
}

func (aNodeId NodeId) String() string {
	var value = uint64(aNodeId)
	return fmt.Sprintf("%016X", value)
}

func UndefinedNodeId() NodeId {
	return kUndefinedNodeId
}
