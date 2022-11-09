package transport

import (
	"github.com/galenliu/chip/messageing/transport/raw"
	"net/netip"
)

type MessageCounterManagerBase interface {
	StartSync(handle *SessionHandle, session *SecureSession) error
	QueueReceivedMessageAndStartSync(
		packetHeader *raw.PacketHeader,
		session *SessionHandle,
		state *SecureSession,
		peerAdders netip.AddrPort,
		buf *raw.PacketBuffer,
	) error
}
