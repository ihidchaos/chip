package transport

import (
	"github.com/galenliu/chip/messageing/transport/raw"
	"net/netip"
)

const kMessageCounterRandomInitMask uint32 = 0x0FFFFFFF
const (
	MessageCounterTypeGlobalUnencrypted = iota
	MessageCounterTypeGlobalEncrypted
	MessageCounterTypeSession
)

type MessageCounter struct {
}

type GlobalUnencryptedMessageCounter struct {
	*MessageCounter
}

type LocalSessionMessageCounter struct {
	*MessageCounter
}

func NewGlobalUnencryptedMessageCounterImpl() *GlobalUnencryptedMessageCounter {
	return &GlobalUnencryptedMessageCounter{}
}

func (g *GlobalUnencryptedMessageCounter) Init() {

}

func (g *GlobalUnencryptedMessageCounter) StartSync(handle *SessionHandle, session SecureSessionBase) error {
	//TODO implement me
	panic("implement me")
}

func (g *GlobalUnencryptedMessageCounter) QueueReceivedMessageAndStartSync(header *raw.PacketHeader, handle *SessionHandle, state uint8, peerAdders netip.AddrPort, buf *raw.PacketBuffer) error {
	//TODO implement me
	panic("implement me")
}
