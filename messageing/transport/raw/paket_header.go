package raw

import (
	"bytes"
	"encoding/binary"
	"github.com/galenliu/chip"
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/bitflags"
	"github.com/moznion/go-optional"
	"io"
)

const (
	fSourceNodeIdPresent       uint8 = 0b00000100
	fDestinationNodeIdPresent  uint8 = 0b00000001
	fDestinationGroupIdPresent uint8 = 0b00000010
	//FDSIZReserved              uint8 = 0b00000011

	fPrivacy      uint8 = 0b10000000
	fControlMsg   uint8 = 0b01000000
	fMsgExtension uint8 = 0b00100000

	fSessionTypeMask uint8 = 0b00000011
	fVersionIdMask   uint8 = 0b11110000
	fVersionShift          = 4
	fMsgFlagsMask          = 0x07
)

type SessionType uint8

const (
	UnicastSession SessionType = 0
	GroupSession   SessionType = 1
)

func (t SessionType) String() string {
	var value = uint8(t)
	switch value {
	case 0:
		return "unicase"
	case 1:
		return "GroupSession"
	default:
		return "unknown"
	}
}

/**********************************************
 * PacketHeader format (little endian):
 *
 * -------- Unencrypted header -----------------------------------------------------
 *[0:1]   8 bit:  | Message Flags: VERSION: 4 bit | S: 1 bit | RESERVED: 1 bit | DSIZ: 2 bit |
 *[1:2]   8 bit:  | Security Flags: P: 1 bit | C: 1 bit | MX: 1 bit | RESERVED: 3 bit | Session TransportType: 2 bit |
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
 *  <var>:  | (Unencrypted) Message Authentication NextTag                             |
 *
 **********************************************/

type PacketHeader struct {
	DestinationGroupId optional.Option[lib.GroupId]

	mMessageFlags     bitflags.Flags[uint8]
	mSecFlags         bitflags.Flags[uint8]
	SessionId         uint16
	MessageCounter    uint32
	SourceNodeId      optional.Option[lib.NodeId]
	DestinationNodeId optional.Option[lib.NodeId]

	mSessionType SessionType
}

type paketHeaderOption func(*PacketHeader)

func NewPacketHeader(opts ...paketHeaderOption) *PacketHeader {
	h := &PacketHeader{
		mMessageFlags:      bitflags.Some[uint8](0),
		mSecFlags:          bitflags.Some[uint8](0),
		SourceNodeId:       nil,
		DestinationNodeId:  nil,
		DestinationGroupId: nil,
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
	return header.mSecFlags.Has(fPrivacy)
}

func (header *PacketHeader) IsGroupSession() bool {
	return header.mSessionType == GroupSession
}

func (header *PacketHeader) IsUnicastSession() bool {
	return header.mSessionType == UnicastSession
}

func (header *PacketHeader) SecurityFlags() uint8 {
	return header.mSecFlags.Raw()
}

func (header *PacketHeader) IsSessionTypeValid() bool {
	switch header.mSessionType {
	case UnicastSession:
		return true
	case GroupSession:
		return true
	default:
		return false
	}
}

func (header *PacketHeader) IsValidGroupMsg() bool {
	return header.IsGroupSession() && header.SourceNodeId.IsSome() && header.DestinationGroupId.IsSome() &&
		!header.IsSecureSessionControlMsg() && header.HasPrivacyFlag()
}

func (header *PacketHeader) IsValidMCSPMsg() bool {
	return header.IsGroupSession() && header.SourceNodeId.IsSome() && header.DestinationNodeId.IsSome() &&
		header.IsSecureSessionControlMsg() && header.HasPrivacyFlag()
}

func (header *PacketHeader) IsEncrypted() bool {
	return !(header.SessionId == kMsgUnsecuredUnicastSessionId && header.IsUnicastSession())
}

func (header *PacketHeader) VersionId() uint8 {
	return (header.mMessageFlags.Raw() & fVersionIdMask) >> 4
}

func (header *PacketHeader) MICTagLength() int {
	if header.IsEncrypted() {
		return crypto.AEADMicLengthBytes
	}
	return 0
}

func (header *PacketHeader) IsSecureSessionControlMsg() bool {
	return header.mSecFlags.Has(fControlMsg)
}

func (header *PacketHeader) SetSecureSessionControlMsg(b bool) *PacketHeader {
	header.mSecFlags.Set(b, fControlMsg)
	return header
}

func (header *PacketHeader) SetSourceNodeId(id lib.NodeId) *PacketHeader {
	header.SourceNodeId = optional.Some(id)
	header.mMessageFlags.Sets(fSourceNodeIdPresent)
	return header
}

func (header *PacketHeader) ClearSourceNodeId() {
	header.SourceNodeId = nil
}

func (header *PacketHeader) SetDestinationNodeId(id lib.NodeId) {
	header.DestinationNodeId = optional.Some(id)
}

func (header *PacketHeader) ClearDestinationGroupId() {
	header.DestinationNodeId = nil
}

func (header *PacketHeader) SetSessionId(id uint16) {
	header.SessionId = id
}

func (header *PacketHeader) SetMessageCounter(u uint32) {
	header.MessageCounter = u
}

func (header *PacketHeader) SetUnsecured() {
	header.SessionId = kMsgUnsecuredUnicastSessionId
	header.mSessionType = UnicastSession
}

func (header *PacketHeader) Encode(w io.Writer) (err error) {

	messageFlags := header.mMessageFlags
	messageFlags.Set(header.SourceNodeId.IsSome(), fSourceNodeIdPresent)
	messageFlags.Set(header.DestinationNodeId.IsSome(), fDestinationNodeIdPresent)
	messageFlags.Set(header.DestinationGroupId.IsSome(), fDestinationGroupIdPresent)

	msgFlags := (kMsgHeaderVersion << fVersionShift) | (messageFlags.Raw() & fMsgFlagsMask)
	buf := new(bytes.Buffer)
	if err = buf.WriteByte(msgFlags); err != nil {
		return
	}
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, header.SessionId)
	if _, err = buf.Write(data); err != nil {
		return
	}
	if err = buf.WriteByte(header.mSecFlags.Raw()); err != nil {
		return
	}
	data = make([]byte, 4)
	binary.LittleEndian.PutUint32(data, header.MessageCounter)
	if _, err = buf.Write(data); err != nil {
		return
	}

	if header.SourceNodeId.IsSome() {
		data = make([]byte, 8)
		binary.LittleEndian.PutUint64(data, uint64(header.SourceNodeId.Unwrap()))
		if _, err = buf.Write(data); err != nil {
			return
		}
	}
	if header.DestinationNodeId.IsSome() {
		data = make([]byte, 8)
		binary.LittleEndian.PutUint64(data, uint64(header.DestinationNodeId.Unwrap()))
		if _, err = buf.Write(data); err != nil {
			return
		}
	} else if header.DestinationGroupId.IsSome() {
		data = make([]byte, 2)
		binary.LittleEndian.PutUint16(data, uint16(header.DestinationGroupId.Unwrap()))
		if _, err = buf.Write(data); err != nil {
			return
		}
	}
	if buf.Len() != int(header.EncodeSizeBytes()) {
		return chip.MATTER_ERROR_INVALID_MESSAGE_LENGTH
	}
	_, err = buf.WriteTo(w)
	return err
}

func (header *PacketHeader) DecodeAndConsume(buf io.Reader) error {

	data := make([]byte, 1)
	_, err := buf.Read(data)
	if err != nil {
		return err
	}
	header.mMessageFlags = bitflags.Some(data[0])

	_, err = buf.Read(data)
	if err != nil {
		return err
	}
	header.setSecurityFlags(data[0])

	data = make([]byte, 2)
	_, err = buf.Read(data)
	header.SessionId = binary.LittleEndian.Uint16(data)
	if err != nil {
		return err
	}

	data = make([]byte, 4)
	_, err = buf.Read(data)
	header.MessageCounter = binary.LittleEndian.Uint32(data)
	if err != nil {
		return err
	}

	if header.mMessageFlags.Has(fSourceNodeIdPresent) {
		data = make([]byte, 8)
		_, err = buf.Read(data)
		if err != nil {
			return err
		}
		v := binary.LittleEndian.Uint64(data)
		header.SourceNodeId = optional.Some(lib.NodeId(v))
	}
	if header.mMessageFlags.Has(fDestinationNodeIdPresent) {
		data = make([]byte, 8)
		_, err = buf.Read(data)
		if err != nil {
			return err
		}
		v := binary.LittleEndian.Uint64(data)
		header.DestinationNodeId = optional.Some(lib.NodeId(v))
	}
	if header.mMessageFlags.Has(fDestinationGroupIdPresent) {
		data = make([]byte, 2)
		_, err = buf.Read(data)
		if err != nil {
			return err
		}
		v := binary.LittleEndian.Uint16(data)
		header.DestinationGroupId = optional.Some(lib.GroupId(v))
	}
	return nil
}

func (header *PacketHeader) EncodeSizeBytes() uint16 {
	size := kFixedUnencryptedHeaderSizeBytes
	if header.SourceNodeId.IsSome() {
		size += kNodeIdSizeBytes
	}
	if header.DestinationNodeId.IsSome() {
		size += kNodeIdSizeBytes
	} else if header.DestinationGroupId.IsSome() {
		size += kGroupIdSizeBytes
	}
	return uint16(size)
}

func (header *PacketHeader) setSecurityFlags(flags uint8) {
	header.mSecFlags = bitflags.Some(flags)
	header.mSessionType = SessionType(flags & fSessionTypeMask)
}

func (header *PacketHeader) setMessageFlags(flags uint8) {
	header.mMessageFlags = bitflags.Some(flags)

}

func (header *PacketHeader) SetDestinationGroupId(id lib.GroupId) *PacketHeader {
	header.DestinationGroupId = optional.Some(id)
	header.mMessageFlags.Sets(fDestinationGroupIdPresent)
	return header
}

func (header *PacketHeader) SetSessionType(sessionType SessionType) *PacketHeader {
	header.mSessionType = sessionType
	header.mSecFlags.Sets(fSessionTypeMask)
	return header
}
