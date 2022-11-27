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

//type MessageType interface {
//	bdx.MsgType | interaction_model.MsgType | secure_channel.MsgType | echo.MsgType
//}

func NewId(protocolId uint16, option optional.Option[lib.VendorId]) Id {
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

//func (id *Id) ProtocolName() string {
//	if id.mVendorId != lib.VidCommon {
//		return sUnknownTypeName
//	}
//	switch id.ProtocolId() {
//	case secure_channel.ProtocolId:
//		return secure_channel.ProtocolName
//	case interaction_model.ProtocolId:
//		return interaction_model.ProtocolName
//	case bdx.ProtocolId:
//		return bdx.ProtocolName
//	case user_directed_commissioning.ProtocolId:
//		return user_directed_commissioning.ProtocolName
//	case echo.ProtocolId:
//		return echo.ProtocolName
//	default:
//		return sUnknownTypeName
//	}
//}

//func (id *Id) MessageTypeName(messageType uint8) string {
//	if id.mVendorId != lib.VidCommon {
//		return sUnknownTypeName
//	}
//	switch id.ProtocolId() {
//	case secure_channel.ProtocolId:
//		msg := secure_channel.MsgType(messageType)
//		return msg.String()
//	case interaction_model.ProtocolId:
//		msg := interaction_model.MsgType(messageType)
//		return msg.String()
//	case bdx.ProtocolId:
//		msg := bdx.MsgType(messageType)
//		return msg.String()
//	case user_directed_commissioning.ProtocolId:
//		msg := user_directed_commissioning.MsgType(messageType)
//		return msg.String()
//	case echo.ProtocolId:
//		msg := echo.MsgType(messageType)
//		return msg.String()
//	default:
//		return sUnknownTypeName
//	}
//}

func (id *Id) LogValue() log.Value {
	return log.GroupValue(
		log.String("id", fmt.Sprintf("%04X", id.mProtocolId)),
	)
}
