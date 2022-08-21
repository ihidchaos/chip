package lib

type DefaultStorageKeyAllocator struct {
	mKeyName       string
	mKeyNameBuffer string
}

func (d *DefaultStorageKeyAllocator) FabricIndexInfo() string {
	return d.SetConst("g/fidx")
}

func (d *DefaultStorageKeyAllocator) SetConst(keyName string) string {
	d.mKeyName = keyName
	return d.mKeyName
}

func (d *DefaultStorageKeyAllocator) FabricNOC(index FabricIndex) string {
	return d.Format("f/%x/n", index)
}

func (d *DefaultStorageKeyAllocator) Format(s string, index FabricIndex) string {
	return ""
}
