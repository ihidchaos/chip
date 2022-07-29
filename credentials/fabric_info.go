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
	GetScopedNodeIdForNode(node lib.NodeID) lib.ScopedNodeId

	GetPeerId() device.PeerId
	GetPeerIdForNode(id lib.NodeID) device.PeerId

	GetFabricId() lib.FabricId
	GetFabricIndex() FabricIndex

	GetCompressedFabricId() lib.CompressedFabricId

	GetVendorId() lib.VendorID

	IsInitialized() bool
	HasOperationalKey() bool
}

type FabricInfo struct {
	mFabricLabel         string
	mRootPublicKey       crypto.P256PublicKey
	mNodeId              lib.NodeID
	mFabricId            lib.FabricId
	mFabricIndex         FabricIndex
	mCompressedFabriceId lib.CompressedFabricId
	mVendorId            lib.VendorID
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

func (info *FabricInfo) GetScopedNodeIdForNode(node lib.NodeID) lib.ScopedNodeId {
	//TODO implement me
	panic("implement me")
}

func (info *FabricInfo) GetPeerIdForNode(id lib.NodeID) device.PeerId {
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

func (info *FabricInfo) GetVendorId() lib.VendorID {
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
	NodeId                    lib.NodeID
	FabriceId                 lib.FabricId
	FabricIndex               FabricIndex
	CompressedFabricId        lib.CompressedFabricId
	RootPublicKey             crypto.P256PublicKey
	VendorId                  uint16
	OperationalKeypair        crypto.P256Keypair
	HasExternallyOwnedKeypair bool
}
