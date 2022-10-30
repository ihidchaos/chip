package raw

import (
	"github.com/galenliu/chip/lib"
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

func (c *MessageAuthenticationCode) Decode(header *PacketHeader, msg *PacketBuffer, size uint16) error {
	tagLen := header.MICTagLength()
	if tagLen == 0 {
		return lib.MatterErrorWrongEncryptionTypeFromPeer
	}
	if size < tagLen {
		return lib.MatterErrorInvalidArgument
	}
	c.mTag = msg.Bytes()[msg.TotLength()-int(tagLen) : msg.DataLength()]
	return nil
}

func (c *MessageAuthenticationCode) Tag() []byte {
	return c.mTag
}
