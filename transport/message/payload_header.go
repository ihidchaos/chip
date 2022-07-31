package message

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

type PayloadHeader struct {
	mExchangeFlags     uint8
	mExchangeID        uint16
	mProtocolID        uint16
	mVendorId          uint16
	mAckMessageCounter uint32
	mLength            uint8
	mProtocolOpcode    uint8
}

func NewPayloadHeader(data []byte) *PayloadHeader {

	header := &PayloadHeader{}

	header.mExchangeFlags = data[0]
	header.mProtocolOpcode = data[1]
	header.mExchangeID = binary.LittleEndian.Uint16(data[2:4])
	header.mProtocolID = binary.LittleEndian.Uint16(data[4:6])
	header.mLength = 6
	if header.HaveVendorId() {
		header.mVendorId = binary.LittleEndian.Uint16(data[header.mLength : header.mLength+2])
	}
	if header.IsAckMsg() {
		header.mVendorId = binary.LittleEndian.Uint16(data[header.mLength : header.mLength+2])
	}

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

func (header *PayloadHeader) DecodeAndConsume(data []byte) {

}
