package transport

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
)

func Decrypt(context *CryptoContext, nonce []byte, packetHeader *raw.PacketHeader, msg *lib.PacketBuffer) error {

	if msg.IsNull() {
		return lib.ChipErrorInvalidArgument
	}
	footerLen := packetHeader.MICTagLength()
	if footerLen >= msg.DataLength() {
		return lib.ChipErrorInvalidArgument
	}
	mac := raw.NewMessageAuthenticationCode()
	err := mac.Decode(packetHeader, msg, footerLen)
	if err != nil {
		return err
	}
	return nil
}
