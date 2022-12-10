package system

import (
	"bytes"
	"errors"
)

type PacketBufferHandle struct {
	*bytes.Buffer
}

func NewPacketBufferHandle(data []byte) *PacketBufferHandle {
	return &PacketBufferHandle{
		Buffer: bytes.NewBuffer(data),
	}
}

func (p *PacketBufferHandle) IsNull() bool {
	return p.Buffer.Len() == 0
}

func (p *PacketBufferHandle) IsValid() error {
	if p.Len() < 36 {
		return errors.New("message length is too short")
	}
	return nil
}
