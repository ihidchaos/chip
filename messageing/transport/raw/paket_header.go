package raw

import (
	"bytes"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	buffer "github.com/galenliu/chip/platform/system/buffer"
	"io"
)

const (
	kMsgUnsecuredUnicastSessionId uint16 = 0x0000
	KMsgHeaderVersion             uint8  = 0x0000

	fSourceNodeIdPresent       uint8 = 0b00000100
	fDestinationNodeIdPresent  uint8 = 0b00000001
	fDestinationGroupIdPresent uint8 = 0b00000010
	//FDSIZReserved              uint8 = 0b00000011

	FPrivacyFlag    uint8 = 0b10000000
	FControlMsgFlag uint8 = 0b01000000

	//FMsgExtensionFlag uint8 = 0b00100000

	fSessionTypeMask uint8 = 0b00000011
	FVersionIdMask   uint8 = 0b11110000
)

type SessionType uint8

const (
	unicast SessionType = 0
	group   SessionType = 1
)

func (t SessionType) String() string {
	var value = uint8(t)
	switch value {
	case 0:
		return "unicase"
	case 1:
		return "group"
	default:
		return "unknown"
	}
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
	DestinationGroupId lib.GroupId

	mMessageFlags     uint8
	mSecFlags         uint8
	SessionId         uint16
	MessageCounter    uint32
	SourceNodeId      lib.NodeId
	DestinationNodeId lib.NodeId

	mSessionType SessionType
}

type paketHeaderOption func(*PacketHeader)

func NewPacketHeader(opts ...paketHeaderOption) *PacketHeader {
	h := &PacketHeader{
		mMessageFlags:      0,
		mSecFlags:          0,
		SourceNodeId:       lib.UndefinedNodeId,
		DestinationNodeId:  lib.UndefinedNodeId,
		DestinationGroupId: lib.UndefinedGroupId,
		MessageCounter:     0,
		mSessionType:       0,
		SessionId:          0,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (header *PacketHeader) HasPrivacyFlag() bool {
	return lib.HasFlags(header.mSecFlags, FPrivacyFlag)
}

func (header *PacketHeader) IsGroupSession() bool {
	return header.mSessionType == group
}

func (header *PacketHeader) IsUnicastSession() bool {
	return header.mSessionType == unicast
}

func (header *PacketHeader) SecurityFlags() uint8 {
	return header.mSecFlags
}

func (header *PacketHeader) IsSessionTypeValid() bool {
	switch header.mSessionType {
	case unicast:
		return true
	case group:
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
	return !(header.SessionId == kMsgUnsecuredUnicastSessionId && header.IsUnicastSession())
}

func (header *PacketHeader) VersionId() uint8 {
	return (header.mMessageFlags & FVersionIdMask) >> 4
}

func (header *PacketHeader) MICTagLength() uint16 {
	if header.IsEncrypted() {
		return crypto.MatterCryptoAEADMicLengthBytes
	}
	return 0
}

func (header *PacketHeader) IsSecureSessionControlMsg() bool {
	return (header.mSecFlags & FControlMsgFlag) != 0
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

func (header *PacketHeader) SetSessionId(id uint16) {
	header.SessionId = id
}

func (header *PacketHeader) SetMessageCounter(u uint32) {
	header.MessageCounter = u
}

func (header *PacketHeader) SetUnsecured() {
	header.SessionId = kMsgUnsecuredUnicastSessionId
	header.mSessionType = unicast
}

func (header *PacketHeader) Encode() (*bytes.Buffer, error) {

	var msgFlags = header.mMessageFlags
	msgFlags = lib.SetFlag(header.SourceNodeId.HasValue(), msgFlags, fSourceNodeIdPresent)
	msgFlags = lib.SetFlag(header.DestinationNodeId.HasValue(), msgFlags, fDestinationNodeIdPresent)
	msgFlags = lib.SetFlag(header.DestinationGroupId.HasValue(), msgFlags, fDestinationGroupIdPresent)
	msgFlags = KMsgHeaderVersion<<4 | msgFlags
	buf := bytes.NewBuffer(nil)
	err := buffer.Write8(buf, msgFlags)
	err = buffer.LittleEndianWrite16(buf, header.SessionId)
	err = buffer.Write8(buf, header.mSecFlags)
	err = buffer.LittleEndianWrite32(buf, header.MessageCounter)
	if header.SourceNodeId.HasValue() {
		err = buffer.LittleEndianWrite64(buf, uint64(header.SourceNodeId))
	}
	if header.DestinationNodeId.HasValue() {
		err = buffer.LittleEndianWrite64(buf, uint64(header.DestinationNodeId))
	} else if header.DestinationGroupId.HasValue() {
		err = buffer.LittleEndianWrite16(buf, uint16(header.DestinationGroupId))
	}
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (header *PacketHeader) DecodeAndConsume(buf io.Reader) error {

	var err error
	header.mMessageFlags, err = buffer.Read8(buf)

	secFlags, err := buffer.Read8(buf)
	header.setSecurityFlags(secFlags)

	header.SessionId, err = buffer.LittleEndianRead16(buf)
	header.MessageCounter, err = buffer.LittleEndianRead32(buf)
	if err != nil {
		return err
	}

	if lib.HasFlags(header.mMessageFlags, fSourceNodeIdPresent) {
		v, _ := buffer.LittleEndianRead64(buf)
		header.SourceNodeId = lib.NodeId(v)
	}
	if lib.HasFlags(header.mMessageFlags, fDestinationNodeIdPresent) {
		v, _ := buffer.LittleEndianRead64(buf)
		header.DestinationNodeId = lib.NodeId(v)
	}
	if lib.HasFlags(header.mMessageFlags, fDestinationGroupIdPresent) {
		v, _ := buffer.LittleEndianRead16(buf)
		header.DestinationGroupId = lib.GroupId(v)
	}
	return nil
}

func (header *PacketHeader) setSecurityFlags(flags uint8) {
	header.mSecFlags = flags
	header.mSessionType = SessionType(flags & fSessionTypeMask)
}
