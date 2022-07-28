package credentials

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/device"
)

type FabricIndex uint8

type FabricInfoProvider interface {
	GetFabricLabel()
	CommitToStorage()
}

type FabricInfo struct {
	mFabricLabel         string
	mRootPublicKey       crypto.P256PublicKey
	mNodeId              device.NodeId
	mFabricId            device.FabricId
	mFabricIndex         FabricIndex
	mCompressedFabriceId device.CompressedFabricId
	mVendorId            uint16
}

func (info *FabricInfo) GetPeerId() device.PeerId {
	return device.PeerId{}
}

type FabricInfoInitParams struct {
	NodeId                    device.NodeId
	FabriceId                 device.FabricId
	FabricIndex               FabricIndex
	CompressedFabricId        device.CompressedFabricId
	RootPublicKey             crypto.P256PublicKey
	VendorId                  uint16
	OperationalKeypair        crypto.P256Keypair
	HasExternallyOwnedKeypair bool
}
