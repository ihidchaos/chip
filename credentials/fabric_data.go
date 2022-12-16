package credentials

import (
	"bytes"
	"github.com/galenliu/chip"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/store"
	"github.com/galenliu/chip/lib/tlv"
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
	FabricIndex lib.FabricIndex
	FirstGroup  lib.GroupId
	GroupCount  uint16
	FirstMap    uint16
	MapCount    uint16
	KeysetCount uint16
	FirstKeyset lib.KeysetId
	Next        lib.FabricIndex
	key         string
}

func (f *FabricData) UpdateKey() (lib.StorageKeyName, error) {
	if f.FabricIndex == lib.UndefinedFabricIndex() {
		return "", chip.ErrorInvalidFabricIndex
	}
	key := lib.FabricGroups(f.FabricIndex)
	return key, nil
}

func (f *FabricData) Serialize(e *tlv.Encoder) (err error) {
	var container = tlv.TypeUnknownContainer
	if container, err = e.StartContainer(tlv.AnonymousTag(), tlv.TypeStructure); err != nil {
		return err
	}
	if err = e.Put(tagNext, f.FirstGroup); err != nil {
		return err
	}
	if err = e.Put(tagGroupCount, f.GroupCount); err != nil {
		return err
	}
	if err = e.Put(tagGroupCount, f.GroupCount); err != nil {
		return err
	}
	if err = e.Put(tagFirstMap, f.FirstMap); err != nil {
		return err
	}
	if err = e.Put(tagMapCount, f.MapCount); err != nil {
		return err
	}
	if err = e.Put(tagFirstKeyset, f.FirstKeyset); err != nil {
		return err
	}
	if err = e.Put(tagKeysetCount, f.KeysetCount); err != nil {
		return err
	}
	if err = e.Put(tagNext, f.Next); err != nil {
		return err
	}
	return e.EndContainer(container)
}

func (f *FabricData) Deserialize(d *tlv.Decoder) (err error) {

	if err = d.NextType(tlv.TypeStructure, tlv.AnonymousTag()); err != nil {
		return
	}
	container := tlv.TypeNotSpecified
	if container, err = d.EnterContainer(); err != nil {
		return
	}
	if err = d.NextValue(tagFirstGroup, &f.FirstGroup); err != nil {
		return err
	}

	if err = d.NextValue(tagGroupCount, &f.GroupCount); err != nil {
		return err
	}

	if err = d.NextValue(tagFirstMap, &f.FirstMap); err != nil {
		return err
	}

	if err = d.NextValue(tagMapCount, &f.MapCount); err != nil {
		return err
	}

	if err = d.NextValue(tagFirstKeyset, &f.FirstKeyset); err != nil {
		return err
	}

	if err = d.NextValue(tagKeysetCount, &f.KeysetCount); err != nil {
		return err
	}

	if err = d.NextValue(tagNext, &f.Next); err != nil {
		return err
	}
	return d.ExitContainer(container)
}

func (f *FabricData) Load(storage store.PersistentStorageDelegate) error {
	key, _ := f.UpdateKey()
	var data []byte
	if err := storage.GetKeyValue(key.Name(), data); err != nil {
		return err
	} else {
		tlvReader := tlv.NewDecoder(bytes.NewBuffer(data))
		return f.Deserialize(tlvReader)
	}
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
	LinkedData
}

func (d KeyMapData) Load(storage store.PersistentStorageDelegate) error {
	return nil
}

func NewKeyMapData(index lib.FabricIndex, linkId uint16) *KeyMapData {
	return &KeyMapData{
		fabricIndex: index,
		GroupKey: &GroupKey{
			groupId:  0,
			keysetId: 0,
		},
		LinkedData: LinkedData{
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

func (d *KeySetData) Find(mStorage store.PersistentStorageDelegate, fabric *FabricData, targetId lib.KeysetId) bool {
	d.fabricIndex = fabric.FabricIndex
	d.keysetId = fabric.FirstKeyset
	d.first = true
	for i := 0; i < int(fabric.KeysetCount); i++ {
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

func (d *KeySetData) load(mStorage store.PersistentStorageDelegate) interface{} {
	return nil
}
