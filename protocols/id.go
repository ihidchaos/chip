package protocols

import (
	"fmt"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/protocols/bdx"
	"github.com/galenliu/chip/protocols/echo"
	"github.com/galenliu/chip/protocols/interaction_model"
	"github.com/galenliu/chip/protocols/secure_channel"
	"github.com/galenliu/chip/protocols/user_directed_commissioning"
	log "golang.org/x/exp/slog"
)

var sUnknownTypeName = "----"

const (
	SecureChannel             uint16 = 0x0000
	InteractionModel          uint16 = 0x0001
	BDX                       uint16 = 0x0002
	UserDirectedCommissioning uint16 = 0x0003
	Echo                      uint16 = 0x0004
)

type Id struct {
	vendorId   lib.VendorId
	protocolId uint16
}

var SecureChannelId = &Id{
	vendorId:   lib.VidCommon,
	protocolId: 0x0000,
}

var NotSpecifiedProtocolId = &Id{
	vendorId:   lib.VidNotSpecified,
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

func (id *Id) ProtocolName() string {
	if id.vendorId != lib.VidCommon {
		return sUnknownTypeName
	}
	switch id.ProtocolId() {
	case SecureChannel:
		return secure_channel.ProtocolName
	case InteractionModel:
		return interaction_model.ProtocolName
	case BDX:
		return bdx.ProtocolName
	case UserDirectedCommissioning:
		return user_directed_commissioning.ProtocolName
	case Echo:
		return echo.ProtocolName
	default:
		return sUnknownTypeName
	}
}

func (id *Id) MessageTypeName(messageType uint8) string {
	if id.vendorId != lib.VidCommon {
		return sUnknownTypeName
	}
	switch id.ProtocolId() {
	case SecureChannel:
		msg := secure_channel.MsgType(messageType)
		return msg.String()
	case InteractionModel:
		msg := interaction_model.MsgType(messageType)
		return msg.String()
	case BDX:
		msg := bdx.MsgType(messageType)
		return msg.String()
	case UserDirectedCommissioning:
		msg := user_directed_commissioning.MsgType(messageType)
		return msg.String()
	case Echo:
		msg := echo.MsgType(messageType)
		return msg.String()
	default:
		return sUnknownTypeName

	}
}

func (id *Id) LogValue() log.Value {
	return log.GroupValue(
		log.String("id", fmt.Sprintf("%04X", id.protocolId)),
		log.String("name", id.ProtocolName()),
	)
}
