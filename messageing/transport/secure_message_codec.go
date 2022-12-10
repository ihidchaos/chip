package transport

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/platform/system"
)

func Decrypt(context *session.CryptoContext, nonce []byte, packetHeader *raw.PacketHeader, msg *system.PacketBufferHandle) (*raw.PayloadHeader, error) {

	if msg.IsNull() {
		return nil, lib.InvalidArgument
	}
	footerLen := packetHeader.MICTagLength()
	if int(footerLen) >= msg.Len() {
		return nil, lib.InvalidArgument
	}
	mac := raw.NewMessageAuthenticationCode()
	err := mac.Decode(packetHeader, msg, footerLen)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
