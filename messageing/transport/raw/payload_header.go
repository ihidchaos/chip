package raw

import (
	"encoding/binary"
	"fmt"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/bitflags"
	"github.com/moznion/go-optional"
	log "golang.org/x/exp/slog"
	"io"
)

const (
	fInitiator uint8 = 0b1
	// fAckMsg / Sets when current message is an acknowledgment for a previously received message.
	fAckMsg uint8 = 0b10
	// fNeedsAck / Sets when current message is requesting an acknowledgment from the recipient.
	fNeedsAck uint8 = 0b100

	// FSecuredExtension / Secured Extension block is present.
	//FSecuredExtension uint8 = 0b1000

	// fVendorIdPresent / Sets when a vendor id is prepended to the Message Protocol Id field.
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
	mExchangeFlags     bitflags.Flags[uint8]
	mProtocolOpcode    uint8
	mExchangeId        uint16
	mProtocolId        uint16
	mVendorId          optional.Option[lib.VendorId]
	mAckMessageCounter optional.Option[uint32]
}

type payloadHeaderOption func(header *PayloadHeader)

func NewPayloadHeader(opts ...payloadHeaderOption) *PayloadHeader {
	header := &PayloadHeader{
		mExchangeFlags:  bitflags.Flags[uint8]{},
		mProtocolOpcode: 0,
		mExchangeId:     0,
		mProtocolId:     0,
	}
	for _, opt := range opts {
		opt(header)
	}
	return header
}

func (header *PayloadHeader) IsInitiator() bool {
	return header.mExchangeFlags.Has(fInitiator)
}

func (header *PayloadHeader) HasMessageType(t uint8) bool {
	return header.mProtocolOpcode == t
}

func (header *PayloadHeader) IsAckMsg() bool {
	return header.mExchangeFlags.Has(fAckMsg)
}

func (header *PayloadHeader) NeedsAck() bool {
	return header.mExchangeFlags.Has(fNeedsAck)
}

func (header *PayloadHeader) HaveVendorId() bool {
	return header.mExchangeFlags.Has(fVendorIdPresent)
}

func (header *PayloadHeader) DecodeAndConsume(buf io.Reader) error {

	data := make([]byte, 1)
	_, err := buf.Read(data)
	if err != nil {
		return err
	}
	header.mExchangeFlags = bitflags.Some(data[0])

	_, err = buf.Read(data)
	if err != nil {
		return err
	}
	header.mProtocolOpcode = data[0]

	data = make([]byte, 2)
	_, err = buf.Read(data)
	if err != nil {
		return err
	}
	header.mExchangeId = binary.LittleEndian.Uint16(data)

	_, err = buf.Read(data)
	if err != nil {
		return err
	}
	header.mProtocolId = binary.LittleEndian.Uint16(data)

	if header.HaveVendorId() {
		_, err = buf.Read(data)
		if err != nil {
			return err
		}
		vendorId := lib.VendorId(binary.LittleEndian.Uint16(data))
		header.mVendorId = optional.Some(vendorId)
	}

	if header.IsAckMsg() {
		data = make([]byte, 4)
		_, err = buf.Read(data)
		if err != nil {
			return err
		}
		header.mAckMessageCounter = optional.Some(binary.LittleEndian.Uint32(data))
	}
	return nil
}

func (header *PayloadHeader) AckMessageCounter() optional.Option[uint32] {
	return header.mAckMessageCounter
}

func (header *PayloadHeader) ProtocolID() uint16 {
	return header.mProtocolId
}

func (header *PayloadHeader) VendorId() optional.Option[lib.VendorId] {
	return header.mVendorId
}

func (header *PayloadHeader) MessageType() uint8 {
	return header.mProtocolOpcode
}

func (header *PayloadHeader) ExchangeId() uint16 {
	return header.mExchangeId
}

func (header *PayloadHeader) LogValue() log.Value {
	return log.GroupValue(
		log.String("exchangeId", fmt.Sprintf("%04X", header.mExchangeId)),
		log.Bool("IsInitiator", header.IsInitiator()),
	)
}
