package raw

import (
	"fmt"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/platform/system/buffer"
	"github.com/galenliu/chip/protocols"
	log "golang.org/x/exp/slog"
	"io"
)

const (
	fInitiator uint8 = 0b1
	// fAckMsg / Set when current message is an acknowledgment for a previously received message.
	fAckMsg uint8 = 0b10
	// fNeedsAck / Set when current message is requesting an acknowledgment from the recipient.
	fNeedsAck uint8 = 0b100

	// FSecuredExtension / Secured Extension block is present.
	//FSecuredExtension uint8 = 0b1000

	// fVendorIdPresent / Set when a vendor id is prepended to the Message Protocol Id field.
	fVendorIdPresent uint8 = 0b10000
)

/**********************************************
 * Header format (little endian):
 *
 * -------- Unencrypted header -----------------------------------------------------
 *  8 bit:  | Message Flags: VERSION: 4 bit | S: 1 bit | RESERVED: 1 bit | DSIZ: 2 bit |
 *  8 bit:  | Security Flags: P: 1 bit | C: 1 bit | MX: 1 bit | RESERVED: 3 bit | Session Type: 2 bit |
 *  16 bit: | Session ID                                                           |
 *  32 bit: | Message Counter                                                      |
 *  64 bit: | SOURCE_NODE_ID (iff source node flag is set)                         |
 *  64 bit: | DEST_NODE_ID (iff destination node flag is set)                      |
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

type PayloadHeader struct {
	mExchangeFlags     uint8
	mProtocolOpcode    uint8
	mExchangeId        uint16
	mProtocolId        protocols.Id
	mVendorId          lib.VendorId
	mAckMessageCounter uint32
}

type payloadHeaderOption func(header *PayloadHeader)

func NewPayloadHeader(opts ...payloadHeaderOption) *PayloadHeader {
	header := &PayloadHeader{
		mExchangeFlags:     0,
		mProtocolOpcode:    0,
		mExchangeId:        0,
		mProtocolId:        protocols.Id{},
		mVendorId:          0,
		mAckMessageCounter: 0,
	}
	for _, opt := range opts {
		opt(header)
	}
	return header
}

func (header *PayloadHeader) IsInitiator() bool {
	return lib.HasFlags(header.mExchangeFlags, fInitiator)
}

func (header *PayloadHeader) HasMessageType(t uint8) bool {
	return header.mProtocolOpcode == t
}

func (header *PayloadHeader) IsAckMsg() bool {
	return lib.HasFlags(header.mExchangeFlags, fAckMsg)
}

func (header *PayloadHeader) NeedsAck() bool {
	return lib.HasFlags(header.mExchangeFlags, fNeedsAck)
}

func (header *PayloadHeader) HaveVendorId() bool {
	return lib.HasFlags(header.mExchangeFlags, fVendorIdPresent)
}

func (header *PayloadHeader) DecodeAndConsume(buf io.Reader) error {
	var err error
	header.mExchangeFlags, err = buffer.Read8(buf)
	header.mProtocolOpcode, err = buffer.Read8(buf)
	header.mExchangeId, err = buffer.LittleEndianRead16(buf)
	protocolId, err := buffer.LittleEndianRead16(buf)
	if err != nil {
		return err
	}
	var vendorId = lib.VidCommon
	if header.HaveVendorId() {
		vid, err := buffer.LittleEndianRead16(buf)
		if err != nil {
			return err
		}
		vendorId = lib.VendorId(vid)
	}
	header.mProtocolId = protocols.NewProtocolId(vendorId, protocolId)
	if header.IsAckMsg() {
		ackCounter, err := buffer.LittleEndianRead32(buf)
		if err != nil {
			return err
		}
		header.mAckMessageCounter = ackCounter
	}
	return nil
}

func (header *PayloadHeader) AckMessageCounter() uint32 {
	return header.mAckMessageCounter
}

func (header *PayloadHeader) ProtocolID() *protocols.Id {
	return &header.mProtocolId
}

func (header *PayloadHeader) MessageType() uint8 {
	return header.mProtocolOpcode
}

func (header *PayloadHeader) ExchangeId() uint16 {
	return header.mExchangeId
}

func (header *PayloadHeader) HasProtocol(id *protocols.Id) bool {
	return header.mProtocolId.Equal(id)
}

func (header *PayloadHeader) LogValue() log.Value {
	return log.GroupValue(
		log.String("exchangeId", fmt.Sprintf("%04X", header.mExchangeId)),
		log.Bool("IsInitiator", header.IsInitiator()),
	)
}
