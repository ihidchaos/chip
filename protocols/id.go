package protocols

import "github.com/galenliu/chip/lib"

const (
	SecureChannel             = 0x0000
	InteractionModel          = 0x0001
	BDX                       = 0x0002
	UserDirectedCommissioning = 0x0003
	Echo                      = 0x0004
)

type Id struct {
	vendorId   lib.VendorId
	protocolId uint16
}

var SecureChannelId = &Id{
	vendorId:   lib.VendorIdCommon,
	protocolId: 0x0000,
}

var NotSpecifiedProtocolId = &Id{
	vendorId:   lib.VendorIdNotSpecified,
	protocolId: 0xFFFF,
}

func NewProtocolId(vendorId lib.VendorId, protocolId uint16) Id {
	return Id{
		vendorId:   vendorId,
		protocolId: protocolId,
	}
}

//var StandardSecureChannel = &Id{vendorId: lib.VendorIdMatterStandard, protocolId: 0x0000}

func FromFullyQualifiedSpecForm(aSpecForm uint32) Id {
	return Id{vendorId: lib.VendorId(aSpecForm >> 16), protocolId: uint16(aSpecForm&(1<<16) - 1)}
}

func (id *Id) ToFullyQualifiedSpecForm() uint32 {
	return id.toUint32()
}

func (id *Id) VendorId() lib.VendorId { return id.vendorId }

func (id *Id) toUint32() uint32 {
	return (uint32(id.vendorId) << 16) | uint32(id.protocolId)
}

func (id *Id) ProtocolId() uint16 { return id.protocolId }

func (id *Id) Equal(other *Id) bool {
	return id.vendorId == other.vendorId && id.protocolId == other.protocolId
}
