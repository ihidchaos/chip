package raw

import (
	"encoding/binary"
	"fmt"
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/lib/bitflags"
	"github.com/galenliu/chip/protocols"
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
 *  8 bit:  | Security Flags: P: 1 bit | C: 1 bit | MX: 1 bit | RESERVED: 3 bit | Session TransportType: 2 bit |
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
 *  <var>:  | (Unencrypted) Message Authentication NextTag                             |
 *
 **********************************************/

type PayloadHeader struct {
	mExchangeFlags    bitflags.Flags[uint8]
	MessageType       uint8
	ExchangeId        uint16
	ProtocolId        protocols.Id
	AckMessageCounter optional.Option[uint32]
}

type payloadHeaderOption func(header *PayloadHeader)

func NewPayloadHeader(opts ...payloadHeaderOption) *PayloadHeader {
	header := &PayloadHeader{
		mExchangeFlags: bitflags.Flags[uint8]{},
		MessageType:    0,
		ExchangeId:     0,
		ProtocolId:     protocols.NotSpecifiedId(),
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
	return header.MessageType == t
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

func (header *PayloadHeader) DecodeAndConsume(buf io.Reader) (err error) {

	data := make([]byte, 1)
	if _, err = buf.Read(data); err != nil {
		return err
	}
	header.mExchangeFlags = bitflags.Some(data[0])

	if _, err = buf.Read(data); err != nil {
		return err
	}
	header.MessageType = data[0]

	data = make([]byte, 2)
	if _, err = buf.Read(data); err != nil {
		return err
	}
	header.ExchangeId = binary.LittleEndian.Uint16(data)

	if _, err = buf.Read(data); err != nil {
		return err
	}
	protocolId := binary.LittleEndian.Uint16(data)

	vendorId := lib.VidCommon
	if header.HaveVendorId() {
		if _, err = buf.Read(data); err != nil {
			return err
		}
		vendorId = lib.VendorId(binary.LittleEndian.Uint16(data))
	}
	header.ProtocolId = protocols.New(protocolId, optional.Some(vendorId))

	if header.IsAckMsg() {
		data = make([]byte, 4)
		if _, err = buf.Read(data); err != nil {
			return err
		}
		header.AckMessageCounter = optional.Some(binary.LittleEndian.Uint32(data))
	}
	return nil
}

func (header *PayloadHeader) LogValue() log.Value {
	return log.GroupValue(
		log.String("ExchangeId", fmt.Sprintf("%04X", header.ExchangeId)),
		log.Uint64("MessageType", uint64(header.MessageType)),
		log.Bool("isInitiator", header.IsInitiator()),
		log.Any("ProtocolId", header.ProtocolId),
		log.Uint64("AckMessageCounter", uint64(header.AckMessageCounter.Unwrap())),
	)
}

func (header *PayloadHeader) SetInitiator(initiator bool) *PayloadHeader {
	header.mExchangeFlags.Set(initiator, fInitiator)
	return header
}

func (header *PayloadHeader) SetAckMessageCounter(counter uint32) *PayloadHeader {
	header.AckMessageCounter = optional.Some(counter)
	header.mExchangeFlags.Sets(fAckMsg)
	return header
}

func (header *PayloadHeader) SetNeedsAck(inNeedsAck bool) *PayloadHeader {
	header.mExchangeFlags.Set(inNeedsAck, fNeedsAck)
	return header
}

func (header *PayloadHeader) SetProtocol(id protocols.Id) *PayloadHeader {
	header.mExchangeFlags.Set(id.VendorId != lib.VidCommon, fVendorIdPresent)
	header.ProtocolId = id
	return header
}

func (header *PayloadHeader) SetMessageType(id protocols.Id, typ uint8) *PayloadHeader {
	header.MessageType = typ
	header.SetProtocol(id)
	return header
}

func (header *PayloadHeader) Encode(buf io.Writer) (err error) {
	if _, err = buf.Write([]byte{header.mExchangeFlags.Raw()}); err != nil {
		return err
	}
	if _, err = buf.Write([]byte{header.MessageType}); err != nil {
		return err
	}
	var data = make([]byte, 2)
	binary.LittleEndian.PutUint16(data, header.ExchangeId)
	if _, err = buf.Write(data); err != nil {
		return err
	}

	if header.HaveVendorId() {
		binary.LittleEndian.PutUint16(data, uint16(header.ProtocolId.VendorId))
		if _, err = buf.Write(data); err != nil {
			return err
		}
	}

	binary.LittleEndian.PutUint16(data, header.ProtocolId.ProtocolId)
	if _, err = buf.Write(data); err != nil {
		return err
	}

	if header.AckMessageCounter.IsSome() {
		data = make([]byte, 4)
		binary.LittleEndian.PutUint32(data, header.AckMessageCounter.Unwrap())
		if _, err = buf.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func (header *PayloadHeader) EncodeSizeBytes() uint16 {
	var size = kEncryptedHeaderSizeBytes
	if header.HaveVendorId() {
		size += kVendorIdSizeBytes
	}
	if header.AckMessageCounter.IsSome() {
		size += kAckMessageCounterSizeBytes
	}
	return uint16(size)
}
