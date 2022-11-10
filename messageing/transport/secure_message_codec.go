package transport

import (
	"github.com/galenliu/chip/lib"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/platform/system"
)

func Decrypt(context *CryptoContext, nonce []byte, packetHeader *raw.PacketHeader, msg *system.PacketBufferHandle) error {

	if msg.IsNull() {
		return lib.MatterErrorInvalidArgument
	}
	footerLen := packetHeader.MICTagLength()
	if int(footerLen) >= msg.DataLength() {
		return lib.MatterErrorInvalidArgument
	}
	mac := raw.NewMessageAuthenticationCode()
	err := mac.Decode(packetHeader, msg, footerLen)
	if err != nil {
		return err
	}
	return nil
}
