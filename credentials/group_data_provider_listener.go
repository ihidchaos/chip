package credentials

import (
	"github.com/galenliu/chip/lib"
)

type GroupListener interface {
	OnGroupAdded(fabricIndex lib.FabricIndex, newGroup *GroupInfo)
	OnGroupRemoved(fabricIndex lib.FabricIndex, newGroup *GroupInfo)
}
