package secure_channel

import (
	"github.com/galenliu/chip/messageing"
	transport2 "github.com/galenliu/chip/messageing/transport"
)

type PairingSession struct {
	mRole uint8

	SessionManager       transport2.SessionManager
	mExchangeCtxt        messageing.ExchangeContext
	mSecureSessionHolder transport2.SessionHolderWithDelegate
	mLocalMRPConfig      *transport2.ReliableMessageProtocolConfig
}

func NewPairingSession() *PairingSession {
	return &PairingSession{}
}
