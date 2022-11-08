package raw

import (
	"bytes"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	buffer2 "github.com/galenliu/chip/pkg/tlv/buffer"
)

const (
	kMsgUnicastSessionIdUnsecured uint16 = 0x0000
	KMsgHeaderVersion             uint8  = 0x0000

	FSourceNodeIdPresent       uint8 = 0b00000100
	FDestinationNodeIdPresent  uint8 = 0b00000001
	FDestinationGroupIdPresent uint8 = 0b00000010
	//FDSIZReserved              uint8 = 0b00000011

	FPrivacyFlag    uint8 = 0b10000000
	FControlMsgFlag uint8 = 0b01000000

	//FMsgExtensionFlag uint8 = 0b00100000

	FSessionTypeMask uint8 = 0b00000001
	FVersionIdMask   uint8 = 0b11110000
)

type TSessionType uint8

const (
	KUnicast TSessionType = 0
	KGroup   TSessionType = 1
)

func (T TSessionType) Uint8() uint8 {
	return uint8(T)
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
	MessageFlags       uint8
	SessionId          uint16
	SecFlags           uint8
	MessageCounter     uint32
	SourceNodeId       lib.NodeId
	DestinationNodeId  lib.NodeId
	DestinationGroupId lib.GroupId
	SessionType        TSessionType
}

func NewPacketHeader() *PacketHeader {
	return &PacketHeader{
		MessageFlags:       0,
		SecFlags:           0,
		SourceNodeId:       lib.UndefinedNodeId,
		DestinationNodeId:  lib.UndefinedNodeId,
		DestinationGroupId: lib.UndefinedGroupId,
		MessageCounter:     0,
		SessionType:        0,
		SessionId:          0,
	}
}

func (header *PacketHeader) HasPrivacyFlag() bool {
	return lib.HasFlags(header.SecFlags, FPrivacyFlag)
}

func (header *PacketHeader) IsGroupSession() bool {
	return lib.HasFlags(header.SecFlags, FSessionTypeMask)
}

func (header *PacketHeader) IsUnicastSession() bool {
	return !lib.HasFlags(header.SecFlags, FSessionTypeMask)
}

func (header *PacketHeader) IsSessionTypeValid() bool {
	switch header.SessionType {
	case KUnicast:
		return true
	case KGroup:
		return true
	default:
		return false
	}
}

func (header *PacketHeader) IsValidGroupMsg() bool {
	return header.IsGroupSession() && header.SourceNodeId != 0 && header.DestinationGroupId != 0 &&
		!header.IsSecureSessionControlMsg() && header.HasPrivacyFlag()
}

func (header *PacketHeader) IsValidMCSPMsg() bool {
	return header.IsGroupSession() && header.SourceNodeId != 0 && header.DestinationNodeId != 0 &&
		header.IsSecureSessionControlMsg() && header.HasPrivacyFlag()
}

func (header *PacketHeader) IsEncrypted() bool {
	return !(header.SessionId == kMsgUnicastSessionIdUnsecured && header.IsUnicastSession())
}

func (header *PacketHeader) GetVersionId() uint8 {
	return (header.MessageFlags & FVersionIdMask) >> 4
}

func (header *PacketHeader) MICTagLength() uint16 {
	if header.IsEncrypted() {
		return crypto.MatterCryptoAEADMicLengthBytes
	}
	return 0
}

func (header *PacketHeader) IsSecureSessionControlMsg() bool {
	return (header.SecFlags & FControlMsgFlag) != 0
}

func (header *PacketHeader) SetSecureSessionControlMsg(value bool) {
	//TODO implement me
	panic("implement me")
}

func (header *PacketHeader) SetSourceNodeId(id lib.NodeId) {
	header.SourceNodeId = id
}

func (header *PacketHeader) ClearSourceNodeId() {
	header.SourceNodeId = 0
}

func (header *PacketHeader) SetDestinationNodeId(id lib.NodeId) {
	header.DestinationNodeId = id
}

func (header *PacketHeader) ClearDestinationGroupId() {
	header.DestinationNodeId = 0
}

func (header *PacketHeader) SetSessionType(t uint8) {
	header.SessionType = TSessionType(t)
}

func (header *PacketHeader) SetSessionId(id uint16) {
	header.SessionId = id
}

func (header *PacketHeader) SetMessageCounter(u uint32) {
	header.MessageCounter = u
}

func (header *PacketHeader) SetUnsecured() {
	header.SessionId = kMsgUnicastSessionIdUnsecured
	header.SessionType = KUnicast
}

func (header *PacketHeader) Encode() (*bytes.Buffer, error) {

	var msgFlags = header.MessageFlags
	msgFlags = lib.SetFlag(header.SourceNodeId.HasValue(), msgFlags, FSourceNodeIdPresent)
	msgFlags = lib.SetFlag(header.DestinationNodeId.HasValue(), msgFlags, FDestinationNodeIdPresent)
	msgFlags = lib.SetFlag(header.DestinationGroupId.HasValue(), msgFlags, FDestinationGroupIdPresent)
	msgFlags = KMsgHeaderVersion<<4 | msgFlags
	buf := bytes.NewBuffer(nil)
	err := buffer2.Write8(buf, msgFlags)
	err = buffer2.LittleEndianWrite16(buf, header.SessionId)
	err = buffer2.Write8(buf, header.SecFlags)
	err = buffer2.LittleEndianWrite32(buf, header.MessageCounter)
	if header.SourceNodeId.HasValue() {
		err = buffer2.LittleEndianWrite64(buf, uint64(header.SourceNodeId))
	}
	if header.DestinationNodeId.HasValue() {
		err = buffer2.LittleEndianWrite64(buf, uint64(header.DestinationNodeId))
	} else if header.DestinationGroupId.HasValue() {
		err = buffer2.LittleEndianWrite16(buf, uint16(header.DestinationGroupId))
	}
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (header *PacketHeader) GetMessageCounter() uint32 {
	return header.MessageCounter
}
