package credentials

import (
	"bytes"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/tlv"
	"github.com/galenliu/chip/pkg/storage"
)

var (
	tagFirstGroup  = tlv.ContextTag(1)
	tagGroupCount  = tlv.ContextTag(2)
	tagFirstMap    = tlv.ContextTag(3)
	tagMapCount    = tlv.ContextTag(4)
	tagFirstKeyset = tlv.ContextTag(5)
	tagKeysetCount = tlv.ContextTag(6)
	tagNext        = tlv.ContextTag(7)
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
	key         string
}

func (f *FabricData) deserialize(tlvReader tlv.Reader) (err error) {

	err = tlvReader.NextE(tlv.AnonymousTag())
	if err != nil {
		return
	}
	container := tlv.Type_Structure
	container, err = tlvReader.EnterContainer()
	if err != nil {
		return
	}

	err = tlvReader.NextE(tagFirstGroup)
	val, err := tlvReader.GetUint16()
	f.firstGroup = lib.GroupId(val)
	if err != nil {
		return err
	}

	err = tlvReader.NextE(tagGroupCount)
	val, err = tlvReader.GetUint16()
	f.groupCount = val
	if err != nil {
		return err
	}

	err = tlvReader.NextE(tagFirstMap)
	val, err = tlvReader.GetUint16()
	f.firstMap = val
	if err != nil {
		return err
	}

	err = tlvReader.NextE(tagMapCount)
	val, err = tlvReader.GetUint16()
	f.mapCount = val
	if err != nil {
		return err
	}

	err = tlvReader.NextE(tagFirstKeyset)
	val, err = tlvReader.GetUint16()
	f.firstKeyset = lib.KeysetId(val)
	if err != nil {
		return err
	}

	err = tlvReader.NextE(tagKeysetCount)
	val, err = tlvReader.GetUint16()
	f.keysetCount = val
	if err != nil {
		return err
	}

	err = tlvReader.NextE(tagNext)
	val8, err := tlvReader.GetUint8()
	f.next = lib.FabricIndex(val8)
	if err != nil {
		return err
	}
	return tlvReader.ExitContainer(container)
}

func (f *FabricData) Load(storage storage.KvsPersistentStorageDelegate) error {
	data, err := storage.GetBytes(f.key)
	if err != nil {
		return err
	}
	tlvReader := tlv.NewReader(bytes.NewBuffer(data))
	return f.deserialize(tlvReader)
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

type GroupKey struct {
	groupId  lib.GroupId
	keysetId lib.KeysetId
}

type LinkedData struct {
	id, index, next, prev, maxId uint16
	first                        bool
}

type KeyMapData struct {
	fabricIndex lib.FabricIndex
	groupId     lib.GroupId
	keysetId    lib.KeysetId
	*GroupKey
	*LinkedData
}

func NewKeyMapData(index lib.FabricIndex, linkId uint16) *KeyMapData {
	return &KeyMapData{
		fabricIndex: index,
		GroupKey: &GroupKey{
			groupId:  0,
			keysetId: 0,
		},
		LinkedData: &LinkedData{
			id:    linkId,
			index: 0,
			next:  0,
			prev:  0,
			maxId: 0,
			first: true,
		},
	}
}

type KeySetData struct {
	fabricIndex     lib.FabricIndex
	next            lib.KeysetId
	prev            lib.KeysetId
	first           bool
	keysetId        lib.KeysetId
	keysetCount     uint8
	policy          any
	operationalKeys []crypto.GroupOperationalCredentials
}

func (d *KeySetData) Find(mStorage storage.KvsPersistentStorageDelegate, fabric *FabricData, targetId lib.KeysetId) bool {
	d.fabricIndex = fabric.fabricIndex
	d.keysetId = fabric.firstKeyset
	d.first = true
	for i := 0; i < int(fabric.keysetCount); i++ {
		err := d.load(mStorage)
		if err != nil {
			continue
		}
		if d.keysetId == targetId {
			return true
		}
		d.first = false
		d.prev = d.keysetId
		d.keysetId = d.next
	}
	return false
}

func (d *KeySetData) load(mStorage storage.KvsPersistentStorageDelegate) interface{} {
	return nil
}
