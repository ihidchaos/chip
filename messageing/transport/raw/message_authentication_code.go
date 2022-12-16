package raw

import (
	"bytes"
	"github.com/galenliu/chip"
	"github.com/galenliu/chip/platform/system"
	"io"
)

const kMaxTagLen = 16

func NewMessageAuthenticationCode() *MessageAuthenticationCode {
	return &MessageAuthenticationCode{
		mMaxTagLen: kMaxTagLen,
		mTag:       nil,
	}
}

type MessageAuthenticationCode struct {
	mMaxTagLen uint8
	mTag       []byte
}

func (c *MessageAuthenticationCode) SetTag(header *PacketHeader, tag []byte) {
	c.mTag = tag
}

func (c *MessageAuthenticationCode) Tag() []byte {
	return c.mTag
}

func (c *MessageAuthenticationCode) Decode(header *PacketHeader, msg *system.PacketBufferHandle, size uint16) error {
	tagLen := header.MICTagLength()
	if tagLen == 0 {
		return chip.ErrorWrongEncryptionTypeFromPeer
	}
	if size < tagLen {
		return chip.ErrorInvalidArgument
	}
	c.mTag = msg.Bytes()[msg.Len()-int(tagLen) : msg.Len()]
	return nil
}

func (c *MessageAuthenticationCode) Encode(header *PacketHeader, writer io.Writer) error {
	buf := new(bytes.Buffer)
	if len(c.mTag) != int(header.MICTagLength()) {
		return chip.ErrorInvalidArgument
	}
	_, err := buf.WriteTo(writer)
	return err
}
