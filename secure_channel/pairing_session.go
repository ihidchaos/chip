package secure_channel

import (
	"github.com/galenliu/chip/messageing"
	"github.com/galenliu/chip/transport"
)

type PairingSession struct {
	mRole           uint8
	mSessionManger  transport.SessionManager
	mLocalMRPConfig *messageing.ReliableMessageProtocolConfig
}
