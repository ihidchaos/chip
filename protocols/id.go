package protocols

import (
	"fmt"
	"github.com/galenliu/chip/lib"
	"github.com/moznion/go-optional"
	log "golang.org/x/exp/slog"
)

var sUnknownTypeName = "----"

type Id struct {
	VendorId   lib.VendorId
	ProtocolId uint16
}

func New(protocolId uint16, option optional.Option[lib.VendorId]) Id {
	id := Id{
		VendorId:   lib.VidCommon,
		ProtocolId: protocolId,
	}
	if option.IsSome() {
		id.VendorId = option.Unwrap()
	}
	return id
}

func NotSpecifiedId() Id {
	return Id{
		VendorId:   lib.VidNotSpecified,
		ProtocolId: 0xFFFF,
	}
}

//var StandardSecureChannel = &Id{VendorId: lib.VendorIdMatterStandard, ProtocolId: 0x0000}

func FromFullyQualifiedSpecForm(aSpecForm uint32) Id {
	return Id{VendorId: lib.VendorId(aSpecForm >> 16), ProtocolId: uint16(aSpecForm & 0x0000FFFF)}
}

func (id *Id) ToFullyQualifiedSpecForm() uint32 {
	return id.toUint32()
}

func (id *Id) toUint32() uint32 {
	return (uint32(id.VendorId) << 16) | uint32(id.ProtocolId)
}

func (id *Id) Equal(other Id) bool {
	return id.VendorId == other.VendorId && id.ProtocolId == other.ProtocolId
}

func (id *Id) LogValue() log.Value {
	return log.GroupValue(
		log.String("ProtocolId", fmt.Sprintf("%04X", id.ProtocolId)),
		log.String("VendorId", fmt.Sprintf("%04X", id.VendorId)),
	)
}
