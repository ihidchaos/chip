package message

import "github.com/galenliu/chip/transport"

type SessionMessageDelegate interface {
	OnMessageReceived(header *Header, header2 PayloadHeader, session transport.Session, duplicate uint8, data []byte)
}

type ExchangeManager struct {
}
