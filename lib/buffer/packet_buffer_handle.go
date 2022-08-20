package buffer

import (
	"bytes"
)

type PacketBuffer struct {
	data []byte
	*bytes.Buffer
}

func NewPacketBuffer(data []byte) *PacketBuffer {
	return &PacketBuffer{
		data:   data,
		Buffer: bytes.NewBuffer(data),
	}
}

func (buf *PacketBuffer) IsNull() bool {
	return buf.Buffer.Len() == 0
}

func (buf *PacketBuffer) TotLength() int {
	return len(buf.data)
}

func (buf *PacketBuffer) DataLength() uint16 {
	return uint16(buf.Buffer.Len())
}
