package message

import (
	"encoding/binary"
	"github.com/galenliu/chip/internal"
	"github.com/galenliu/chip/lib"
)

const (
	kMsgUnicastSessionIdUnsecured uint16 = 0x0000

	SessionTypeUnicast uint8 = 0
	SessionTypeGroup   uint8 = 1

	//MsgFlagValues
	kSourceNodeIdPresent       = 0b00000100
	kDestinationNodeIdPresent  = 0b00000001
	kDestinationGroupIdPresent = 0b00000010
	kDSIZReserved              = 0b00000011

	//SecFlagValues
	kPrivacyFlag      = 0b10000000
	kControlMsgFlag   = 0b01000000
	kMsgExtensionFlag = 0b00100000
)

type PacketHeader interface {
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
	EncodeSizeBytes()
	Decode(data []byte) []byte
}

type Header struct {
	mMessageFlags uint8
	mSecFlags     uint8

	mSourceNodeId       lib.NodeId
	mDestinationNodeId  lib.NodeId
	mDestinationGroupId lib.GroupId

	mSecFlagMask uint8

	mMessageCounter uint32
	mSessionType    uint8
	mSessionFlags   uint8

	mSessionId uint16
	mLength    uint8
}

func (h Header) GetMessageCounter() uint32 {
	//TODO implement me
	panic("implement me")
}

func (h Header) GetSourceNodeId() lib.NodeId {
	//TODO implement me
	panic("implement me")
}

func (h Header) GetDestinationNodeId() lib.NodeId {
	//TODO implement me
	panic("implement me")
}

func (h Header) GetDestinationGroupId() lib.GroupId {
	//TODO implement me
	panic("implement me")
}

func (h Header) GetSessionId() uint16 {
	//TODO implement me
	panic("implement me")
}

func (h Header) GetSessionType() uint8 {
	return h.mSessionType
}

func (h Header) GetMessageFlags() uint8 {
	return h.mMessageFlags
}

func (h Header) GetSecurityFlags() uint8 {
	return h.mSessionFlags
	//TODO implement me
	panic("implement me")
}

func (h Header) HasPrivacyFlag() bool {
	//TODO implement me
	panic("implement me")
}

func (h Header) IsGroupSession() bool {
	return h.mSessionType == SessionTypeGroup
}

func (h Header) IsUnicastSession() bool {
	return h.mSessionType == SessionTypeUnicast
}

func (h Header) IsSessionTypeValid() bool {
	//TODO implement me
	panic("implement me")
}

func (h Header) IsValidGroupMsg() bool {
	//TODO implement me
	panic("implement me")
}

func (h Header) IsValidMCSPMsg() bool {
	//TODO implement me
	panic("implement me")
}

func (h Header) IsEncrypted() bool {
	return !(h.mSessionId == kMsgUnicastSessionIdUnsecured && h.IsUnicastSession())
}

func (h Header) MICTagLength() uint16 {
	//TODO implement me
	panic("implement me")
}

func (h Header) IsSecureSessionControlMsg() bool {
	//TODO implement me
	panic("implement me")
}

func (h Header) SetSecureSessionControlMsg(value bool) {
	//TODO implement me
	panic("implement me")
}

func (h Header) SetSourceNodeId(id lib.NodeId) {
	//TODO implement me
	panic("implement me")
}

func (h Header) ClearSourceNodeId() {
	//TODO implement me
	panic("implement me")
}

func (h Header) SetDestinationNodeId(id lib.NodeId) {
	//TODO implement me
	panic("implement me")
}

func (h Header) ClearDestinationGroupId() {
	//TODO implement me
	panic("implement me")
}

func (h Header) SetSessionType(uint82 uint8) {
	//TODO implement me
	panic("implement me")
}

func (h Header) SetSessionId(uint162 uint16) {
	//TODO implement me
	panic("implement me")
}

func (h Header) SetMessageCounter(u uint32) {
	//TODO implement me
	panic("implement me")
}

func (h Header) SetUnsecured() {
	//TODO implement me
	panic("implement me")
}

func (h Header) EncodeSizeBytes() {
	//TODO implement me
	panic("implement me")
}

func (h Header) Decode(data []byte) []byte {
	//TODO implement me
	panic("implement me")
}

func (h Header) Len() uint8 {
	return h.mLength
}

/**********************************************
 * Header format (little endian):
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
 * 32 bit:  | Acknowledged Message Counter (if A flag in the Header is set)        |
 * -------- Encrypted Application Data Start ---------------------------------------
 *  <var>:  | Encrypted Data                                                       |
 * -------- Encrypted Application Data End -----------------------------------------
 *  <var>:  | (Unencrypted) Message Authentication Tag                             |
 *
 **********************************************/

func DecodeHeader(data []byte) (*Header, error) {

	if len(data) < 36 {
		return nil, internal.ChipErrorInvalidArgument
	}
	header := &Header{}
	header.mLength = 0
	header.mMessageFlags = data[0]
	header.mLength = header.mLength + 1

	header.mSessionFlags = data[1]
	header.mLength = header.mLength + 1

	header.mSessionId = binary.LittleEndian.Uint16(data[2:4])
	header.mLength = header.mLength + 2

	header.mMessageCounter = binary.LittleEndian.Uint32(data[4:8])
	header.mLength = header.mLength + 4
	if header.mMessageFlags&kSourceNodeIdPresent != 0 {
		header.mSourceNodeId = lib.NodeId(binary.LittleEndian.Uint64(data[8:16]))
		header.mLength = header.mLength + 8
	}
	if header.mMessageFlags&kDestinationNodeIdPresent != 0 {
		header.mDestinationNodeId = lib.NodeId(binary.LittleEndian.Uint64(data[8:16]))
		header.mLength = header.mLength + 8
	}
	return header, nil
}
