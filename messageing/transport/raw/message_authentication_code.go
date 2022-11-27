package raw

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/platform/system"
)

const kMaxTagLen = 16
const kMaxAppMessageLen = 1200

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

func (c *MessageAuthenticationCode) Decode(header *PacketHeader, msg *system.PacketBufferHandle, size uint16) error {
	tagLen := header.MICTagLength()
	if tagLen == 0 {
		return lib.WrongEncryptionTypeFromPeer
	}
	if size < tagLen {
		return lib.InvalidArgument
	}
	c.mTag = msg.Bytes()[msg.TotLength()-int(tagLen) : msg.Length()]
	return nil
}

func (c *MessageAuthenticationCode) Tag() []byte {
	return c.mTag
}
