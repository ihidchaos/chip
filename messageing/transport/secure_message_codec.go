package transport

import (
	"bytes"
	"github.com/galenliu/chip"
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/platform/system"
)

func Encrypt(context *CryptoContext, nonce []byte, payloadHeader *raw.PayloadHeader, header *raw.PacketHeader, buf *bytes.Buffer) (err error) {
	if buf.Len() == 0 || buf.Len() > raw.MaxAppMessageLen {
		return err.MATTER_ERROR_INVALID_MESSAGE_LENGTH
	}
	if err = payloadHeader.Encode(buf); err != nil {
		return err
	}
	cipherTag, err := context.Encrypt(buf.Bytes(), nonce, header)
	if err != nil {
		return err
	}

	return nil
}

func Decrypt(context *CryptoContext, nonce []byte, packetHeader *raw.PacketHeader, msg *system.PacketBufferHandle) (*raw.PayloadHeader, error) {
	if msg.IsNull() {
		return nil, chip.ErrorInvalidArgument
	}
	footerLen := packetHeader.MICTagLength()
	if int(footerLen) >= msg.Len() {
		return nil, chip.ErrorInvalidArgument
	}
	mac := raw.NewMessageAuthenticationCode()
	err := mac.Decode(packetHeader, msg, footerLen)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
