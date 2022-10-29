package raw

import (
	"bytes"
	"errors"
	"github.com/galenliu/chip/lib"
	buffer "github.com/galenliu/chip/pkg/tlv/buffer"
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

func (buf *PacketBuffer) IsValid() error {
	if buf.DataLength() < 36 {
		return errors.New("message length is too short")
	}
	return nil
}

func (buf *PacketBuffer) PacketHeader() (*PacketHeader, error) {
	h := NewPacketHeader()
	var err error
	h.MessageFlags, err = buffer.Read8(buf)
	h.SessionId, err = buffer.LittleEndianRead16(buf)
	h.SecFlags, err = buffer.Read8(buf)
	h.MessageCounter, err = buffer.LittleEndianRead32(buf)
	if err != nil {
		return nil, err
	}
	if lib.HasFlags(h.MessageFlags, FSourceNodeIdPresent) {
		v, _ := buffer.LittleEndianRead64(buf)
		h.SourceNodeId = lib.NodeId(v)
	}
	if lib.HasFlags(h.MessageFlags, FDestinationNodeIdPresent) {
		v, _ := buffer.LittleEndianRead64(buf)
		h.DestinationNodeId = lib.NodeId(v)
	}
	if lib.HasFlags(h.MessageFlags, FDestinationGroupIdPresent) {
		v, _ := buffer.LittleEndianRead16(buf)
		h.DestinationGroupId = lib.GroupId(v)
	}
	return h, nil
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
