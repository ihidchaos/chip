package protocols

import (
	"fmt"
	"github.com/galenliu/chip/lib"
	"github.com/moznion/go-optional"
	log "golang.org/x/exp/slog"
)

var sUnknownTypeName = "----"

type Id struct {
	mVendorId   lib.VendorId
	mProtocolId uint16
}

func New(protocolId uint16, option optional.Option[lib.VendorId]) Id {
	id := Id{
		mVendorId:   lib.VidCommon,
		mProtocolId: protocolId,
	}
	if option.IsSome() {
		id.mVendorId = option.Unwrap()
	}
	return id
}

//var StandardSecureChannel = &Id{mVendorId: lib.VendorIdMatterStandard, mProtocolId: 0x0000}

func FromFullyQualifiedSpecForm(aSpecForm uint32) Id {
	return Id{mVendorId: lib.VendorId(aSpecForm >> 16), mProtocolId: uint16(aSpecForm & 0x0000FFFF)}
}

func (id *Id) ToFullyQualifiedSpecForm() uint32 {
	return id.toUint32()
}

func (id *Id) VendorId() lib.VendorId { return id.mVendorId }

func (id *Id) toUint32() uint32 {
	return (uint32(id.mVendorId) << 16) | uint32(id.mProtocolId)
}

func (id *Id) ProtocolId() uint16 { return id.mProtocolId }

func (id *Id) Equal(other Id) bool {
	return id.mVendorId == other.mVendorId && id.mProtocolId == other.mProtocolId
}

func (id *Id) LogValue() log.Value {
	return log.GroupValue(
		log.String("id", fmt.Sprintf("%04X", id.mProtocolId)),
	)
}
