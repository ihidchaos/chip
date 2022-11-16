package transport

import (
	"github.com/galenliu/chip/messageing/transport/raw"
	"github.com/galenliu/chip/messageing/transport/session"
	"github.com/galenliu/chip/platform/system"
	"net/netip"
)

type MessageCounterManagerBase interface {
	StartSync(handle *SessionHandle, session *session.Secure) error
	QueueReceivedMessageAndStartSync(
		packetHeader *raw.PacketHeader,
		session *SessionHandle,
		state *session.Secure,
		peerAdders netip.AddrPort,
		buf *system.PacketBufferHandle,
	) error
}
