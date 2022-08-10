package raw

import (
	"bytes"
	"encoding/binary"
)

type PacketBuffer struct {
	data   []byte
	buffer *bytes.Buffer
}

func NewPacketBuffer(data []byte) *PacketBuffer {
	return &PacketBuffer{
		data:   data,
		buffer: bytes.NewBuffer(data),
	}
}

func (buf *PacketBuffer) Read8() (uint8, error) {
	value, err := buf.buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (buf *PacketBuffer) Read16() (uint16, error) {
	tSize := 2
	data := make([]byte, tSize)
	_, err := buf.buffer.Read(data)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint16(data)
	return value, nil
}

func (buf *PacketBuffer) Read32() (uint32, error) {
	tSize := 4
	data := make([]byte, tSize)
	_, err := buf.buffer.Read(data)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint32(data)
	return value, nil
}

func (buf *PacketBuffer) Read64() (uint64, error) {
	tSize := 8
	data := make([]byte, tSize)
	_, err := buf.buffer.Read(data)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint64(data)
	return value, nil
}

func (buf *PacketBuffer) ReadBytes(data []byte) error {
	_, err := buf.buffer.Read(data)
	if err != nil {
		return err
	}
	return nil
}

func (buf *PacketBuffer) IsNull() bool {
	return buf.buffer.Len() == 0
}

func (buf *PacketBuffer) TotLength() int {
	return len(buf.data)
}

func (buf *PacketBuffer) DataLength() uint16 {
	return uint16(buf.buffer.Len())
}
