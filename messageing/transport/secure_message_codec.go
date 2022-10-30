package transport

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
)

func Decrypt(context *CryptoContext, nonce []byte, packetHeader *raw.PacketHeader, msg *raw.PacketBuffer) error {

	if msg.IsNull() {
		return lib.MatterErrorInvalidArgument
	}
	footerLen := packetHeader.MICTagLength()
	if footerLen >= msg.DataLength() {
		return lib.MatterErrorInvalidArgument
	}
	mac := raw.NewMessageAuthenticationCode()
	err := mac.Decode(packetHeader, msg, footerLen)
	if err != nil {
		return err
	}
	return nil
}
