package credentials

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/storage"
)

type FabricData struct {
	fabricIndex lib.FabricIndex
	firstGroup  lib.GroupId
	groupCount  uint16
	firstMap    uint16
	mapCount    uint16
	keysetCount uint16
	firstKeyset lib.KeysetId
	next        lib.FabricIndex
}

func (d FabricData) Load(storage storage.KvsPersistentStorageDelegate) {

}

func NewFabricData(index lib.FabricIndex) *FabricData {
	fd := fabricData()
	fd.fabricIndex = index
	return fd
}

func fabricData() *FabricData {
	return &FabricData{
		fabricIndex: lib.KUndefinedFabricIndex,
		firstGroup:  lib.KUndefinedGroupId,
		groupCount:  0,
		firstMap:    0,
		mapCount:    0,
		firstKeyset: lib.KInvalidKeysetId,
		keysetCount: 0,
		next:        lib.KUndefinedFabricIndex,
	}
}
