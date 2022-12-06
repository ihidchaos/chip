package system

import (
	"bytes"
	"errors"
)

type PacketBufferHandle struct {
	data []byte
	r    *bytes.Buffer
}

func NewPacketBufferHandle(data []byte) *PacketBufferHandle {
	return &PacketBufferHandle{
		data: data,
		r:    bytes.NewBuffer(data),
	}
}

func (p *PacketBufferHandle) Read(data []byte) (int, error) {
	return p.r.Read(data)
}

func (p *PacketBufferHandle) Bytes() []byte {
	return p.r.Bytes()
}

func (p *PacketBufferHandle) IsNull() bool {
	return p.r.Len() == 0
}

func (p *PacketBufferHandle) Length() int {
	return p.r.Len()
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
