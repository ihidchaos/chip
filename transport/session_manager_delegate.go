package transport

import "github.com/galenliu/chip/transport/message"

type SessionMessageDelegate interface {
	OnMessageReceived(packetHeader *message.PacketHeader, payloadHeader *message.PayloadHeader, session Session, duplicate uint8, data []byte)
}
