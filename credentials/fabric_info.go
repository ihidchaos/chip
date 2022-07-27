package credentials

import "github.com/galenliu/chip/core"

type FabricInfo struct {
}

func (info *FabricInfo) GetPeerId() core.PeerId {
	return core.PeerId{}
}
