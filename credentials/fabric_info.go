package credentials

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/device"
	"github.com/galenliu/chip/lib"
)

type FabricIndex uint8

type FabricInfoProvider interface {
	GetFabricLabel() string
	SetFabricLabel(label string)

	GetScopedNodeId() lib.ScopedNodeId
	GetScopedNodeIdForNode(node lib.NodeId) lib.ScopedNodeId

	GetPeerId() device.PeerId
	GetPeerIdForNode(id lib.NodeId) device.PeerId

	GetFabricId() lib.FabricId
	GetFabricIndex() FabricIndex

	GetCompressedFabricId() lib.CompressedFabricId

	GetVendorId() lib.VendorId

	IsInitialized() bool
	HasOperationalKey() bool
}

type FabricInfo struct {
	mFabricLabel         string
	mRootPublicKey       crypto.P256PublicKey
	mNodeId              lib.NodeId
	mFabricId            lib.FabricId
	mFabricIndex         FabricIndex
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

func (info *FabricInfo) GetPeerIdForNode(id lib.NodeId) device.PeerId {
	return device.NewPeerId(id, info.mCompressedFabriceId)
}

func (info *FabricInfo) GetFabricId() lib.FabricId {
	return info.mFabricId
}

func (info *FabricInfo) GetFabricIndex() FabricIndex {
	return info.mFabricIndex
}

func (info *FabricInfo) GetCompressedFabricId() lib.CompressedFabricId {
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

func (info *FabricInfo) GetPeerId() device.PeerId {
	return device.NewPeerId(info.mNodeId, info.mCompressedFabriceId)
}

type FabricInfoInitParams struct {
	NodeId                    lib.NodeId
	FabriceId                 lib.FabricId
	FabricIndex               FabricIndex
	CompressedFabricId        lib.CompressedFabricId
	RootPublicKey             crypto.P256PublicKey
	VendorId                  uint16
	OperationalKeypair        crypto.P256Keypair
	HasExternallyOwnedKeypair bool
}
