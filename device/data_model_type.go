package device

import (
	"fmt"
)

type InstanceName uint64
type FabricIndex uint8

func (name InstanceName) String() string {
	var v = uint(name)
	return fmt.Sprintf("%016X", v)
}
