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

func (c *MessageAuthenticationCode) Decode(header *PacketHeader, msg *lib.PacketBuffer, size uint16) error {
	tagLen := header.MICTagLength()
	if tagLen == 0 {
		return lib.ChipErrorWrongEncryptionTypeFromPeer
	}
	if size < tagLen {
		return lib.ChipErrorInvalidArgument
	}
	c.mTag = msg.GetData()[msg.TotLength()-int(tagLen) : msg.DataLength()]
	return nil
}

func (c *MessageAuthenticationCode) GetTag() []byte {
	return c.mTag
}
