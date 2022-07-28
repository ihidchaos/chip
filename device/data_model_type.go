package device

import (
	"fmt"
)

type CompressedFabricId uint64
type FabricId uint64
type NodeId uint64
type InstanceName uint64
type FabricIndex uint8

func (id NodeId) String() string {
	var v = uint64(id)
	return fmt.Sprintf("%016X", v)
}

func (name InstanceName) String() string {
	var v = uint(name)
	return fmt.Sprintf("%016X", v)
}
