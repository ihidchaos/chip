package raw

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/protocols"
)

type TExchangeFlag = uint8

const (
	FInitiator TExchangeFlag = 0b1
	// FAckMsg / Set when current message is an acknowledgment for a previously received message.
	FAckMsg TExchangeFlag = 0b10
	// FNeedsAck / Set when current message is requesting an acknowledgment from the recipient.
	FNeedsAck TExchangeFlag = 0b100
	// FSecuredExtension / Secured Extension block is present.
	FSecuredExtension TExchangeFlag = 0b1000
	// FVendorIdPresent / Set when a vendor id is prepended to the Message Protocol Id field.
	FVendorIdPresent TExchangeFlag = 0b10000
)

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
 *  8 bit:  | Protocol Opcode   /Sigma1/Sigma2//Sigma4/Sigma1 Fin                 |
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

type PayloadHeader struct {
	mExchangeFlags     uint8
	mExchangeId        uint16
	mProtocolId        protocols.Id
	mVendorId          uint16
	mAckMessageCounter uint32
	mSize              uint8
	mProtocolOpcode    uint8
	mMessageType       uint8
}

func NewPayloadHeader() *PayloadHeader {
	header := &PayloadHeader{}
	return header
}

func (header *PayloadHeader) IsInitiator() bool {
	return lib.HasFlags(header.mExchangeFlags, FInitiator)
}

func (header *PayloadHeader) HasMessageType(t uint8) bool {
	return header.mMessageType == t
}

func (header *PayloadHeader) IsAckMsg() bool {
	return lib.HasFlags(header.mExchangeFlags, FAckMsg)
}

func (header *PayloadHeader) NeedsAck() bool {
	return lib.HasFlags(header.mExchangeFlags, FNeedsAck)
}

func (header *PayloadHeader) HaveVendorId() bool {
	return lib.HasFlags(header.mExchangeFlags, FVendorIdPresent)
}

func (header *PayloadHeader) DecodeAndConsume(data *PacketBuffer) error {
	return header.Decode(data)
}

func (header *PayloadHeader) Decode(buf *PacketBuffer) error {
	header.mExchangeFlags, _ = buf.Read8()
	header.mProtocolOpcode, _ = buf.Read8()
	header.mExchangeId, _ = buf.Read16()
	protocolId, _ := buf.Read16()
	var vendorId = lib.KVidCommon
	if header.HaveVendorId() {
		vendorId, _ = buf.Read16()
	}
	header.mProtocolId = protocols.Id{
		VendorId:   vendorId,
		ProtocolId: protocolId,
	}
	if header.IsAckMsg() {
		header.mVendorId, _ = buf.Read16()
	}
	return nil
}

func (header *PayloadHeader) GetProtocolID() protocols.Id {
	return header.mProtocolId
}

func (header *PayloadHeader) GetMessageType() uint8 {
	return header.mMessageType
}

func (header *PayloadHeader) GetExchangeID() uint16 {
	return header.mExchangeId
}

func (header *PayloadHeader) HasProtocol(id *protocols.Id) bool {
	return header.mProtocolId.Equal(id)
}
