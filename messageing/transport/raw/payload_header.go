package raw

import "encoding/binary"

const (
	//ExFlagValues
	kExchangeFlagInitiator uint8 = 0b1
	/// Set when current message is an acknowledgment for a previously received message.
	kExchangeFlagAckMsg uint8 = 0b10
	/// Set when current message is requesting an acknowledgment from the recipient.
	kExchangeFlagNeedsAck uint8 = 0b100
	/// Secured Extension block is present.
	kExchangeFlagSecuredExtension uint8 = 0b1000
	/// Set when a vendor id is prepended to the Message Protocol Id field.
	kExchangeFlagVendorIdPresent uint8 = 0b10000
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
 *  8 bit:  | Protocol Opcode   /Sigma1/Sigma2//Sigma4//Sigma1 Fin                 |
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
	mExchangeID        uint16
	mProtocolID        uint16
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
	return header.mExchangeFlags&kExchangeFlagInitiator != 0
}

func (header *PayloadHeader) IsAckMsg() bool {
	return header.mExchangeFlags&kExchangeFlagAckMsg != 0
}

func (header *PayloadHeader) NeedsAck() bool {
	return header.mExchangeFlags&kExchangeFlagNeedsAck != 0
}

func (header *PayloadHeader) HaveVendorId() bool {
	return header.mExchangeFlags&kExchangeFlagVendorIdPresent != 0
}

func (header *PayloadHeader) DecodeAndConsume(data []byte) error {

	header.mExchangeFlags = data[0]
	header.mProtocolOpcode = data[1]

	header.mExchangeID = binary.LittleEndian.Uint16(data[2:4])
	header.mProtocolID = binary.LittleEndian.Uint16(data[4:6])
	header.mSize = 6
	if header.HaveVendorId() {
		header.mVendorId = binary.LittleEndian.Uint16(data[header.mSize : header.mSize+2])
		header.mSize = header.mSize + 2
	}
	if header.IsAckMsg() {
		header.mVendorId = binary.LittleEndian.Uint16(data[header.mSize : header.mSize+2])
		header.mSize = header.mSize + 2
	}
	return nil
}

func (header *PayloadHeader) GetProtocolID() uint16 {
	return header.mProtocolID
}

func (header *PayloadHeader) GetMessageType() uint8 {
	return header.mMessageType
}
