package message

import "github.com/galenliu/chip/transport"

type SessionMessageDelegate interface {
	OnMessageReceived(packetHeader *PacketHeader, payloadHeader *PayloadHeader, session transport.Session, duplicate uint8, data []byte)
}
