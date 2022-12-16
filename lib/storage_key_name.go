package lib

import (
	"fmt"
)

type StorageKeyName string

func Formatted(format string, args ...any) StorageKeyName {
	return StorageKeyName(fmt.Sprintf(format, args...))
}

func (s *StorageKeyName) Name() string {
	return string(*s)
}

func (s *StorageKeyName) IsInitialized() bool {
	return len(string(*s)) != 0
}

func (s *StorageKeyName) IsUninitialized() bool {
	return len(string(*s)) == 0
}

func FabricGroups(index FabricIndex) StorageKeyName {
	return StorageKeyName(fmt.Sprintf("f/%x/g", index))
}

func FabricGroup(index FabricIndex, group GroupId) StorageKeyName {
	return StorageKeyName(fmt.Sprintf("f/%x/g/%x", index, group))
}

func FabricGroupKey(fabric FabricIndex, index uint16) StorageKeyName {
	return StorageKeyName(fmt.Sprintf("f/%x/gk/%x", fabric, index))
}

func FabricGroupEndpoint(fabric FabricIndex, group GroupId, endpoint EndpointId) StorageKeyName {
	return StorageKeyName(fmt.Sprintf("f/%x/g/%x/e/%x", fabric, group, endpoint))
}

func FabricKeyset(fabric FabricIndex, keyset uint16) StorageKeyName {
	return StorageKeyName(fmt.Sprintf("f/%x/k/%x", fabric, keyset))
}

func AttributeValue(endpointId EndpointId, clusterId ClusterId, attributeId AttributeId) StorageKeyName {
	return StorageKeyName(fmt.Sprintf("g/a/%x/%d/%d", endpointId, clusterId, attributeId))
}
