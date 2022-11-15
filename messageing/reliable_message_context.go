package messageing

import (
	"github.com/galenliu/chip/lib"
	"time"
)

const (
	kFlagInitiator                uint16 = 1 << iota
	kFlagResponseExpected         uint16 = 1 << 1
	kFlagAutoRequestAck           uint16 = 1 << 2
	kFlagMessageNotAcked          uint16 = 1 << 3
	kFlagAckPending               uint16 = 1 << 4
	kFlagAckMessageCounterIsValid uint16 = 1 << 5
	kFlagWillSendMessage          uint16 = 1 << 6
	kFlagClosed                   uint16 = 1 << 7
	kFlagActiveMode               uint16 = 1 << 8
	kFlagEphemeralExchange        uint16 = 1 << 9
	kFlagIgnoreSessionRelease     uint16 = 1 << 10
)

type ReliableMessageContext struct {
	mFlags                        uint16
	mNextAckTime                  time.Time
	mPendingPeerAckMessageCounter uint32
}

func NewReliableMessageContext() ReliableMessageContext {
	return ReliableMessageContext{
		mNextAckTime:                  time.Time{},
		mPendingPeerAckMessageCounter: 0,
	}
}

func (c *ReliableMessageContext) HandleRcvdAck(ackCounter uint32) {
	//TODO implement me
	panic("implement me")
}

func (c *ReliableMessageContext) HandleNeedsAck(messageCounter, flags uint32) {
	//TODO implement me
	panic("implement me")
}

func (c *ReliableMessageContext) FlushAcks() error {
	//TODO implement me
	panic("implement me")
}

func (c *ReliableMessageContext) IsEphemeralExchange() bool {
	return lib.HasFlags(c.mFlags, kFlagMessageNotAcked)
}

func (c *ReliableMessageContext) IsMessageNotAcked() bool {
	return lib.HasFlags(c.mFlags, kFlagEphemeralExchange)
}

func (c *ReliableMessageContext) HasPiggybackAckPending() bool {
	return lib.HasFlags(c.mFlags, kFlagAckMessageCounterIsValid)
}

func (c *ReliableMessageContext) SetAckPending(b bool) {
	c.mFlags = lib.SetFlag(b, c.mFlags, kFlagAckPending)
}

func (c *ReliableMessageContext) IsAckPending() bool {
	return lib.HasFlags(c.mFlags, kFlagAckPending)
}

func (c *ReliableMessageContext) AutoRequestAck() bool {
	return lib.HasFlags(c.mFlags, kFlagAutoRequestAck)
}

func (c *ReliableMessageContext) IsRequestingActiveMode() bool {
	return lib.HasFlags(c.mFlags, kFlagActiveMode)
}

func (c *ReliableMessageContext) ShouldIgnoreSessionRelease() bool {
	return lib.HasFlags(c.mFlags, kFlagIgnoreSessionRelease)
}
