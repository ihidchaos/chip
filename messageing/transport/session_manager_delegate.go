package transport

import (
	message2 "github.com/galenliu/chip/messageing/transport/raw"
)

type SessionMessageDelegate interface {
	OnMessageReceived(packetHeader *message2.PacketHeader, payloadHeader *message2.PayloadHeader, session Session, duplicate uint8, data []byte)
}
