package messageing

import (
	"github.com/galenliu/chip/config"
	"github.com/galenliu/chip/messageing/transport"
	log "golang.org/x/exp/slog"
	"time"
)

type RetransTableEntry struct {
	ec              *ExchangeHandle
	retainedBuf     *transport.EncryptedPacketBufferHandle
	nextRetransTime time.Time
	sendCount       uint8
}

type ReliableMessageMgr struct {
	mRetransTable [config.RMPRetransTableSize]*RetransTableEntry
}

func (m *ReliableMessageMgr) ClearRetransTable(entry *RetransTableEntry) {

}

func (m *ReliableMessageMgr) clearRetransTable(rc *ReliableMessageContext) {
	for _, entry := range m.mRetransTable {
		if entry.ec.ReliableMessageContext == rc {
			m.ClearRetransTable(entry)
		}
	}
}

func (m *ReliableMessageMgr) CheckAndRemRetransTable(rc *ReliableMessageContext, ackMessageCounter uint32) bool {
	var removed = false
	for _, entry := range m.mRetransTable {
		if entry.ec.ReliableMessageContext == rc && entry.retainedBuf.MessageCounter() == ackMessageCounter {
			m.ClearRetransTable(entry)
			log.Info("ExchangeManager Rxd Ack; Removing :", "MessageCounter", ackMessageCounter,
				" from Retrans Table on exchange ",
				ackMessageCounter, "from Retrans Table on exchange", rc.mExchangeContext)
			removed = true
		}
	}
	return removed
}
