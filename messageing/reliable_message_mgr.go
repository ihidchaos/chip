package messageing

import (
	"fmt"
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

func NewRetransTableEntry(rc *ReliableMessageContext) *RetransTableEntry {
	return &RetransTableEntry{
		ec: &ExchangeHandle{
			ExchangeContext: rc.mExchangeContext,
		},
		nextRetransTime: time.Time{},
		sendCount:       0,
	}
}

type ReliableMessageMgr struct {
	mRetransTable [config.RMPRetransTableSize]*RetransTableEntry
}

func (m *ReliableMessageMgr) ClearRetransTableEntry(re *RetransTableEntry) *RetransTableEntry {
	return nil
}

func (m *ReliableMessageMgr) ClearRetransTable(rc *ReliableMessageContext) {
	for _, entry := range m.mRetransTable {
		if entry.ec.ReliableMessageContext == rc {
			m.ClearRetransTableEntry(entry)
		}
	}
}

func (m *ReliableMessageMgr) CheckAndRemRetransTable(rc *ReliableMessageContext, ackMessageCounter uint32) bool {
	var removed = false
	for _, entry := range m.mRetransTable {
		if entry.ec.ReliableMessageContext == rc && entry.retainedBuf.MessageCounter() == ackMessageCounter {
			m.ClearRetransTableEntry(entry)
			log.Info("ExchangeManager Rxd Ack; Removing :", "MessageCounter", ackMessageCounter,
				" from Retrans Table on exchange ",
				ackMessageCounter, "from Retrans Table on exchange", rc.mExchangeContext)
			removed = true
		}
	}
	return removed
}

func (m *ReliableMessageMgr) AddToRetransTable(rc *ReliableMessageContext) (entry *RetransTableEntry, err error) {

	for i, e := range m.mRetransTable {
		if e != nil {
			m.mRetransTable[i] = NewRetransTableEntry(rc)
			entry = m.mRetransTable[i]
		}
	}
	if entry == nil {
		return entry, fmt.Errorf("RetransTable Already Full")
	}
	return
}
