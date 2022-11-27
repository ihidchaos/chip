package system

import (
	"bytes"
	"errors"
)

type PacketBufferHandle struct {
	data []byte
	*bytes.Buffer
}

func NewPacketBufferHandle(data []byte) *PacketBufferHandle {
	return &PacketBufferHandle{
		data:   data,
		Buffer: bytes.NewBuffer(data),
	}
}

func (p *PacketBufferHandle) IsNull() bool {
	return p.Buffer.Len() == 0
}

func (p *PacketBufferHandle) Length() int {
	return p.Buffer.Len()
}

func (p *PacketBufferHandle) TotLength() int {
	return len(p.data)
}

func (p *PacketBufferHandle) IsValid() error {
	if p.Length() < 36 {
		return errors.New("message length is too short")
	}
	return nil
}
