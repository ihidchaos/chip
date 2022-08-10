package raw

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
)

const (
	kMsgUnicastSessionIdUnsecured uint16 = 0x0000

	SessionTypeUnicast uint8 = 0
	SessionTypeGroup   uint8 = 1

	FSourceNodeIdPresent       uint8 = 0b00000100
	FDestinationNodeIdPresent  uint8 = 0b00000001
	FDestinationGroupIdPresent uint8 = 0b00000010
	FDSIZReserved              uint8 = 0b00000011

	FPrivacyFlag      uint8 = 0b10000000
	FControlMsgFlag   uint8 = 0b01000000
	FMsgExtensionFlag uint8 = 0b00100000
	FSessionTypeMask  uint8 = 0b00100000
)

type TSessionType uint8

const (
	KUnicast TSessionType = 0
	KGroup   TSessionType = 1
)

func (T TSessionType) Uint8() uint8 {
	return uint8(T)
}

type PacketHeaderBase interface {
	GetMessageCounter() uint32
	GetSourceNodeId() lib.NodeId
	GetDestinationNodeId() lib.NodeId

	GetDestinationGroupId() lib.GroupId
	GetSessionId() uint16

	GetSessionType() uint8
	GetMessageFlags() uint8
	GetSecurityFlags() uint8
	HasPrivacyFlag() bool

	IsGroupSession() bool
	IsUnicastSession() bool
	IsSessionTypeValid() bool
	IsValidGroupMsg() bool
	IsValidMCSPMsg() bool
	IsEncrypted() bool
	MICTagLength() uint16
	IsSecureSessionControlMsg() bool

	SetSecureSessionControlMsg(value bool)
	SetSourceNodeId(lib.NodeId)
	ClearSourceNodeId()
	SetDestinationNodeId(id lib.NodeId)
	ClearDestinationGroupId()

	SetSessionType(uint82 uint8)
	SetSessionId(uint162 uint16)
	SetMessageCounter(uint32)
	SetUnsecured()
	DecodeAndConsume([]byte) error
}

/**********************************************
 * PacketHeader format (little endian):
 *
 * -------- Unencrypted header -----------------------------------------------------
 *[0:1]   8 bit:  | Message Flags: VERSION: 4 bit | S: 1 bit | RESERVED: 1 bit | DSIZ: 2 bit |
 *[1:2]   8 bit:  | Security Flags: P: 1 bit | C: 1 bit | MX: 1 bit | RESERVED: 3 bit | Session Type: 2 bit |
 *[2:4]   16 bit: | Session ID                                                           |
 *[4:8]   32 bit: | Message Counter                                                      |
 *[8:16]  64 bit: | SOURCE_NODE_ID (iff source node flag is set)                         |
 *[16:24] 64 bit: | DEST_NODE_ID (iff destination node flag is set)                      |
 * -------- Encrypted header -------------------------------------------------------
 *  8 bit:  | Exchange Flags: RESERVED: 3 bit | V: 1 bit | SX: 1 bit | R: 1 bit | A: 1 bit | I: 1 bit |
 *  8 bit:  | Protocol Opcode                                                      |
 * 16 bit:  | Exchange ID                                                          |
 * 16 bit:  | Protocol ID                                                          |
 * 16 bit:  | Optional Vendor ID                                                   |
 * 32 bit:  | Acknowledged Message Counter (if A flag in the PacketHeader is set)        |
 * -------- Encrypted Application Data Start ---------------------------------------
 *  <var>:  | Encrypted Data                                                       |
 * -------- Encrypted Application Data End -----------------------------------------
 *  <var>:  | (Unencrypted) Message Authentication Tag                             |
 *
 **********************************************/

type PacketHeader struct {
	mMessageFlags uint8
	mSecFlags     uint8

	mSourceNodeId       lib.NodeId
	mDestinationNodeId  lib.NodeId
	mDestinationGroupId lib.GroupId

	mSecFlagMask uint8

	mMessageCounter uint32
	mSessionType    TSessionType
	mSessionFlags   uint8

	mSessionId uint16
}

func NewPacketHeader() *PacketHeader {
	return &PacketHeader{
		mMessageFlags:       0,
		mSecFlags:           0,
		mSourceNodeId:       0,
		mDestinationNodeId:  0,
		mDestinationGroupId: 0,
		mSecFlagMask:        0,
		mMessageCounter:     0,
		mSessionType:        0,
		mSessionFlags:       0,
		mSessionId:          0,
	}
}

func (h *PacketHeader) GetMessageCounter() uint32 {
	return h.mMessageCounter
}

func (h *PacketHeader) GetSourceNodeId() lib.NodeId {
	return h.mSourceNodeId
}

func (h *PacketHeader) GetDestinationNodeId() lib.NodeId {
	return h.mDestinationNodeId
}

func (h *PacketHeader) GetDestinationGroupId() lib.GroupId {
	return h.mDestinationGroupId
}

func (h *PacketHeader) GetSessionId() uint16 {
	return h.mSessionId
}

func (h *PacketHeader) GetSessionType() TSessionType {
	return h.mSessionType
}

func (h *PacketHeader) GetMessageFlags() uint8 {
	return h.mMessageFlags
}

func (h *PacketHeader) GetSecurityFlags() uint8 {
	return h.mSessionFlags
}

func (h *PacketHeader) HasPrivacyFlag() bool {
	return lib.HasFlags(h.mSecFlags, FPrivacyFlag)
}

func (h *PacketHeader) IsGroupSession() bool {
	return h.mSessionType.Uint8() == SessionTypeGroup
}

func (h *PacketHeader) IsUnicastSession() bool {
	return h.mSessionType.Uint8() == SessionTypeUnicast
}

func (h *PacketHeader) IsSessionTypeValid() bool {
	switch h.mSessionType {
	case KUnicast:
		return true
	case KGroup:
		return true
	default:
		return false
	}
}

func (h *PacketHeader) IsValidGroupMsg() bool {
	return h.IsGroupSession() && h.GetSourceNodeId() != 0 && h.GetDestinationGroupId() != 0 &&
		!h.IsSecureSessionControlMsg() && h.HasPrivacyFlag()
}

func (h *PacketHeader) IsValidMCSPMsg() bool {
	return h.IsGroupSession() && h.GetSourceNodeId() != 0 && h.GetDestinationNodeId() != 0 &&
		h.IsSecureSessionControlMsg() && h.HasPrivacyFlag()
}

func (h *PacketHeader) IsEncrypted() bool {
	return !(h.mSessionId == kMsgUnicastSessionIdUnsecured && h.IsUnicastSession())
}

func (h *PacketHeader) MICTagLength() uint16 {
	if h.IsEncrypted() {
		return crypto.ChipCryptoAeadMicLengthBytes
	}
	return 0
}

func (h *PacketHeader) IsSecureSessionControlMsg() bool {
	return (h.mSecFlags & FControlMsgFlag) != 0
}

func (h *PacketHeader) SetSecureSessionControlMsg(value bool) {
	//TODO implement me
	panic("implement me")
}

func (h *PacketHeader) SetSourceNodeId(id lib.NodeId) {
	h.mSourceNodeId = id
}

func (h *PacketHeader) ClearSourceNodeId() {
	h.mSourceNodeId = 0
}

func (h *PacketHeader) SetDestinationNodeId(id lib.NodeId) {
	h.mDestinationNodeId = id
}

func (h *PacketHeader) ClearDestinationGroupId() {
	h.mDestinationNodeId = 0
}

func (h *PacketHeader) SetSessionType(t uint8) {
	h.mSessionType = TSessionType(t)
}

func (h *PacketHeader) SetSessionId(id uint16) {
	h.mSessionId = id
}

func (h *PacketHeader) SetMessageCounter(u uint32) {
	h.mMessageCounter = u
}

func (h *PacketHeader) SetUnsecured() {
	h.mSessionId = kMsgUnicastSessionIdUnsecured
	h.mSessionType = KUnicast
}

func (h *PacketHeader) DecodeAndConsume(buf *PacketBuffer) error {

	if buf.DataLength() < 36 {
		return lib.ChipErrorInvalidArgument
	}
	h.mMessageFlags, _ = buf.Read8()
	h.mSecFlags, _ = buf.Read8()
	h.mSessionType = TSessionType(h.mSecFlags & FSessionTypeMask)
	h.mSessionFlags, _ = buf.Read8()
	h.mSessionId, _ = buf.Read16()
	h.mMessageCounter, _ = buf.Read32()

	if lib.HasFlags(h.mMessageFlags, FSourceNodeIdPresent) {
		v, _ := buf.Read64()
		h.mSourceNodeId = lib.NodeId(v)
	}

	if lib.HasFlags(h.mMessageFlags, FDestinationNodeIdPresent) {
		v, _ := buf.Read64()
		h.mDestinationNodeId = lib.NodeId(v)
	}
	if lib.HasFlags(h.mMessageFlags, FDestinationGroupIdPresent) {
		v, _ := buf.Read16()
		h.mDestinationGroupId = lib.GroupId(v)
	} else {
		h.mDestinationNodeId = lib.NodeId(h.GetDestinationNodeId().GetGroupId())
	}
	return nil
}
