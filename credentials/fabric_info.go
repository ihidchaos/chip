package credentials

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/pkg/storage"
)

type FabricInfoProvider interface {
	GetFabricLabel() string
	SetFabricLabel(label string)

	GetNodeId() lib.NodeId

	GetScopedNodeId() lib.ScopedNodeId
	GetScopedNodeIdForNode(node lib.NodeId) lib.ScopedNodeId

	FabricId() lib.FabricId
	FabricIndex() lib.FabricIndex

	CompressedFabricId() lib.CompressedFabricId

	GetVendorId() lib.VendorId

	IsInitialized() bool
	HasOperationalKey() bool
}

type FabricInfo struct {
	mFabricLabel         string
	mRootPublicKey       *crypto.P256PublicKey
	mNodeId              lib.NodeId
	mFabricId            lib.FabricId
	mFabricIndex         lib.FabricIndex
	mCompressedFabriceId lib.CompressedFabricId
	mVendorId            lib.VendorId
}

func (info *FabricInfo) GetFabricLabel() string {
	//TODO implement me
	panic("implement me")
}

func (info *FabricInfo) SetFabricLabel(label string) {
	//TODO implement me
	panic("implement me")
}

func (info *FabricInfo) GetScopedNodeId() lib.ScopedNodeId {
	//TODO implement me
	panic("implement me")
}

func (info *FabricInfo) GetScopedNodeIdForNode(node lib.NodeId) lib.ScopedNodeId {
	//TODO implement me
	panic("implement me")
}

func (info *FabricInfo) CommitToStorage(storage storage.KvsPersistentStorageDelegate) {

}

func (info *FabricInfo) FabricId() lib.FabricId {
	return info.mFabricId
}

func (info *FabricInfo) FabricIndex() lib.FabricIndex {
	return info.mFabricIndex
}

func (info *FabricInfo) CompressedFabricId() lib.CompressedFabricId {
	return info.mCompressedFabriceId
}

func (info *FabricInfo) GetVendorId() lib.VendorId {
	return info.mVendorId
}

func (info *FabricInfo) IsInitialized() bool {
	//TODO implement me
	panic("implement me")
}

func (info *FabricInfo) HasOperationalKey() bool {
	//TODO implement me
	panic("implement me")
}

func (info *FabricInfo) GetNodeId() lib.NodeId {
	return info.mNodeId
}

func (info *FabricInfo) Reset() {

}

func (info *FabricInfo) FetchRootPubkey() (*crypto.P256PublicKey, error) {
	if !info.IsInitialized() {
		return nil, lib.KeyNotFound
	}
	return info.mRootPublicKey, nil
}

type FabricInfoInitParams struct {
	NodeId                    lib.NodeId
	FabriceId                 lib.FabricId
	FabricIndex               lib.FabricIndex
	CompressedFabricId        lib.CompressedFabricId
	RootPublicKey             crypto.P256PublicKey
	VendorId                  uint16
	OperationalKeypair        crypto.P256Keypair
	HasExternallyOwnedKeypair bool
}
