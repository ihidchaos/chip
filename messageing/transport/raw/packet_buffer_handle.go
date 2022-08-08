package raw

import (
	"encoding/binary"
	log "github.com/sirupsen/logrus"
)

type PacketBuffer struct {
	Bytes   []byte
	Payload int
	Len     int
	TotLen  int
}

func NewPacketBuffer(data []byte) *PacketBuffer {
	return &PacketBuffer{
		Bytes:   data,
		Payload: 0,
		Len:     len(data),
		TotLen:  len(data),
	}
}

func (header *PacketBuffer) Read8() uint8 {
	if header.TotLen < 1 {
		log.Panic("Packet Buffer invalid")
	}
	value := header.Bytes[header.Payload]
	header.ConsumeHead(1)
	return value
}

func (header *PacketBuffer) Read16() uint16 {
	tSize := 2
	if header.TotLen < tSize {
		log.Panic("Packet Buffer invalid")
	}
	value := binary.LittleEndian.Uint16(header.Bytes[header.Payload : header.Payload+tSize])
	header.ConsumeHead(tSize)
	return value
}

func (header *PacketBuffer) Read32() uint32 {
	tSize := 4
	if header.TotLen < tSize {
		log.Panic("Packet Buffer invalid")
	}
	value := binary.LittleEndian.Uint32(header.Bytes[header.Payload : header.Payload+tSize])
	header.ConsumeHead(tSize)
	return value
}

func (header *PacketBuffer) Read64() uint64 {
	tSize := 8
	if header.TotLen < tSize {
		log.Panic("Packet Buffer invalid")
	}
	value := binary.LittleEndian.Uint64(header.Bytes[header.Payload : header.Payload+tSize])
	header.ConsumeHead(tSize)
	return value
}

func (header *PacketBuffer) ConsumeHead(aConsumeLength int) {
	if aConsumeLength > header.Len {
		aConsumeLength = header.Len
	}
	header.Payload = header.Payload + aConsumeLength
	header.Len = header.Len - aConsumeLength
	header.TotLen = header.Len - aConsumeLength
}
